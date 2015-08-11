package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"gitHub.***REMOVED***/monsoon/arc/fact/host"

	log "github.com/Sirupsen/logrus"
	"github.com/kardianos/osext"
)

var runitRunScript = template.Must(template.New("run").Parse(`#!/bin/sh
exec 2>&1
exec {{ .executable }} -c {{.arcDir}}/arc.cfg server
`))

var runitLogScript = template.Must(template.New("log").Parse(`#!/bin/sh
exec {{ .serviceDir }}/svlogd -tt {{ .arcDir }}/log
`))

var runitFinishScript = template.Must(template.New("finish").Parse(`#!/bin/sh
if "$1" == "-1";then
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

func Install(dir string) error {
	executable, err := osext.Executable()
	if err != nil {
		return errors.New("Can't locate running executable")
	}
	if err := os.MkdirAll(path.Join(dir, "log"), 0755); err != nil {
		return err
	}

	serviceDir := path.Join(dir, "service")
	if err = installRunitSupervisor(executable, dir, serviceDir); err != nil {
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
	return fmt.Errorf("Unknown init system: %s", init)
}

func detectInitSystem() (string, error) {
	var hostFacts map[string]interface{}
	var err error
	if hostFacts, err = host.New().Facts(); err != nil {
		return "", errors.New("Can't detect init system")
	}
	init, ok := hostFacts["init_package"].(string)
	if !ok {
		return "", errors.New("Can't detect init system")
	}
	return init, nil
}

func installRunitSupervisor(executable, arcDir, serviceDir string) error {
	log.Info("Installing supervisor")

	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return err
	}
	//Write our the runit binaries
	for _, exe := range []string{"runsv", "svlogd"} {
		err := ioutil.WriteFile(path.Join(serviceDir, exe), FSMustByte(false, "/"+exe), 0755)
		if err != nil {
			return err
		}
	}

	templateVars := map[string]string{
		"arcDir":     arcDir,
		"serviceDir": serviceDir,
		"executable": executable,
	}

	runFile, err := os.OpenFile(path.Join(serviceDir, "run"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer runFile.Close()
	err = runitRunScript.Execute(runFile, templateVars)
	if err != nil {
		return err
	}
	finishFile, err := os.OpenFile(path.Join(serviceDir, "finish"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	err = runitFinishScript.Execute(finishFile, templateVars)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(serviceDir, "log"), 0755)
	if err != nil {
		return err
	}
	logFile, err := os.OpenFile(path.Join(serviceDir, "log", "run"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer logFile.Close()
	err = runitLogScript.Execute(logFile, templateVars)
	if err != nil {
		return err
	}

	return nil
}

func systemdService(cmd string) error {
	log.Info("Creating systemd service")

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

	if out, err := exec.Command("systemctl", "enable", "arc").CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to enable systemd service: %s", string(out))
	}

	if out, err := exec.Command("systemctl", "start", "arc").CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to start systemd service: %s", string(out))
	}

	return nil

}

func upstartService(cmd string) error {
	log.Info("Creating upstart job")
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

	if out, err := exec.Command("start", "arc").CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to job: %s", string(out))
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
	inittabFile, err := os.OpenFile(path.Join("/etc/inittab"), os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if _, err = inittabFile.WriteString(fmt.Sprintf("\nRI:2345:respawn:%s\n", cmd)); err != nil {
		inittabFile.Close()
		return err
	}
	inittabFile.Close()
	if out, err := exec.Command("telinit", "q").CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to reload inittab: %s", string(out))
	}
	return nil
}
