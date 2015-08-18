package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/kardianos/osext"
)

func (s service) Status() (State, string, error) {
	out, err := s.nssmCmd("status").CombinedOutput()
	message := strings.TrimSpace(string(filterNullBytes(out)))
	switch message {
	case "SERVICE_RUNNING":
		return RUNNING, message, err
	case "SERVICE_PAUSED", "SERVICE_STOPPED":
		return STOPPED, message, err
	}
	return UNKNOWN, message, err
}

func (s service) Restart() error {
	_, err := s.nssmCmd("restart").CombinedOutput()
	return err
}

func (s service) Start() error {
	_, err := s.nssmCmd("start").CombinedOutput()
	return err
}
func (s service) Stop() error {
	_, err := s.nssmCmd("stop").CombinedOutput()
	return err
}

func (s service) nssmCmd(cmd string, args ...string) *exec.Cmd {
	args = append([]string{cmd, serviceName}, args...)
	return exec.Command(filepath.Join(s.dir, "nssm.exe"), args...)
}

func (s service) Install() error {
	executable, err := osext.Executable()
	if err != nil {
		return errors.New("Can't locate running executable")
	}

	if err := os.MkdirAll(path.Join(s.dir, "log"), 0755); err != nil {
		return err
	}

	return s.installNSSM(executable)

}

func (s service) installNSSM(executable string) error {
	log.Info("Installing the NSSM supervisor")
	nssm := filepath.Join(s.dir, "nssm.exe")
	err := ioutil.WriteFile(nssm, FSMustByte(false, "/nssm.exe"), 0755)
	if err != nil {
		return err
	}

	//Remove any previously created service
	s.nssmCmd("stop").Run()
	s.nssmCmd("remove", "confirm").Run()

	if out, err := s.nssmCmd("install", executable, "-c", path.Join(s.dir, "arc.cfg"), "server").CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to install service: %s err: %v", filterNullBytes(out), err)
	}

	settings := [][]string{
		[]string{"Description", serviceDescription},
		[]string{"DisplayName", serviceDisplayName},
		[]string{"AppStdout", path.Join(s.dir, "log", "current")},
		[]string{"AppStderr", path.Join(s.dir, "log", "current")},
		[]string{"AppRotateFiles", "1"},
		[]string{"AppRotateBytes", "100000"},
		[]string{"AppRotateOnline", "1"},
		[]string{"AppStopMethodSkip", "6"},
		[]string{"AppStopMethodConsole", "2000"},
	}
	for _, setting := range settings {
		if out, err := s.nssmCmd("set", setting...).CombinedOutput(); err != nil {
			return fmt.Errorf("Failed to set service option %s: %s", setting[0], string(out))
		}
	}

	log.Info("Starting service")
	if out, err := s.nssmCmd("start", serviceName).CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to start service: %s", filterNullBytes(out))
	}

	return nil
}

// No idea why, but the output of the nssm binary execution contains a null byte after every character
// we filter this out to keep the errors readble
func filterNullBytes(in []byte) []byte {
	filtered := make([]byte, 0, len(in)/2)
	for _, b := range in {
		if b != 0 {
			filtered = append(filtered, b)
		}
	}
	return filtered
}
