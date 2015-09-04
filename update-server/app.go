package main

import (
	"html/template"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gitHub.***REMOVED***/monsoon/arc/version"
	"gitHub.***REMOVED***/monsoon/arc/update-server/storage"	
)

const appName = "arc-update-server"

var (
	st    				storage.Storage
	templates     map[string]*template.Template
)

func main() {
	app := cli.NewApp()

	app.Name = appName
	app.Version = version.String()
	app.Authors = []cli.Author{
		{
			Name:  "Fabian Ruff",
			Email: "fabian.ruff@sap.com",
		},
		{
			Name:  "Arturo Reuschenbach Puncernau",
			Email: "a.reuschenbach.puncernau@sap.com",
		},
	}
	app.Usage = "web server to to check and deliver from updates"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log-level,l",
			Usage:  "log level",
			Value:  "info",
			EnvVar: "LOG_LEVEL",
		},
		cli.StringFlag{
			Name:   "bind-address,b",
			Usage:  "listen address for the update server",
			Value:  "0.0.0.0:3000",
			EnvVar: "BIND_ADDRESS",
		},
	}

	app.Commands = cliCommands

	app.Before = func(c *cli.Context) error {
		lvl, err := log.ParseLevel(c.GlobalString("log-level"))
		if err != nil {
			log.Fatalf("Invalid log level: %s\n", c.GlobalString("log-level"))
			return err
		}
		log.SetLevel(lvl)
		return nil
	}

	app.Run(os.Args)
}
