package main

import (
	"os"

	"gitHub.***REMOVED***/monsoon/arc/janitor"

	"github.com/codegangsta/cli"
)

const (
	appName   = "Arc janitor"
	envPrefix = "ARC_"
)

func main() {
	app := cli.NewApp()

	app.Name = appName
	app.Version = janitor.VersionString()
	app.Authors = []cli.Author{
		{
			Name:  "Arturo Reuschenbach Puncernau",
			Email: "a.reuschenbach.puncernau@sap.com",
		},
		{
			Name:  "Fabian Ruff",
			Email: "fabian.ruff@sap.com",
		},
	}
	app.Usage = "Arc clean jobs scheduler"
	app.Action = runJanitor
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-config,c",
			Usage:  "Database configuration file",
			Value:  "db/dbconf.yml",
			EnvVar: envPrefix + "DB_CONFIG",
		},
		cli.StringFlag{
			Name:   "env",
			Usage:  "Environment to use (development, test, production)",
			Value:  "development",
			EnvVar: envPrefix + "ENV",
		},
		cli.StringFlag{
			Name:   "bind-address,b",
			Usage:  "Listen address for live cron jobs monitoring and metrics HTTP endpoint",
			Value:  "0.0.0.0:3000",
			EnvVar: "BIND_ADDRESS",
		},
	}

	app.Run(os.Args)
}

func runJanitor(c *cli.Context) {
	// init janitor
	conf := janitor.JanitorConf{
		BindAddress:  c.GlobalString("bind-address"),
		DbConfigFile: c.GlobalString("db-config"),
		Environment:  c.GlobalString("env"),
	}
	janitor := janitor.InitJanitor(conf)

	// start scheduler
	janitor.InitScheduler()

	// start http server
	janitor.InitServer()
}
