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
		//config-file is only here for the generated help, it is actually handled above
		//so the the settings from the file are set as env vars before app.Run(...) is called
		cli.StringFlag{
			Name:   "config-file,c",
			Usage:  "Load config file",
			Value:  defaultConfigFile,
			EnvVar: envPrefix + "CONFIGFILE",
		},
		cli.StringFlag{
			Name:   "transport,t",
			Usage:  "Transport backend driver",
			Value:  "mqtt",
			EnvVar: envPrefix + "TRANSPORT",
		},
		cli.StringSliceFlag{
			Name:   "endpoint,e",
			Usage:  "Endpoint url(s) for selected transport",
			EnvVar: envPrefix + "ENDPOINT",
			Value:  new(cli.StringSlice),
		},
		cli.StringFlag{
			Name:   "tls-ca-cert",
			Usage:  "CA to verify transport endpoints",
			EnvVar: envPrefix + "TLS_CA_CERT",
		},
		cli.StringFlag{
			Name:   "tls-client-cert",
			Usage:  "Client cert to use for TLS",
			EnvVar: envPrefix + "TLS_CLIENT_CERT",
		},
		cli.StringFlag{
			Name:   "tls-client-key",
			Usage:  "Private key used in client TLS auth",
			EnvVar: envPrefix + "TLS_CLIENT_KEY",
		},
		cli.StringFlag{
			Name:   "log-level,l",
			Usage:  "Log level",
			EnvVar: envPrefix + "LOG_LEVEL",
			Value:  "info",
		},
		cli.BoolFlag{
			Name:   "no-auto-update",
			Usage:  "Should NO trigger auto updates",
			EnvVar: envPrefix + "NO_AUTO_UPDATE",
		},
		cli.IntFlag{
			Name:   "update-interval",
			Usage:  "Time update interval in seconds",
			EnvVar: envPrefix + "UPDATE_INTERVAL",
			Value:  21600,
		},
		cli.StringFlag{
			Name:   "update-uri",
			Usage:  "Update server uri",
			EnvVar: envPrefix + "UPDATE_URI",
			Value:  "http://localhost:3000/updates",
		},
	}

	app.Commands = commands

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
