package chef

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	chefClientBinary = "/usr/bin/chef-client"
	chefSoloBinary   = "/usr/bin/chef-solo"
)

func install(installer string) error {
	installerType := regexp.MustCompile(`[^.]+$`).FindString(installer)
	var cmd *exec.Cmd
	switch installerType {
	case "rpm":
		cmd = exec.Command("/bin/rpm", "-U", installer)
	case "deb":
		cmd = exec.Command("/usr/bin/dpkg", "-i", installer)
	default:
		return fmt.Errorf("Unknown package format: %s", installerType)
	} // #nosec

	log.Infof("Running %s", strings.Join(cmd.Args, " "))

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warnf("Command failed: %s", output)
		return err
	}

	if err := addSAPCAsToChefBundle("/opt/chef/embedded/ssl/certs/cacert.pem"); err != nil {
		return err
	}

	return nil
}
