package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/sapcc/arc/service"
)

func Restart(c *cli.Context) (int, error) {
	err := service.New(c.String("install-dir")).Restart()
	if err != nil {
		fmt.Println(err.Error())
		return 1, err
	}
	return 0, nil
}
