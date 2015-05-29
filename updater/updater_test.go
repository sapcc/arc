package updater

import (
	"fmt"
	"github.com/inconshreveable/go-update/check"
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

	if up.updateUri != validOptions["updateUri"] {
		t.Error("Expected upater attribute 'updateUri' set to", validOptions["updateUri"], ", got ", up.updateUri)
	}
	if up.params.AppId != validOptions["appName"] {
		t.Error("Expected upater attribute 'AppId' set to", validOptions["appName"], ", got ", up.params.AppId)
	}
	if up.params.AppVersion != validOptions["version"] {
		t.Error("Expected upater attribute 'AppVersion' set to", validOptions["version"], ", got ", up.params.AppVersion)
	}
}

func TestUpdaterUpdateNotAvailable(t *testing.T) {
	server := testTools(204, "")
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	_, err := up.Update()
	if err != check.NoUpdateAvailable {
		t.Error("Expected get one error, got ", err)
	}
}

func TestUpdaterUpdateSuccess(t *testing.T) {
	// mock apply upload
	origApplyUpdate := applyUpdate
	applyUpdate = mock_apply_update
	defer func() { applyUpdate = origApplyUpdate }()

	// mock server
	server := testTools(200, `{"initiative":"automatically","url":"MIAU://non_valid_url","patch_url":null,"patch_type":null,"version":"999","checksum":null,"signature":null}`)
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	_, err := up.Update()
	if err != nil {
		t.Error("Expected get no error, got ", err)
	}
}

// private

func mock_apply_update(r *check.Result) error {
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
