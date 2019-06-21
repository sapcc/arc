package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/sapcc/arc/service"
)

func Start(c *cli.Context) (int, error) {
	err := service.New(c.String("install-dir")).Start()
	if err != nil {
		fmt.Println(err.Error())
		return 1, err
	}
	return 0, nil
}
