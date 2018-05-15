package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

//lint:file-ignore U1000 Ignore all unused code, it is just for testing

const appName = "arc-update-server"

var (
	templates map[string]*template.Template
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
	app.Action = runServer
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
	app.Before = func(c *cli.Context) error {
		lvl, err := log.ParseLevel(c.GlobalString("log-level"))
		if err != nil {
			log.Fatalf("Invalid log level: %s\n", c.GlobalString("log-level"))
			return err
		}
		log.SetLevel(lvl)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func runServer(c *cli.Context) {
	log.Infof("Starting update server version %s.", version.String())

	// get the router
	router := NewRouter()

	// run server
	log.Infof("Listening on %q...", c.GlobalString("bind-address"))
	err := http.ListenAndServe(c.GlobalString("bind-address"), router)
	checkErrAndPanic(err, fmt.Sprintf("Failed to bind on %s: ", c.GlobalString("bind-address")))
}

func checkErrAndPanic(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
