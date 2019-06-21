/* In cases where Gas reports a failure that has been manually verified as being
safe it is possible to annotate the code with a '#nosec' comment. */
package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"text/template"

	"github.com/sapcc/arc/fact/host"

	log "github.com/Sirupsen/logrus"
	"github.com/kardianos/osext"
)

var runitRunScript = template.Must(template.New("run").Parse(`#!/bin/sh
exec 2>&1
[ -f /etc/profile.d/proxy_settings.sh ] && . /etc/profile.d/proxy_settings.sh
exec {{ .executable }} -c {{.arcDir}}/arc.cfg server
`))

var runitLogScript = template.Must(template.New("log").Parse(`#!/bin/sh
exec {{ .serviceDir }}/svlogd {{ .arcDir }}/log
`))

var runitFinishScript = template.Must(template.New("finish").Parse(`#!/bin/sh
if [ "$1" = "-1" ]; then
  echo -n "process exited by signal "
  kill -l $2
else
  echo process exited with exit code $1
fi
`))

var upstartScript = template.Must(template.New("upstart").Parse(`
# arc supervisor
start on runlevel [2345]
stop on runlevel [^2345]
normal exit 0 111
respawn
exec {{ .cmd }}
`))

var systemdScript = template.Must(template.New("systemd").Parse(`
[Unit]
Description=Arc Process Supervisor

[Service]
ExecStart={{ .cmd }}
Restart=always
RestartSec=1

[Install]
WantedBy=multi-user.target
`))

func (s service) Status() (State, string, error) {
	cmd := s.svCmd("status", "service")
	out, err := cmd.CombinedOutput()
	message := strings.TrimSuffix(string(out), "\n")
	if rc, ok := err.(*exec.ExitError); ok {
		if rc.Sys().(syscall.WaitStatus).ExitStatus() == 3 {
			return STOPPED, message, nil
		}
		return UNKNOWN, message, err
	}
	if err != nil {
		return UNKNOWN, err.Error(), err
	}
	return RUNNING, message, nil
}

func (s service) Start() error {
	return s.svCmd("start", "service").Run()
}

func (s service) Stop() error {
	return s.svCmd("stop", "service").Run()
}

func (s service) Restart() error {
	return s.svCmd("term", "service").Run()
}

func (s service) Install() error {
	executable, err := osext.Executable()
	if err != nil {
		return errors.New("can't locate running executable")
	}
	if err := os.MkdirAll(path.Join(s.dir, "log"), 0700); err != nil {
		return err
	}

	serviceDir := path.Join(s.dir, "service")
	if err = installRunitSupervisor(executable, s.dir, serviceDir); err != nil {
		return err
	}
	init, err := detectInitSystem()
	if err != nil {
		return err
	}

	serviceCmd := fmt.Sprintf("%s %s", path.Join(serviceDir, "runsv"), serviceDir)

	switch init {
	case "systemd":
		return systemdService(serviceCmd)
	case "upstart":
		return upstartService(serviceCmd)
	case "sysv":
		return sysvService(serviceCmd)
	}
	return fmt.Errorf("unknown init system: %s", init)
}

func (s service) svCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(path.Join(s.dir, "service", "service"), args...) // #nosec
	cmd.Env = []string{fmt.Sprintf("SVDIR=%s", s.dir)}
	return cmd
}

func detectInitSystem() (string, error) {
	var hostFacts map[string]interface{}
	var err error
	if hostFacts, err = host.New(nil).Facts(); err != nil {
		return "", errors.New("can't detect init system")
	}
	init, ok := hostFacts["init_package"].(string)
	if !ok {
		return "", errors.New("can't detect init system")
	}
	return init, nil
}

func installRunitSupervisor(executable, arcDir, serviceDir string) error {
	log.Info("Installing supervisor")

	if err := os.MkdirAll(serviceDir, 0755); /* #nosec */ err != nil {
		return err
	}
	//Write our the runit binaries
	for _, exe := range []string{"runsv", "svlogd", "sv"} {
		err := ioutil.WriteFile(path.Join(serviceDir, exe), FSMustByte(false, "/"+exe), 0755)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(path.Join(serviceDir, "service")); err == nil {
		os.Remove(path.Join(serviceDir, "service"))
	}
	if err := os.Symlink("sv", path.Join(serviceDir, "service")); err != nil {
		return err
	}

	templateVars := map[string]string{
		"arcDir":     arcDir,
		"serviceDir": serviceDir,
		"executable": executable,
	}
	/* #nosec */
	runFile, err := os.OpenFile(path.Join(serviceDir, "run"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer runFile.Close()
	err = runitRunScript.Execute(runFile, templateVars)
	if err != nil {
		return err
	}
	/* #nosec */
	finishFile, err := os.OpenFile(path.Join(serviceDir, "finish"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	err = runitFinishScript.Execute(finishFile, templateVars)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(serviceDir, "log"), 0700)
	if err != nil {
		return err
	}
	/* #nosec */
	logFile, err := os.OpenFile(path.Join(serviceDir, "log", "run"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer logFile.Close()
	return runitLogScript.Execute(logFile, templateVars)
}

func systemdService(cmd string) error {
	log.Info("Creating systemd service")

	/* #nosec */
	unitFile, err := os.OpenFile(path.Join("/etc/systemd/system/arc.service"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	err = systemdScript.Execute(unitFile, map[string]string{"cmd": cmd})
	if err != nil {
		unitFile.Close()
		return err
	}
	unitFile.Close()

	if out, err := exec.Command("/bin/systemctl", "enable", "arc").CombinedOutput(); /* #nosec */ err != nil {
		return fmt.Errorf("failed to enable systemd service: %s", string(out))
	}

	if out, err := exec.Command("/bin/systemctl", "start", "arc").CombinedOutput(); /* #nosec */ err != nil {
		return fmt.Errorf("failed to start systemd service: %s", string(out))
	}

	return nil

}

func upstartService(cmd string) error {
	log.Info("Creating upstart job")
	/* #nosec */
	upstartFile, err := os.OpenFile(path.Join("/etc/init/arc.conf"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	err = upstartScript.Execute(upstartFile, map[string]string{"cmd": cmd})
	if err != nil {
		upstartFile.Close()
		return err
	}
	upstartFile.Close()

	if out, err := exec.Command("start", "arc").CombinedOutput(); /* #nosec */ err != nil {
		return fmt.Errorf("failed to job: %s", string(out))
	}

	return nil
}

func sysvService(cmd string) error {
	log.Info("Adding service to inittab")
	inittab, err := ioutil.ReadFile("/etc/inittab")
	if err != nil {
		return err
	}
	if strings.Contains(string(inittab), cmd) {
		//nothing to do cmd is already in inittab
		return nil
	}
	/* #nosec */
	inittabFile, err := os.OpenFile(path.Join("/etc/inittab"), os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if _, err = inittabFile.WriteString(fmt.Sprintf("\nRI:2345:respawn:%s\n", cmd)); err != nil {
		inittabFile.Close()
		return err
	}
	inittabFile.Close()
	if out, err := exec.Command("telinit", "q").CombinedOutput(); /* #nosec */ err != nil {
		return fmt.Errorf("failed to reload inittab: %s", string(out))
	}
	return nil
}
