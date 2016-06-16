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
	if up.Params.AppVersion != validOptions["version"] {
		t.Error("Expected upater attribute 'AppVersion' set to", validOptions["version"], ", got ", up.Params.AppVersion)
	}
}

func TestUpdaterCheckAndUpdateNotAvailable(t *testing.T) {
	server := testTools(204, "")
	defer server.Close()

	// add the server url to the valid options to get a mock response
	validOptions["updateUri"] = server.URL

	up := New(validOptions)
	success, err := up.CheckAndUpdate()
	if success {
		t.Error("Expected to be false")
	}
	if err != nil {
		t.Error("Expected not get an error, got ", err)
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
