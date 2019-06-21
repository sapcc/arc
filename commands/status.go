package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/sapcc/arc/service"
)

func Status(c *cli.Context) (int, error) {
	state, message, err := service.New(c.String("install-dir")).Status()
	fmt.Println(message)
	switch state {
	case service.RUNNING:
		return 0, err
	case service.STOPPED:
		return 3, err
	}
	return 4, err
}
