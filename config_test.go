// +build !integration

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
	os.Args = []string{os.Args[0], "-config-file=FromArg"}

	c := configFile()

	if c != "FromArg" {
		t.Errorf("%s != FromArg", c)
	}
}

func TestConfigFileFromEnv(t *testing.T) {
	defer resetEnv()
	os.Setenv(configFileEnvName, "FromEnv")
	c := configFile()
	if configFile() != "FromEnv" {
		t.Errorf("%s != 'FromEnv'", c)
	}
}

func resetEnv() {
	os.Args = []string{os.Args[0]}
	os.Unsetenv(configFileEnvName)
}
