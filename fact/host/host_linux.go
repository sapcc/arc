package host

import (
	"io/ioutil"
	"os/exec"
	"strings"

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

	if fqdn, domain := fqdn_and_domain(); fqdn != "" {
		facts["fqdn"] = fqdn
		facts["domain"] = domain
	}

	//init system detection
	if contents, err := ioutil.ReadFile("/proc/1/comm"); err == nil && strings.TrimSpace(string(contents)) == "systemd" {
		//upstart and sysvinit report "init" so we can only use this for systemd
		facts["init_package"] = "systemd"
	} else {
		var cmd *exec.Cmd
		switch facts["platform_family"] {
		case "debian":
			cmd = exec.Command("dpkg", "-S", "/sbin/init")
		case "rhel", "fedora", "suse":
			cmd = exec.Command("rpm", "--qf", "%{name}", "-qf", "/sbin/init")
		}
		if cmd != nil {
			if out, err := cmd.Output(); err == nil {
				pkg_name := string(out)
				switch {
				case strings.Contains(pkg_name, "systemd"):
					facts["init_package"] = "systemd"
				case strings.Contains(pkg_name, "upstart"):
					facts["init_package"] = "upstart"
				case strings.Contains(pkg_name, "sysv"):
					facts["init_package"] = "sysv"
				}
			}
		}
	}

	return facts, nil
}
