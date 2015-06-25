package main

import (
	"os"
	"runtime"

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
			Usage:  "load config file",
			Value:  defaultConfigFile,
			EnvVar: envPrefix + "CONFIGFILE",
		},
		cli.StringFlag{
			Name:   "transport,t",
			Usage:  "transport backend driver",
			Value:  "mqtt",
			EnvVar: envPrefix + "TRANSPORT",
		},
		cli.StringSliceFlag{
			Name:   "endpoint,e",
			Usage:  "endpoint url(s) for selected transport",
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
			Usage:  "log level",
			EnvVar: envPrefix + "LOG_LEVEL",
			Value:  "info",
		},
		cli.BoolFlag{
			Name:   "no-auto-update",
			Usage:  "Specifies if the server should NO trigger auto updates",
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
		config.Endpoints = c.GlobalStringSlice("endpoint")
		config.Transport = c.GlobalString("transport")
		if c.GlobalString("tls-client-cert") != "" || c.GlobalString("tls-client-key") != "" || c.GlobalString("tls-ca-cert") != "" {
			if err := config.LoadTLSConfig(c.GlobalString("tls-client-cert"), c.GlobalString("tls-client-key"), c.GlobalString("tls-ca-cert")); err != nil {
				return err
			}
		} else {
			//This is only for testing when running without a tls certificate
			config.Identity = runtime.GOOS
			config.Project = "test-project"
			config.Organization = "test-org"
		}
		lvl, err := log.ParseLevel(c.GlobalString("log-level"))
		if err != nil {
			log.Fatalf("Invalid log level: %s\n", c.GlobalString("log-level"))
			return err
		}
		config.LogLevel = c.GlobalString("log-level")
		log.SetLevel(lvl)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
