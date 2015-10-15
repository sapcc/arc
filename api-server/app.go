package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/databus23/keystone"
	"github.com/databus23/keystone/cache/postgres"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

const (
	appName   = "arc-api-server"
	envPrefix = "ARC_"
)

var (
	config = arc_config.New()
	db     *sql.DB
	tp     transport.Transport
	ks     = keystone.Auth{}
	env    string
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
	app.Usage = "api server"
	app.Action = runServer
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "transport,T",
			Usage:  "Transport backend driver",
			Value:  "mqtt",
			EnvVar: envPrefix + "TRANSPORT",
		},
		cli.StringFlag{
			Name:   "log-level,l",
			Usage:  "Log level",
			EnvVar: envPrefix + "LOG_LEVEL",
			Value:  "info",
		},
		cli.StringSliceFlag{
			Name:   "endpoint,e",
			Usage:  "Endpoint url(s) for selected transport",
			EnvVar: envPrefix + "ENDPOINT",
			Value:  new(cli.StringSlice),
		},
		cli.StringFlag{
			Name:  "bind-address,b",
			Usage: "Update server URL",
			Value: "0.0.0.0:3000",
		},
		cli.StringFlag{
			Name:   "env",
			Usage:  "Environment to use (development, test, production)",
			Value:  "development",
			EnvVar: envPrefix + "ENV",
		},
		cli.StringFlag{
			Name:   "db-config,c",
			Usage:  "Database configuration file",
			Value:  "db/dbconf.yml",
			EnvVar: envPrefix + "DB_CONFIG",
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
			Name:   "keystone-endpoint, ke",
			Usage:  "Endpoint url for Keystone",
			EnvVar: envPrefix + "KEYSTONE_ENDPOINT",
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
	// save the environment
	env = c.GlobalString("env")
	log.Infof("Starting api server version %s. Environment: %s", version.Version, env)

	// check endpoint
	if len(config.Endpoints) == 0 {
		log.Fatal("No endpoints for MQTT given")
	}

	var err error
	db, err = ownDb.NewConnection(c.GlobalString("db-config"), env)
	checkErrAndPanic(err, "Error connecting to the DB:")
	defer db.Close()

	// global transport instance
	tp, err = arcNewConnection(config)
	checkErrAndPanic(err, "")
	defer tp.Disconnect()

	// keystone initialization
	if c.GlobalString("keystone-endpoint") != "" {
		ks.Endpoint = c.GlobalString("keystone-endpoint")
		ks.TokenCache = postgres.New(db, 30*time.Second, "token_cache")
		log.Infof("Keystone binded. Endpoint %q", c.GlobalString("keystone-endpoint"))
	}

	// subscribe to all replies
	go arcSubscribeReplies(tp)

	// start the routine scheduler
	go routineScheduler(db, 60*time.Second)

	// init the router
	router := newRouter(env)

	// run server
	log.Infof("Listening on %q...", c.GlobalString("bind-address"))
	err = http.ListenAndServe(c.GlobalString("bind-address"), router)
	checkErrAndPanic(err, fmt.Sprintf("Failed to bind on %s: ", c.GlobalString("bind-address")))
}

func checkErrAndPanic(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
