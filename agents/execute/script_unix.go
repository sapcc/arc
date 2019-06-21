//+build darwin linux

package execute

import "github.com/sapcc/arc/arc"

var scriptSuffix = ".sh"

func scriptCommand(file string) *arc.Subprocess {
	process := arc.NewSubprocess("/bin/bash", file)
	return process
}
