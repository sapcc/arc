package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"gitHub.***REMOVED***/monsoon/arc/service"
)

func Status(c *cli.Context) (int, error) {
	status, err := service.New("/opt/arc").Status()
	fmt.Println(status)
	if err != nil {
		return 1, err
	}
	return 0, nil
}
