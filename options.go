package main

import (
	"github.com/codegangsta/cli"
)

var optConfigFile = cli.StringFlag{
	Name:   "config-file,c",
	Usage:  "load config file",
	Value:  defaultConfigFile,
	EnvVar: envPrefix + "CONFIGFILE",
}

var optTransport = cli.StringFlag{
	Name:   "transport,t",
	Usage:  "transport backend driver",
	Value:  "mqtt",
	EnvVar: envPrefix + "TRANSPORT",
}

var optEndpoint = cli.StringSliceFlag{
	Name:   "endpoint,e",
	Usage:  "endpoint url(s) for selected transport",
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
	Usage:  "log level",
	EnvVar: envPrefix + "LOG_LEVEL",
	Value:  "info",
}

var optNoAutoUpdate = cli.BoolFlag{
	Name:   "no-auto-update",
	Usage:  "Specifies if the server should NO trigger auto updates",
	EnvVar: envPrefix + "NO_AUTO_UPDATE",
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
	Value:  "http://localhost:3000/updates",
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
	Usage: "installation directory",
	Value: defaultConfigDir(),
}