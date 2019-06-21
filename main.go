package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	_ "github.com/sapcc/arc/agents/chef"
	_ "github.com/sapcc/arc/agents/execute"
	_ "github.com/sapcc/arc/agents/rpc"
	"github.com/sapcc/arc/utils/cliDescriptionGenerator/templates"
	"github.com/sapcc/arc/version"
)

var exitCode = 0

func main() {

	//If we have a config file load it
	if configFile := configFile(); configFile != "" {
		if err := loadConfigFile(configFile); err != nil {
			log.Fatalf("invalid config file: %s\n", err)
		}
	}

	// override cli template
	cli.AppHelpTemplate = templates.AppTemplate

	app := cli.NewApp()

	app.Name = appName
	app.Usage = fmt.Sprint(cmdUsage["docs-commands"], "\n\n", cmdDescription["docs-commands"])
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
