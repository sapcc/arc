package main

import (
	"html/template"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gorilla/handlers"

	"gitHub.***REMOVED***/monsoon/arc/version"
)

const appName = "arc-update-server"

var (
	buildsRootPath string
	templates      map[string]*template.Template
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
			Name:   "path,p",
			Usage:  "Directory containig update artifacts",
			EnvVar: "ARTIFACTS_PATH",
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

	app.Run(os.Args)
}

// private

func runServer(c *cli.Context) {

	// check mandatory params
	buildsRootPath = c.GlobalString("path")
	if buildsRootPath == "" {
		log.Fatal("No path to update artifacts given.")
	}
	log.Infof("Starting update server version %s.", version.String())
	log.Infof("Serving artificats from %s.", buildsRootPath)

	// cache the templates
	templates = getTemplates()

	// get the router
	router := newRouter()

	// run server
	log.Infof("Listening on %q...", c.GlobalString("bind-address"))
	accessLogger := handlers.CombinedLoggingHandler(os.Stdout, router)
	if err := http.ListenAndServe(c.GlobalString("bind-address"), accessLogger); err != nil {
		log.Fatalf("Failed to bind on %s: %s", c.GlobalString("bind-address"), err)
	}
}
