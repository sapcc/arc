package host

import "github.com/shirou/gopsutil/host"

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

	return facts, nil
}
