package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/cfssl/csr"
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
	config           = arc_config.New()
	db               *sql.DB
	tp               transport.Transport
	ks               = keystone.Auth{}
	env              string
	pkiEnabled       = false
	agentUpdateURL   = "UPDATE_URL_NOT_CONFIGURED"
	agentEndpointURL = "ENDPOINT_URL_NOT_CONFIGURED"
	agentApiURL      = "" // Should not be initialized because not all regions have TLS termination directly on the api server
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
			Name:   "bind-address",
			Usage:  "Listen address for http api server",
			Value:  "0.0.0.0:3000",
			EnvVar: envPrefix + "LISTEN",
		},
		cli.StringFlag{
			Name:   "bind-address-tls",
			Usage:  "Listen address for https api server",
			Value:  "0.0.0.0:443",
			EnvVar: envPrefix + "LISTEN_TLS",
		},
		cli.StringFlag{
			Name:   "tls-server-cert",
			Usage:  "Server cert to use for TLS",
			EnvVar: envPrefix + "TLS_SERVER_CERT",
		},
		cli.StringFlag{
			Name:   "tls-server-key",
			Usage:  "Private key used in server TLS",
			EnvVar: envPrefix + "TLS_SERVER_KEY",
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
			Usage:  "CA to verify transport endpoints in the messaging middleware",
			EnvVar: envPrefix + "TLS_CA_CERT",
		},
		cli.StringFlag{
			Name:   "tls-client-cert",
			Usage:  "Client cert to use for TLS in the messaging middleware",
			EnvVar: envPrefix + "TLS_CLIENT_CERT",
		},
		cli.StringFlag{
			Name:   "tls-client-key",
			Usage:  "Private key used in client TLS auth in the messaging middleware",
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
			Name:   "pki-ca-cert",
			Usage:  "PKI CA certfiicate used to sign the new certificate",
			EnvVar: envPrefix + "PKI_CA_CERT",
		},
		cli.StringFlag{
			Name:   "pki-ca-key",
			Usage:  "PKI CA private key",
			EnvVar: envPrefix + "PKI_CA_KEY",
		},
		cli.StringFlag{
			Name:   "agent-update-url",
			Usage:  "The default update URL for agents. Only used for agent install script.",
			EnvVar: envPrefix + "AGENT_UPDATE_URL",
		},
		cli.StringFlag{
			Name:   "agent-endpoint-url",
			Usage:  "The default endpoint URL for agents. Only used for agent install script.",
			EnvVar: envPrefix + "AGENT_ENDPOINT_URL",
		},
		cli.StringFlag{
			Name:   "agent-api-url",
			Usage:  "API URL. Only used for renew cert.",
			EnvVar: envPrefix + "AGENT_API_URL",
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
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

	if c.GlobalString("agent-update-url") != "" {
		agentUpdateURL = c.GlobalString("agent-update-url")
	}

	if c.GlobalString("agent-update-url") != "" {
		agentEndpointURL = c.GlobalString("agent-endpoint-url")
	} else if len(config.Endpoints) > 0 {
		agentEndpointURL = config.Endpoints[0]
	}

	if c.GlobalString("agent-api-url") != "" {
		agentApiURL = c.GlobalString("agent-api-url")
	}

	// create db connection
	var err error
	db, err = ownDb.NewConnection(c.GlobalString("db-config"), env)
	FatalfOnError(err, "Error connecting to the DB: %s", err)
	defer db.Close()

	if c.GlobalString("pki-ca-cert") != "" {
		err = pki.SetupSigner(c.GlobalString("pki-ca-cert"), c.GlobalString("pki-ca-key"), c.GlobalString("pki-config"))
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
			csrBytes, clientKey, err := pki.CreateSignReqCertAndPrivKey(cn, "arc-api", "arc-api")
			FatalfOnError(err, "Failed to create CSR: %s", err)
			clientCert, err := pki.Sign(csrBytes, signer.Subject{CN: cn, Names: []csr.Name{csr.Name{O: "arc-api", OU: "arc-api"}}}, "default")
			FatalfOnError(err, "Failed to sign ephemeral certificate: %s", err)
			tlsCert, err := tls.X509KeyPair(clientCert, clientKey)
			FatalfOnError(err, "Failed to use generated certificate: %s", err)
			config.ClientCert = &tlsCert

			caCert, err := ioutil.ReadFile(c.GlobalString("pki-ca-cert"))
			FatalfOnError(err, "Failed to read path %#v: %s", c.GlobalString("pki-ca-cert"), err)
			config.CACerts = x509.NewCertPool()
			if !config.CACerts.AppendCertsFromPEM(caCert) {
				log.Fatalf("Failed to load CA from %#v. Not PEM encoded?", c.GlobalString("pki-ca-cert"))
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

	// init the router
	router := newRouter(env)

	// run server
	server := NewSever(c.GlobalString("tls-server-cert"), c.GlobalString("tls-server-key"), c.GlobalString("pki-ca-cert"), c.GlobalString("bind-address"), c.GlobalString("bind-address-tls"), router)
	go server.run()

	// catch gracefull shutdown and shutdown to close the connetions
	gracefulChan := make(chan os.Signal, 1)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(gracefulChan, syscall.SIGTERM)
	signal.Notify(shutdownChan, syscall.SIGINT)
	for {
		select {
		case s := <-shutdownChan:
			log.Infof("Captured %v", s)
			server.close()
		case s := <-gracefulChan:
			log.Infof("Captured %v", s)
			if err = server.shutdown(); err != nil {
				log.Errorf("error shsuting down server %s\n", err)
			}
		}
	}
}

// FatalfOnError fatal on error
func FatalfOnError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalf(msg, args...)
	}
}
