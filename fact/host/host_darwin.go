package host

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"gitHub.***REMOVED***/monsoon/arc/api-server/pki"
)

func (h Source) Facts() (map[string]interface{}, error) {
	facts := make(map[string]interface{})
	facts["os"] = runtime.GOOS
	facts["platform"] = "mac_os_x"
	facts["platform_family"] = "mac_os_x"
	facts["platform_version"] = nil
	facts["hostname"] = nil
	facts["fqdn"] = nil
	facts["domain"] = nil

	if notAfter, err := pki.CertExpirationDate(h.Config); err == nil {
		hoursLeft := pki.CertExpiresIn(notAfter)
		facts["cert_expiration"] = hoursLeft
	}

	if hostname, err := os.Hostname(); err == nil {
		facts["hostname"] = hostname
	}

	if out, err := exec.Command("/usr/bin/sw_vers").Output(); /* #nosec */ err == nil {
		re := regexp.MustCompile(`ProductVersion:\s+(.+)`)
		for _, line := range strings.Split(string(out), "\n") {
			if match := re.FindStringSubmatch(line); match != nil {
				facts["platform_version"] = match[1]
			}
		}
	}

	if fqdn, domain := fqdn_and_domain(); fqdn != "" {
		facts["fqdn"] = fqdn
		facts["domain"] = domain
	}

	return facts, nil
}
