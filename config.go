package main

import (
	"flag"
	"os"
	"path"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/rakyll/globalconf"
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
	appName           = "onos"
	envPrefix         = "ONOS_"
	defaultConfigFile = path.Join(defaultConfigDir(), appName+".cfg")

	printVersion bool
	configFile   string
	endpoint     string

	config onos.Config // holds the global config struct.
)

//We have this for resetting the flags in the tests
func setupCliFlags() {
	config = *new(onos.Config)
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
	flag.StringVar(&configFile, "config-file", defaultConfigFile, "configuration file")
	flag.StringVar(&config.Transport, "transport", "mqtt", "transport backend")
	flag.StringVar(&endpoint, "endpoint", "", "transport endpoint url")
	flag.StringVar(&config.ClientCa, "client-ca", "", "ca certificate")
	flag.StringVar(&config.ClientCert, "client-cert", "", "client certificate")
	flag.StringVar(&config.ClientKey, "client-key", "", "client key")
	flag.StringVar(&config.LogLevel, "log-level", "info", "log level")
}

func initConfig() error {
	setupCliFlags()

	flag.Parse()

	configFileArg := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "config-file" {
			configFileArg = true
		}
	})
	if !configFileArg {
		if cf := os.Getenv(envPrefix + "CONFIG_FILE"); cf != "" {
			configFile = cf
		}
	}
	globalConfFile := configFile
	if _, err := os.Stat(globalConfFile); os.IsNotExist(err) {
		globalConfFile = ""
	}

	conf, err := globalconf.NewWithOptions(&globalconf.Options{
		Filename:  globalConfFile,
		EnvPrefix: envPrefix,
	})
	if err != nil {
		return err
	}
	conf.ParseAll()

	lvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal("Invalid log level ", config.LogLevel)
	}
	log.SetLevel(lvl)

	config.Endpoints = []string{endpoint}
	config.Identity = "me"
	config.Project = "myproject"

	return nil
}
