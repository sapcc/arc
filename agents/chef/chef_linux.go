package chef

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	chef_binary = "/usr/bin/chef-client"
)

func install(installer string) error {
	installerType := regexp.MustCompile(`[^.]+$`).FindString(installer)
	var cmd *exec.Cmd
	switch installerType {
	case "rpm":
		cmd = exec.Command("rpm", "-U", installer)
	case "deb":
		cmd = exec.Command("dpkg", "-i", installer)
	default:
		return fmt.Errorf("Unknown package format: %s", installerType)
	}

	log.Infof("Running %s", strings.Join(cmd.Args, " "))

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warnf("Command failed: %s", output)
		return err
	}

	return nil
}
