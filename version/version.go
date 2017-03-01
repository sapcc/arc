package version

import (
	"fmt"
	"runtime"
)

//needs to be a var (no const)
//so that we van overwrite during linking with -X main.GITCOMMIT ...
var (
	Version   = ""
	GITCOMMIT = "HEAD"
)

func String() string {
	return fmt.Sprintf("%s (%s), %s", Version, GITCOMMIT, runtime.Version())
}
