package host

import (
	"os/exec"
	"regexp"
	"strings"
)

type Source struct{}

func New() Source {
	return Source{}
}

func (h Source) Name() string {
	return "host"
}

func fqdn_and_domain() (fqdn, domain string) {
	cmd := exec.Command("hostname", "-f")
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
