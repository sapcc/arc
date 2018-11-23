package main

import (
	"github.com/codegangsta/cli"
)

var optConfigFile = cli.StringFlag{
	Name:   "config-file,c",
	Usage:  "Load config file",
	Value:  defaultConfigFile,
	EnvVar: envPrefix + "CONFIGFILE",
}

var optTransport = cli.StringFlag{
	Name:   "transport,T",
	Usage:  "Transport backend driver",
	Value:  "mqtt",
	EnvVar: envPrefix + "TRANSPORT",
}

var optEndpoint = cli.StringSliceFlag{
	Name:   "endpoint,e",
	Usage:  "Endpoint url(s) for selected transport",
	EnvVar: envPrefix + "ENDPOINT",
	Value:  new(cli.StringSlice),
}

var optTlsCaCert = cli.StringFlag{
	Name:   "tls-ca-cert",
	Usage:  "CA to verify transport endpoints",
	EnvVar: envPrefix + "TLS_CA_CERT",
}

var optTlsClientCert = cli.StringFlag{
	Name:   "tls-client-cert",
	Usage:  "Client cert to use for TLS",
	EnvVar: envPrefix + "TLS_CLIENT_CERT",
}

var optTlsClientKey = cli.StringFlag{
	Name:   "tls-client-key",
	Usage:  "Private key used in client TLS auth",
	EnvVar: envPrefix + "TLS_CLIENT_KEY",
}

var optLogLevel = cli.StringFlag{
	Name:   "log-level,l",
	Usage:  "Log level",
	EnvVar: envPrefix + "LOG_LEVEL",
	Value:  "info",
}

var optUpdateInterval = cli.IntFlag{
	Name:   "update-interval",
	Usage:  "Time update interval in seconds",
	EnvVar: envPrefix + "UPDATE_INTERVAL",
	Value:  21600,
}

var optUpdateUri = cli.StringFlag{
	Name:   "update-uri",
	Usage:  "Update server uri",
	EnvVar: envPrefix + "UPDATE_URI",
}

var optCertUpdateInterval = cli.IntFlag{
	Name:   "cert-update-interval",
	Usage:  "Time cert update interval in minutes",
	EnvVar: envPrefix + "CERT_UPDATE_INTERVAL",
	Value:  1440,
}

var optCertUpdateThreshold = cli.IntFlag{
	Name:   "cert-update-threshold",
	Usage:  "Hours threshold before updating cert",
	EnvVar: envPrefix + "CERT_UPDATE_THRESHOLD",
	Value:  744,
}

var optApiUri = cli.StringFlag{
	Name:   "api-uri",
	Usage:  "Api uri",
	EnvVar: envPrefix + "API_URI",
}

var optTimeout = cli.IntFlag{
	Name:  "timeout, t",
	Usage: "Timeout for executing the action",
	Value: 60,
}

var optIdentity = cli.StringFlag{
	Name:  "identity, i",
	Usage: "Target system",
	Value: "",
}

var optPayload = cli.StringFlag{
	Name:  "payload,p",
	Usage: "Payload for action",
	Value: "",
}

var optStdin = cli.BoolFlag{
	Name:  "stdin,s",
	Usage: "Read payload from stdin",
}

var optForce = cli.BoolFlag{
	Name:  "force,f",
	Usage: "No confirmation is needed",
}

var optNoUpdate = cli.BoolFlag{
	Name:  "no-update,n",
	Usage: "No update is triggered",
}

var optRegistrationUrl = cli.StringFlag{
	Name:  "registration-url,r",
	Usage: "Registration url",
}

var optInstallDir = cli.StringFlag{
	Name:  "install-dir,i",
	Usage: "Installation directory",
	Value: defaultConfigDir(),
}

var optCommonName = cli.StringFlag{
	Name:  "common-name,cn",
	Usage: "The name of the arc agent",
}
