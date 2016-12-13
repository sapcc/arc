package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	cffsl_cli "github.com/cloudflare/cfssl/cli"
	cfssl_config "github.com/cloudflare/cfssl/config"
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
	config    = arc_config.New()
	db        *sql.DB
	tp        transport.Transport
	ks        = keystone.Auth{}
	env       string
	pkiConfig cffsl_cli.Config
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
			Name:   "bind-address,b",
			Usage:  "Listen address for api server",
			Value:  "0.0.0.0:3000",
			EnvVar: envPrefix + "LISTEN",
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
		cli.StringFlag{
			Name:   "pki-profile-config",
			Usage:  "Path to PKI profile configuration file",
			Value:  "etc/pki_default_conf.json",
			EnvVar: envPrefix + "PKI_CONFIG",
		},
		cli.StringFlag{
			Name:   "pki-ca",
			Usage:  "PKI CA used to sign the new certificate",
			EnvVar: envPrefix + "PKI_CA",
		},
		cli.StringFlag{
			Name:   "pki-ca-key",
			Usage:  "PKI CA private key",
			EnvVar: envPrefix + "PKI_CA_KEY",
		},
	}

	app.Before = func(c *cli.Context) error {
		// load app configuraion
		err := config.Load(c)
		if err != nil {
			log.Fatalf("Invalid configuration: %s\n", err.Error())
			return err
		}

		// set log level
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

	// create db connection
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

	// load pki configuration
	if c.GlobalString("pki-ca") != "" {
		pkiConfig.CAFile = c.GlobalString("pki-ca")
	}
	if c.GlobalString("pki-ca-key") != "" {
		pkiConfig.CAKeyFile = c.GlobalString("pki-ca-key")
	}
	if c.GlobalString("pki-profile-config") != "" {
		pkiConfig.ConfigFile = c.GlobalString("pki-profile-config")
	}
	pkiConfig.CFG, err = cfssl_config.LoadFile(pkiConfig.ConfigFile)
	checkErrAndPanic(err, fmt.Sprintf("Failed to load PKI profile config file: %s", err))

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
