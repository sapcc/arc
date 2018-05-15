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
	"gitHub.***REMOVED***/monsoon/arc/updater"
	arcVersion "gitHub.***REMOVED***/monsoon/arc/version"
)

var CheckResult = ""
var responseExample = `{
  "app_id": "arc",
  "os": "darwin",
  "arch": "amd64",
  "checksum": "af938b7e52da57de81ce0daf5bc196b703dfe1b54a4b5a6c3a6f8738024eb84d",
  "version": "20170316.01",
  "url":"arc_20170316.01_darwin_amd64"
}`

func TestCmdUpdateMissingUri(t *testing.T) {
	// prepare context flags
	ctx := cli.NewContext(nil, flag.NewFlagSet("local", 0), getParentCtx())

	code, err := Update(ctx, map[string]interface{}{"appName": "test"})
	if err == nil {
		t.Error("Expected to have an error")
	}
	if code != 1 {
		t.Error("Expected to have exit code 1")
	}
	if CheckResult != "" {
		t.Error("Expected to not apply update")
	}
}

func TestCmdUpdateYes(t *testing.T) {
	// mock app version
	tmpVersion := arcVersion.Version
	arcVersion.Version = "20150910.01"

	// mock apply upload
	origApplyUpdate := updater.ApplyUpdate
	updater.ApplyUpdate = mock_apply_update

	// mock server
	server := getMockServer(200, responseExample)
	defer server.Close()

	// prepare context flags
	ctx := cli.NewContext(nil, getLocalSet(false, false, server.URL), getParentCtx())

	// mock input to yes
	in, err := mockConfirmStdInput("yes")
	if err != nil {
		t.Error("Expected to not have an error")
	}
	confirmInput = in

	code, err := Update(ctx, map[string]interface{}{"appName": "test"})
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
		arcVersion.Version = tmpVersion
	}()
}

func TestCmdUpdateNo(t *testing.T) {
	// mock app version
	tmpVersion := arcVersion.Version
	arcVersion.Version = "20150910.01"

	// mock apply upload
	origApplyUpdate := updater.ApplyUpdate
	updater.ApplyUpdate = mock_apply_update

	// mock server
	server := getMockServer(200, responseExample)
	defer server.Close()

	// prepare context flags
	ctx := cli.NewContext(nil, getLocalSet(false, false, server.URL), getParentCtx())

	// mock input to no
	in, err := mockConfirmStdInput("no")
	if err != nil {
		t.Fatal(err)
	}
	confirmInput = in

	code, err := Update(ctx, map[string]interface{}{"appName": "test"})
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
		arcVersion.Version = tmpVersion
	}()
}

func TestCmdUpdateForce(t *testing.T) {
	// mock app version
	tmpVersion := arcVersion.Version
	arcVersion.Version = "20150910.01"

	// mock apply upload
	origApplyUpdate := updater.ApplyUpdate
	updater.ApplyUpdate = mock_apply_update

	// mock server
	server := getMockServer(200, responseExample)
	defer server.Close()

	// prepare context flags
	ctx := cli.NewContext(nil, getLocalSet(false, true, server.URL), getParentCtx())

	code, err := Update(ctx, map[string]interface{}{"appName": "test"})
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
		arcVersion.Version = tmpVersion
	}()
}

func TestCmdUpdateNoUpdate(t *testing.T) {
	// mock app version
	tmpVersion := arcVersion.Version
	arcVersion.Version = "20150910.01"

	// mock apply upload
	origApplyUpdate := updater.ApplyUpdate
	updater.ApplyUpdate = mock_apply_update

	// mock server
	server := getMockServer(200, responseExample)
	defer server.Close()

	// prepare context flags
	ctx := cli.NewContext(nil, getLocalSet(true, false, server.URL), getParentCtx())

	code, err := Update(ctx, map[string]interface{}{"appName": "test"})
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
		arcVersion.Version = tmpVersion
	}()
}

// private

func getLocalSet(noUpdate bool, force bool, serverUrl string) *flag.FlagSet {
	flagSet := flag.NewFlagSet("local", 0)
	flagSet.Bool("no-update", noUpdate, "local")
	flagSet.Bool("force", force, "local")
	flagSet.String("update-uri", serverUrl, "global")
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
	_, err = in.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return in, nil
}

func mock_apply_update(u *updater.Updater, r *updater.CheckResult) error {
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
