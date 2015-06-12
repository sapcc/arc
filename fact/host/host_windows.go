package host

import (
	"os"
	"strings"
	"syscall"

	_ "net" //we need this to ensure the winsock subsystem is initialized

	gopsutil "github.com/shirou/gopsutil/common"
	"gitHub.***REMOVED***/monsoon/arc/vendor/github.com/StackExchange/wmi"
)

type Win32_OperatingSystem struct {
	Version     string
	CSDVersion  string
	OSType      uint16
	Caption     string
	BuildNumber string
}

func (h Source) Facts() (map[string]interface{}, error) {
	facts := make(map[string]interface{})
	if hostname, err := os.Hostname(); err == nil {
		facts["hostname"] = strings.ToLower(hostname)
		if hostent, err := syscall.GetHostByName(hostname); err == nil {
			facts["fqdn"] = strings.ToLower(gopsutil.BytePtrToString(hostent.Name))
		}
	}

	facts["platform"] = "windows"
	facts["platform_family"] = "windows"
	var win32_os []Win32_OperatingSystem
	q := wmi.CreateQuery(&win32_os, "")
	if wmi.Query(q, &win32_os) == nil {
		facts["platform_version"] = win32_os[0].Version
	}

	return facts, nil
}
