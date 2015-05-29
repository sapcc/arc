package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"

	"gitHub.***REMOVED***/monsoon/arc/arc"
)

func defaultConfigDir() string {
	if runtime.GOOS == "windows" {
		return "C:/monsoon"
	} else {
		return "/etc/monsoon"
	}
}

var (
	appName           = "arc"
	envPrefix         = "ARC_"
	defaultConfigFile = path.Join(defaultConfigDir(), appName+".cfg")
	config            arc.Config
)

//returns the path to the config file we want to load
//returns the file the user explicitly specified by flag or env var
//alternativly it returns the default config file if it exists
func configFile() string {
	env := os.Getenv(envPrefix + "CONFIGFILE")
	var filename string
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	fs.StringVar(&filename, "config-file", env, "")
	fs.StringVar(&filename, "c", env, "")
	fs.Parse(os.Args[1:])
	//No file specified by the user
	if filename == "" {
		if _, err := os.Stat(defaultConfigFile); err == nil {
			return defaultConfigFile
		}
	}
	return filename
}

func loadConfigFile(file string) error {
	vars, err := godotenv.Read(file)
	if err != nil {
		return err
	}
	log.Debug("Loaded config file: ", file)
	for name, value := range vars {
		name = strings.Replace(name, "-", "_", -1)
		name = strings.Replace(name, " ", "_", -1)
		name = strings.ToUpper(name)
		if !strings.HasPrefix(name, envPrefix) {
			name = envPrefix + name
		}
		if os.Getenv(name) == "" {
			os.Setenv(name, value)
		}
	}
	return nil
}
