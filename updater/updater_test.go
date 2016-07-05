// +build !integration

package updater

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var validOptions = map[string]string{
	"version":   "2.0",
	"appName":   "Miau",
	"updateUri": "http://localhost:3000/updates",
}

func TestUpdaterNewSuccess(t *testing.T) {
	up := New(validOptions)

	if up.client.Endpoint != validOptions["updateUri"] {
		t.Error("Expected upater attribute 'updateUri' set to", validOptions["updateUri"], ", got ", up.client.Endpoint)
	}
	if up.Params.AppId != validOptions["appName"] {
		t.Error("Expected upater attribute 'AppId' set to", validOptions["appName"], ", got ", up.Params.AppId)
	}
}

func TestUpdaterCheckNotAvailableWhenJSONnotFound(t *testing.T) {
	server := testTools(404, "")
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	res, err := up.Check()
	if res != nil {
		t.Error("Expected to not get a result")
	}
	if err == nil {
		t.Error("Expected to get an error, got ", err)
	}
}

func TestUpdaterCheckNotAvailable(t *testing.T) {
	server := testTools(200, `{
  "app_id": "arc",
  "os": "linux",
  "arch": "amd64",
  "checksum": "06285e7a9d85edee79c5d20732a66c92500a91fe003e1ff0f45de9abcb3d318b",
  "version": "20150910.01",
  "url":"20150910.01_linux_amd64"
}`)
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	res, err := up.Check()
	if res != nil {
		t.Error("Expected to not get a result")
	}
	if err == nil {
		t.Error("Expected to get an error, got ", err)
	}
}

func TestUpdaterCheckAvailable(t *testing.T) {
	server := testTools(200, `{
  "app_id": "arc",
  "os": "linux",
  "arch": "amd64",
  "checksum": "06285e7a9d85edee79c5d20732a66c92500a91fe003e1ff0f45de9abcb3d318b",
  "version": "20160620.01",
  "url":"20160620.01_linux_amd64"
}`)
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	res, err := up.Check()
	if res == nil {
		t.Error("Expected to not get a result")
		return
	}
	if err != nil {
		t.Error("Expected to not get an error, got ", err)
		return
	}

	if res.Version != "20160620.01" {
		t.Error("Expected to match the version, got ", err)
	}
}

func TestUpdaterCheckAndUpdateNotAvailableWhenJSONnotFound(t *testing.T) {
	server := testTools(404, "")
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	success, err := up.CheckAndUpdate()
	if success {
		t.Error("Expected to be false")
	}
	if err == nil {
		t.Error("Expected to get an error, got ", err)
	}
}

func TestUpdaterCheckAndUpdateFail(t *testing.T) {
	server := testTools(401, "")
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	success, err := up.CheckAndUpdate()
	if success {
		t.Error("Expected to be false")
	}
	if err == nil {
		t.Error("Expected to get an error, got ", err)
	}
}

func TestUpdaterCheckAndUpdateSuccess(t *testing.T) {
	// mock apply upload
	origApplyUpdate := ApplyUpdate
	ApplyUpdate = mock_apply_update
	defer func() { ApplyUpdate = origApplyUpdate }()

	// mock server
	server := testTools(200, `{"initiative":"automatically","url":"MIAU://non_valid_url","patch_url":null,"patch_type":null,"version":"999","checksum":null,"signature":null}`)
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	_, err := up.CheckAndUpdate()
	if err != nil {
		t.Error("Expected get no error, got ", err)
	}
}

// private

func mock_apply_update(u *Updater, r *CheckResult) error {
	return nil
}

func testTools(code int, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	return server
}
