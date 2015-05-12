package main

import (
	"flag"
	"os"
	"testing"
)

var (
	configFileEnvName = envPrefix + "CONFIG_FILE"
)

func TestConfigFileViaCommandLine(t *testing.T) {
	ResetForTesting(nil)
	os.Setenv(configFileEnvName, "FromEnvShouldBeIgnored")
	os.Args = append(os.Args, "--config-file=FromArg")

	initConfig()

	if configFile != "FromArg" {
		t.Errorf("%s != FromArg", configFile)
	}
}

func TestConfigFileFromEnv(t *testing.T) {
	ResetForTesting(nil)
	os.Setenv(configFileEnvName, "FromEnv")
	initConfig()
	if configFile != "FromEnv" {
		t.Errorf("%s != 'FromEnv'", configFile)
	}
}

func ResetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.Usage = usage
	os.Args = []string{os.Args[0]}
}
