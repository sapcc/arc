package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/onos/types"
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
	endpoints         types.Endpoints
	clientCa          string
	clientCert        string
	clientKey         string
	configDir         = defaultConfigDir()
	config            types.Config // holds the global confd config.
	transportBackend  string
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
	log.SetLevel(log.DebugLevel)
}

func initConfig() error {
	// Set defaults.
	config = types.Config{
		ConfigDir: defaultConfigDir(),
		Transport: "mqtt",
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

	log.Debug("config dir: ", config.ConfigDir)
	log.Debug("transport: ", config.Transport)
	log.Debug("endpoints: ", config.Endpoints)

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
	}
}
