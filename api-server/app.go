package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/cfssl/signer"
	"github.com/codegangsta/cli"
	"github.com/databus23/keystone"
	"github.com/databus23/keystone/cache/postgres"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pki"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

const (
	appName   = "arc-api-server"
	envPrefix = "ARC_"
)

var (
	config     = arc_config.New()
	db         *sql.DB
	tp         transport.Transport
	ks         = keystone.Auth{}
	env        string
	pkiEnabled = false
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
			Name:   "pki-config",
			Usage:  "Path to PKI profile configuration file",
			Value:  "etc/pki.json",
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
	FatalfOnError(err, "Error connecting to the DB: %s", err)
	defer db.Close()

	if c.GlobalString("pki-ca") != "" {
		err = pki.SetupSigner(c.GlobalString("pki-ca"), c.GlobalString("pki-ca-key"), c.GlobalString("pki-config"))
		FatalfOnError(err, "Failed to initialize PKI subsystem: %s", err)
		pkiEnabled = true
		// dynamically generate transport certificate if CA is given but client certificate is missing
		if config.ClientCert == nil {
			cn := os.Getenv("COMMON_NAME")
			if cn == "" {
				if cn, err = os.Hostname(); err != nil {
					log.Fatalf("Couldn't determine hostname: %s", err)
				}
			}
			log.Infof("Generating ephemeral client certificate for identity %#v", cn)
			csr, clientKey, err := pki.CreateCSR(cn, "", "")
			FatalfOnError(err, "Failed to create CSR: %s", err)
			clientCert, err := pki.Sign(csr, signer.Subject{CN: cn}, "default")
			FatalfOnError(err, "Failed to sign ephemeral certificate: %s", err)
			tlsCert, err := tls.X509KeyPair(clientCert, clientKey)
			FatalfOnError(err, "Failed to use generated certificate: %s", err)
			config.ClientCert = &tlsCert

			caCert, err := ioutil.ReadFile(c.GlobalString("pki-ca"))
			FatalfOnError(err, "Failed to read path %#v: %s", c.GlobalString("pki-ca"), err)
			config.CACerts = x509.NewCertPool()
			if !config.CACerts.AppendCertsFromPEM(caCert) {
				log.Fatalf("Failed to load CA from %#v. Not PEM encoded?", c.GlobalString("pki-ca"))
			}
		}
	}

	// global transport instance
	tp, err = arcNewConnection(config)
	if err != nil {
		log.Fatal(err)
	}
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
	FatalfOnError(err, "Failed to bind on %s: ", c.GlobalString("bind-address"))
}

func FatalfOnError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalf(msg, args...)
	}
}
