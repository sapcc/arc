//+build windows

package execute

import "gitHub.***REMOVED***/monsoon/arc/arc"

var scriptSuffix = ".ps1"

func scriptCommand(file string) *arc.Subprocess {
	process := arc.NewSubprocess("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "RemoteSigned", "-Command", "$ErrorActionPreference = 'Stop'; & "+file)
	return process
}
