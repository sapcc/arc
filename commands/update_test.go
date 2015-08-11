// +build !integration

package commands

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/inconshreveable/go-update/check"
	"gitHub.***REMOVED***/monsoon/arc/updater"
)

var CheckResult = ""
var responseExample = `{"initiative":"automatically","url":"MIAU://non_valid_url","patch_url":null,"patch_type":null,"version":"999","checksum":null,"signature":null}`

func TestCmdUpdateYes(t *testing.T) {
	// mock apply upload
	origApplyUpdate := updater.ApplyUpdate
	updater.ApplyUpdate = mock_apply_update

	// mock server
	server := getMockServer(200, responseExample)
	defer server.Close()

	// prepare context flags
	ctx := cli.NewContext(nil, getLocalSet(false, false), getGlobalSet(server.URL))

	// mock input to yes
	in, err := mockConfirmStdInput("yes")
	if err != nil {
		t.Error("Expected to not have an error")
	}
	confirmInput = in

	code, err := CmdUpdate(ctx, map[string]interface{}{"appName": "test"})
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if code != 0 {
		t.Error("Expected to have exit code 0")
	}
	if CheckResult == "" {
		t.Error("Expected to apply update")
	}

	defer func() {
		in.Close()
		confirmInput = os.Stdin
		CheckResult = ""
		updater.ApplyUpdate = origApplyUpdate
	}()
}

func TestCmdUpdateNo(t *testing.T) {
	// mock apply upload
	origApplyUpdate := updater.ApplyUpdate
	updater.ApplyUpdate = mock_apply_update

	// mock server
	server := getMockServer(200, responseExample)
	defer server.Close()

	// prepare context flags
	ctx := cli.NewContext(nil, getLocalSet(false, false), getGlobalSet(server.URL))

	// mock input to no
	in, err := mockConfirmStdInput("no")
	if err != nil {
		t.Fatal(err)
	}
	confirmInput = in

	code, err := CmdUpdate(ctx, map[string]interface{}{"appName": "test"})
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if code != 0 {
		t.Error("Expected to have exit code 0")
	}

	if CheckResult != "" {
		t.Error("Expected to not apply update")
	}

	defer func() {
		in.Close()
		confirmInput = os.Stdin
		CheckResult = ""
		updater.ApplyUpdate = origApplyUpdate
	}()
}

func TestCmdUpdateForce(t *testing.T) {
	// mock apply upload
	origApplyUpdate := updater.ApplyUpdate
	updater.ApplyUpdate = mock_apply_update

	// mock server
	server := getMockServer(200, responseExample)
	defer server.Close()

	// prepare context flags
	ctx := cli.NewContext(nil, getLocalSet(false, true), getGlobalSet(server.URL))

	code, err := CmdUpdate(ctx, map[string]interface{}{"appName": "test"})
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if code != 0 {
		t.Error("Expected to have exit code 0")
	}

	if CheckResult == "" {
		t.Error("Expected to apply update")
	}

	defer func() {
		CheckResult = ""
		updater.ApplyUpdate = origApplyUpdate
	}()
}

func TestCmdUpdateNoUpdate(t *testing.T) {
	// mock apply upload
	origApplyUpdate := updater.ApplyUpdate
	updater.ApplyUpdate = mock_apply_update

	// mock server
	server := getMockServer(200, responseExample)
	defer server.Close()

	// prepare context flags
	ctx := cli.NewContext(nil, getLocalSet(true, false), getGlobalSet(server.URL))

	code, err := CmdUpdate(ctx, map[string]interface{}{"appName": "test"})
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if code != 0 {
		t.Error("Expected to have exit code 0")
	}

	if CheckResult != "" {
		t.Error("Expected to not apply update")
	}

	defer func() {
		CheckResult = ""
		updater.ApplyUpdate = origApplyUpdate
	}()
}

// private

func getGlobalSet(serverUrl string) *flag.FlagSet {
	flagSet := flag.NewFlagSet("global", 0)
	flagSet.String("update-uri", serverUrl, "global")
	return flagSet
}

func getLocalSet(noUpdate bool, force bool) *flag.FlagSet {
	flagSet := flag.NewFlagSet("local", 0)
	flagSet.Bool("no-update", noUpdate, "local")
	flagSet.Bool("force", force, "local")
	return flagSet
}

func mockConfirmStdInput(input string) (*os.File, error) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	_, err = io.WriteString(in, input)
	if err != nil {
		return nil, err
	}
	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	return in, nil
}

func mock_apply_update(r *check.Result) error {
	CheckResult = "mock apply_update"
	return nil
}

func getMockServer(code int, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	return server
}
