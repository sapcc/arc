package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

const (
	appName   = "arc-api-server"
	envPrefix = "ARC_"
)

var config arc.Config
var db *sql.DB
var tp transport.Transport

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
			Name:   "transport,t",
			Usage:  "transport backend driver",
			Value:  "mqtt",
			EnvVar: envPrefix + "TRANSPORT",
		},
		cli.StringFlag{
			Name:   "log-level,l",
			Usage:  "log level",
			EnvVar: envPrefix + "LOG_LEVEL",
			Value:  "info",
		},
		cli.StringSliceFlag{
			Name:   "endpoint,e",
			Usage:  "endpoint url(s) for selected transport",
			EnvVar: envPrefix + "ENDPOINT",
			Value:  new(cli.StringSlice),
		},
		cli.StringFlag{
			Name:  "bind-address,b",
			Usage: "listen address for the update server",
			Value: "0.0.0.0:3000",
		},
		cli.StringFlag{
			Name:  "db-bind-address,db",
			Usage: "db connection address",
			Value: "postgres://arc:arc@localhost:5432/arc_dev?sslmode=disable",
		},
	}

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

	app.Run(os.Args)
}

// private

func runServer(c *cli.Context) {
	log.Infof("Starting api server version %s. identity: %s, project: %s, organization: %s", version.Version, config.Identity, config.Project, config.Organization)

	// check endpoint
	if len(config.Endpoints) == 0 {
		log.Fatal("No endpoints for MQTT given")
	}

	var err error

	// db global instance
	db, err = ownDb.NewConnection(c.GlobalString("db-bind-address"))
	checkErrAndPanic(err, "Error connecting to the DB or creating tables:")
	defer db.Close()

	// global transport instance
	tp, err = arcNewConnection(config)
	checkErrAndPanic(err, "")
	defer tp.Disconnect()

	// subscribe to all replies
	go arcSubscribeReplies(tp)

	// start the routine scheduler
	go routineScheduler(db)

	// init the router
	router := newRouter()

	// run server
	log.Infof("Listening on %q...", c.GlobalString("bind-address"))
	err = http.ListenAndServe(c.GlobalString("bind-address"), accessLogger(router))
	checkErrAndPanic(err, fmt.Sprintf("Failed to bind on %s: ", c.GlobalString("bind-address")))
}

func accessLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func checkErrAndPanic(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
