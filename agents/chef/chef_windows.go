package chef

import (
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/sapcc/arc/agents/chef/sap-ca"
)

var (
	chefClientBinary = "C:/opscode/chef/bin/chef-client.bat"
	chefSoloBinary   = "C:/opscode/chef/bin/chef-solo.bat"
	eol              = "\r\n"
)

func install(installer string) error {
	cmd := exec.Command("msiexec", "/qn", "/i", installer) // #nosec
	log.Infof("Running %s", strings.Join(cmd.Args, " "))
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warnf("Command failed: %s", output)
		return err
	}

	if err := sapCa.AddSAPCAsToChefBundle("C:/opscode/chef/embedded/ssl/certs/cacert.pem"); err != nil {
		return err
	}

	return nil
}
