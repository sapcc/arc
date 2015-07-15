package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

const (
	appName   = "arc-api-server"
	envPrefix = "ARC_"
)

var config arc_config.Config
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
			Name:   "env",
			Usage:  "environment to use (development, test, production)",
			Value:  "development",
			EnvVar: envPrefix + "ENV",
		},
		cli.StringFlag{
			Name:   "db-config,c",
			Usage:  "database configuration file",
			Value:  "db/dbconf.yml",
			EnvVar: envPrefix + "DB_CONFIG",
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
	log.Infof("Starting api server version %s. Environment: %s", version.Version, c.GlobalString("env"))

	// check endpoint
	if len(config.Endpoints) == 0 {
		log.Fatal("No endpoints for MQTT given")
	}

	var err error
	db, err = ownDb.NewConnection(c.GlobalString("db-config"), c.GlobalString("env"))
	checkErrAndPanic(err, "Error connecting to the DB:")
	defer db.Close()

	// global transport instance
	tp, err = arcNewConnection(config)
	checkErrAndPanic(err, "")
	defer tp.Disconnect()

	// subscribe to all replies
	go arcSubscribeReplies(tp)

	// start the routine scheduler
	go routineScheduler(db, 60*time.Second)

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
