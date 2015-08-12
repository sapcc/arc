package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"gitHub.***REMOVED***/monsoon/arc/service"
)

func Status(c *cli.Context) (int, error) {
	status, err := service.Status("/opt/arc/service")
	if err != nil {
		return 1, err
	}
	fmt.Println(status)
	return 0, nil
}
