package janitor

import (
	"fmt"
	"runtime"
)

var (
	Version   = "develop"
	GITCOMMIT = "HEAD"
)

func VersionString() string {
	return fmt.Sprintf("%s (%s), %s", Version, GITCOMMIT, runtime.Version())
}
