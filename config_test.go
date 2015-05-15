package main

import (
	"os"
	"testing"
)

var (
	configFileEnvName = envPrefix + "CONFIGFILE"
)

func TestConfigFileViaCommandLine(t *testing.T) {
	defer resetEnv()
	os.Setenv(configFileEnvName, "FromEnvShouldBeIgnored")
	os.Args = append(os.Args, "--config-file=FromArg")

	if configFile() != "FromArg" {
		t.Errorf("%s != FromArg", configFile)
	}
}

func TestConfigFileFromEnv(t *testing.T) {
	defer resetEnv()
	os.Setenv(configFileEnvName, "FromEnv")
	if configFile() != "FromEnv" {
		t.Errorf("%s != 'FromEnv'", configFile)
	}
}

func resetEnv() {
	os.Args = []string{os.Args[0]}
	os.Unsetenv(configFileEnvName)
}
