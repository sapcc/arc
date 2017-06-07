package chef

import "fmt"

var (
	chefClientBinary = "/usr/bin/chef-client"
	chefSoloBinary   = "/usr/bin/chef-solo"
	eol              = "\n"
)

func install(installer string) error {
	return fmt.Errorf("Not implemented")
}
