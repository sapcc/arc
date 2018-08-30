package host

import (
	"os/exec"
	"regexp"
	"strings"

	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
)

type Source struct {
	Config *arc_config.Config
}

func New(cfg *arc_config.Config) Source {
	return Source{
		Config: cfg,
	}
}

func (h Source) Name() string {
	return "host"
}

func fqdn_and_domain() (fqdn, domain string) {
	cmd := exec.Command("/bin/hostname", "-f") // #nosec
	if out, err := cmd.Output(); err == nil {
		fqdn_str := strings.TrimSpace(string(out))
		domain_regexp := regexp.MustCompile(`.*?\.(.+)$`)
		if m := domain_regexp.FindStringSubmatch(fqdn_str); m != nil {
			fqdn = fqdn_str
			domain = m[1]
		}
	}

	return
}
