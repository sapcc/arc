package chef

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/sapcc/arc/agents/chef/sap-ca"
)

var (
	chefClientBinary = "/usr/bin/chef-client"
	chefSoloBinary   = "/usr/bin/chef-solo"
	eol              = "\n"
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
		return fmt.Errorf("unknown package format: %s", installerType)
	} // #nosec

	log.Infof("Running %s", strings.Join(cmd.Args, " "))

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warnf("Command failed: %s", output)
		return err
	}

	if err := sapCa.AddSAPCAsToChefBundle("/opt/chef/embedded/ssl/certs/cacert.pem"); err != nil {
		return err
	}

	return nil
}
