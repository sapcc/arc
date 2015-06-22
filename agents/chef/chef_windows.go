package chef

import (
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	chef_binary = "C:/opscode/chef/bin/chef-client.bat"
)

func install(installer string) error {
	cmd := exec.Command("msiexec", "/qn", "/i", installer)
	log.Infof("Running %s", strings.Join(cmd.Args, " "))
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warnf("Command failed: %s", output)
		return err
	}

	return nil
}
