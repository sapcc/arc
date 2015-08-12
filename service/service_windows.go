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

func Status(dir string) (string, error) {
	panic("Not implemented on this platform")
}

func Install(dir string) error {
	executable, err := osext.Executable()
	if err != nil {
		return errors.New("Can't locate running executable")
	}

	if err := os.MkdirAll(path.Join(dir, "log"), 0755); err != nil {
		return err
	}

	return installNSSM(executable, dir)

}

func installNSSM(executable, installDir string) error {
	log.Info("Installing the NSSM supervisor")
	nssm := filepath.Join(installDir, "nssm.exe")
	err := ioutil.WriteFile(nssm, FSMustByte(false, "/nssm.exe"), 0755)
	if err != nil {
		return err
	}

	//Remove any previously created service
	exec.Command(nssm, "stop", serviceName).Run()
	exec.Command(nssm, "remove", serviceName, "confirm").Run()

	installArgs := []string{"install", serviceName, executable, "-c", path.Join(installDir, "arc.cfg"), "server"}
	log.Debugf("Running %s %s", nssm, strings.Join(installArgs, " "))
	if out, err := exec.Command(nssm, installArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to install service: %s err: %v", filterNullBytes(out), err)
	}

	settings := [][]string{
		[]string{"Description", serviceDescription},
		[]string{"DisplayName", serviceDisplayName},
		[]string{"AppStdout", path.Join(installDir, "log", "current")},
		[]string{"AppStderr", path.Join(installDir, "log", "current")},
		[]string{"AppRotateFiles", "1"},
		[]string{"AppRotateBytes", "100000"},
		[]string{"AppRotateOnline", "1"},
		[]string{"AppStopMethodSkip", "6"},
		[]string{"AppStopMethodConsole", "2000"},
	}
	for _, s := range settings {
		if out, err := exec.Command(nssm, "set", serviceName, s[0], s[1]).CombinedOutput(); err != nil {
			return fmt.Errorf("Failed to set service option %s: %s", s[0], string(out))
		}
	}

	log.Info("Starting service")
	if out, err := exec.Command(nssm, "start", serviceName).CombinedOutput(); err != nil {
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
