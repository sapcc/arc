package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gitHub.***REMOVED***/monsoon/arc/version"
	"net/http"
	"os"
)

const appName = "arc-api-server"

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
	app.Usage = "api server"
	app.Action = runServer
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level,l",
			Usage: "log level",
			Value: "info",
		},
		cli.StringFlag{
			Name:  "bind-address,b",
			Usage: "listen address for the update server",
			Value: "0.0.0.0:3000",
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
	// init the router
	router := newRouter()

	// run server
	log.Infof("Listening on %q...", c.GlobalString("bind-address"))
	if err := http.ListenAndServe(c.GlobalString("bind-address"), accessLogger(router)); err != nil {
		log.Fatalf("Failed to bind on %s: %s", c.GlobalString("bind-address"), err)
	}
}

func accessLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}