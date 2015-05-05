package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/onos/onos"
)

func defaultConfigDir() string {
	if runtime.GOOS == "windows" {
		return "C:/monsoon"
	} else {
		return "/etc/monsoon"
	}
}

var (
	configFile        = ""
	defaultConfigFile = path.Join(defaultConfigDir(), "onos.cfg")
	printVersion      bool
	endpoints         onos.Endpoints
	clientCa          string
	clientCert        string
	clientKey         string
	configDir         = defaultConfigDir()
	config            onos.Config // holds the global confd config.
	transportBackend  string
	logLevel          = "info"
)

func init() {
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
	flag.StringVar(&transportBackend, "transport", "mqtt", "transport type")
	flag.Var(&endpoints, "endpoint", "transport endpoint url(s)")
	flag.StringVar(&clientCa, "client-ca", "", "client ca cert")
	flag.StringVar(&clientCert, "client-cert", "", "the client cert")
	flag.StringVar(&clientKey, "client-key", "", "the client key")
	flag.StringVar(&configDir, "config-dir", defaultConfigDir(), "the onos conf directory")
	flag.StringVar(&configFile, "config-file", "", "the onos config file")
	flag.StringVar(&logLevel, "log-level", "info", "log level: debug, info, warn, error, fatal")
}

func initConfig() error {
	// Set defaults.
	config = onos.Config{
		ConfigDir: defaultConfigDir(),
		Transport: "mqtt",
		LogLevel:  "info",
	}
	log.Debug(configFile)
	if configFile == "" {
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			configFile = defaultConfigFile
		}
	}
	if configFile == "" {
		log.Debug("Skipping Onos config file.")
	} else {
		log.Debug("Loading " + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		_, err = toml.Decode(string(configBytes), &config)
		if err != nil {
			return err
		}
	}

	// Update config from commandline flags.
	processFlags()

	lvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		return fmt.Errorf("Invalid log level %s", config.LogLevel)
	}
	log.Infof("Setting log level to %d\n", lvl)
	log.Infof("Setting log level to %s\n", logLevel)
	log.SetLevel(lvl)

	return nil
}

// processFlags iterates through each flag set on the command line and
// overrides corresponding configuration settings.
func processFlags() {
	flag.Visit(setConfigFromFlag)
}

func setConfigFromFlag(f *flag.Flag) {
	switch f.Name {
	case "config-dir":
		config.ConfigDir = configDir
	case "client-cert":
		config.ClientCert = clientCert
	case "client-key":
		config.ClientKey = clientKey
	case "client-ca-keys":
		config.ClientCa = clientCa
	case "endpoint":
		config.Endpoints = endpoints
	case "transport":
		config.Transport = transportBackend
	case "log-level":
		config.LogLevel = logLevel
	}
}
