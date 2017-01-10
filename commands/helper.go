package commands

import (
	"flag"

	"github.com/codegangsta/cli"
)

func getParentCtx() *cli.Context {
	flagSet := flag.NewFlagSet("global", 0)
	globalContext := cli.NewContext(nil, flagSet, nil)
	return globalContext
}
