//+build darwin linux

package host

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/common"
	"github.com/shirou/gopsutil/host"
)

func (h Source) Facts() (map[string]interface{}, error) {

	info, err := host.HostInfo()
	if err != nil {
		return nil, err
	}

	facts := make(map[string]interface{})
	facts["os"] = info.OS
	facts["platform"] = info.Platform
	facts["platform_family"] = info.PlatformFamily
	facts["platform_version"] = info.PlatformVersion
	facts["fqdn"] = nil
	facts["domain"] = nil

	facts["hostname"] = info.Hostname
	cmd := exec.Command("hostname", "-f")
	if out, err := cmd.Output(); err == nil {
		fqdn_str := strings.TrimSpace(string(out))
		domain_regexp := regexp.MustCompile(`.*?\.(.+)$`)
		if m := domain_regexp.FindStringSubmatch(fqdn_str); m != nil {
			facts["fqdn"] = fqdn_str
			facts["domain"] = m[1]
		}
	}

	if common.PathExists("/etc/SuSE-release") {
		contents, err := common.ReadLines("/etc/SuSE-release")
		if err == nil {
			facts["platform_family"] = "suse"
			facts["platform"] = getSusePlatform(contents)
			facts["platform_version"] = getSuseVersion(contents)

		}
	}

	return facts, nil
}

func getSuseVersion(contents []string) string {
	version := ""
	for _, line := range contents {
		if matches := regexp.MustCompile(`VERSION = ([\d.]+)`).FindStringSubmatch(line); matches != nil {
			version = matches[1]
		} else if matches := regexp.MustCompile(`PATCHLEVEL = ([\d]+)`).FindStringSubmatch(line); matches != nil {
			version = version + "." + matches[1]
		}
	}
	return version
}

func getSusePlatform(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))
	if strings.Contains(c, "opensuse") {
		return "opensuse"
	}
	return "suse"
}
