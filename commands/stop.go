package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"gitHub.***REMOVED***/monsoon/arc/service"
)

func Stop(c *cli.Context) (int, error) {
	err := service.New(c.String("install-dir")).Stop()
	if err != nil {
		fmt.Println(err.Error())
		return 1, err
	}
	return 0, nil
}
