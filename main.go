package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	_ "gitHub.***REMOVED***/monsoon/arc/agents/chef"
	_ "gitHub.***REMOVED***/monsoon/arc/agents/execute"
	_ "gitHub.***REMOVED***/monsoon/arc/agents/rpc"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

var exitCode = 0

func main() {

	//If we have a config file load it
	if configFile := configFile(); configFile != "" {
		loadConfigFile(configFile)
	}

	app := cli.NewApp()

	app.Name = appName
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
	app.Usage = "Remote job execution galore"
	app.Version = version.String()
	app.Flags = []cli.Flag{
		optConfigFile,
		optLogLevel,
	}

	app.Commands = cliCommands

	app.Before = func(c *cli.Context) error {
		err := config.Load(c)
		if err != nil {
			log.Fatalf("Invalid configuration: %s\n", err.Error())
			return err
		}

		lvl, err := log.ParseLevel(config.LogLevel)
		if err != nil {
			log.Fatalf("Invalid log level: %s\n", config.LogLevel)
			return err
		}

		log.SetLevel(lvl)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
