package version

//needs to be a var (no const)
//so that we van overwrite during linking with -X main.GITCOMMIT ...
var (
	Version   = "20150910.01"
	GITCOMMIT = "HEAD"
)

func String() string {
	return Version + " (" + GITCOMMIT + ")"
}
