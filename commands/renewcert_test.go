package commands

import (
	"flag"
	"os"
	"testing"

	"github.com/codegangsta/cli"
)

func TestCmdRenewCertUriFromFlag(t *testing.T) {
	// prepare context flags
	flagSet := flag.NewFlagSet("local", 0)
	flagSet.String("api-uri", "https://arc.***REMOVED***", "global")
	ctx := cli.NewContext(nil, flagSet, getParentCtx())

	uri, err := renewCertURI(ctx)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if uri != "https://arc.***REMOVED***/api/v1/agents/renew" {
		t.Error("Expected to get the right renew cert uri")
	}
}

func TestCmdRenewCertUriFromEnv(t *testing.T) {
	// prepare context flags
	flagSet := flag.NewFlagSet("local", 0)
	ctx := cli.NewContext(nil, flagSet, getParentCtx())
	// set env var
	os.Setenv("ARC_UPDATE_URI", "https://beta.arc.***REMOVED***")
	defer func() {
		os.Unsetenv("ARC_UPDATE_URI")
	}()

	uri, err := renewCertURI(ctx)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if uri != "https://arc.***REMOVED***/api/v1/agents/renew" {
		t.Error("Expected to get the right renew cert uri")
	}
}

func TestCmdRenewCertUriFromFlagIgnoreEnv(t *testing.T) {
	// prepare context flags
	flagSet := flag.NewFlagSet("local", 0)
	flagSet.String("api-uri", "https://arc.***REMOVED***", "global")
	ctx := cli.NewContext(nil, flagSet, getParentCtx())
	// set env var
	os.Setenv("ARC_UPDATE_URI", "https://beta.arc.***REMOVED***")
	defer func() {
		os.Unsetenv("ARC_UPDATE_URI")
	}()

	uri, err := renewCertURI(ctx)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if uri != "https://arc.***REMOVED***/api/v1/agents/renew" {
		t.Error("Expected to get the right renew cert uri")
	}
}
