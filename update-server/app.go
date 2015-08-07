package main

import (
	"html/template"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gorilla/handlers"

	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

const appName = "arc-update-server"

var (
	buildsRootPath string
	templates      map[string]*template.Template
	config         arc_config.Config
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
			Name:  "log-level,l",
			Usage: "log level",
			Value: "info",
		},
		cli.StringFlag{
			Name:  "path,p",
			Usage: "Directory containig update artifacts",
		},
		cli.StringFlag{
			Name:  "bind-address,b",
			Usage: "listen address for the update server",
			Value: "0.0.0.0:3000",
		},
	}

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

	app.Run(os.Args)
}

// private

func runServer(c *cli.Context) {
	log.Infof("Starting update server version %s. identity: %s, project: %s, organization: %s", version.Version, config.Identity, config.Project, config.Organization)

	// check mandatory params
	buildsRootPath = c.GlobalString("path")
	if buildsRootPath == "" {
		log.Fatal("No path to update artifacts given.")
	}

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
