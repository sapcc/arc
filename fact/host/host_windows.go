package host

import (
	"os"
	"regexp"
	"strings"
	"syscall"

	_ "net" //we need this to ensure the winsock subsystem is initialized

	"github.com/StackExchange/wmi"
	gopsutil "github.com/shirou/gopsutil/common"
)

type Win32_OperatingSystem struct {
	Version string

	//yields: wmi: cannot load field "CSDVersion" into a "string": unsupported type (<nil>) on windows 2012
	//CSDVersion  string

	OSType      uint16
	Caption     string
	BuildNumber string
}

func (h Source) Facts() (map[string]interface{}, error) {
	facts := make(map[string]interface{})
	facts["hostname"] = nil
	facts["fqdn"] = nil
	facts["domain"] = nil
	facts["platform"] = "windows"
	facts["platform_family"] = "windows"
	facts["platform_version"] = nil

	if hostname, err := os.Hostname(); err == nil {
		facts["hostname"] = strings.ToLower(hostname)
		if hostent, err := syscall.GetHostByName(hostname); err == nil {
			fqdn := strings.ToLower(gopsutil.BytePtrToString(hostent.Name))
			facts["fqdn"] = fqdn
			domain_regexp := regexp.MustCompile(`.*?\.(.+)$`)
			if m := domain_regexp.FindStringSubmatch(fqdn); m != nil {
				facts["domain"] = m[1]
			}
		}
	}

	var win32_os []Win32_OperatingSystem
	q := wmi.CreateQuery(&win32_os, "")
	if wmi.Query(q, &win32_os) == nil {
		facts["platform_version"] = win32_os[0].Version
	}

	return facts, nil
}
