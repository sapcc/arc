// +build !integration

package helpers

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/inconshreveable/go-update/check"
)

//
// IsRelease
//

func TestIsRelease(t *testing.T) {
	windowsParams := check.Params{AppId: "arc", Tags: map[string]string{"os": "windows", "arch": "amd64"}}
	darwinParams := check.Params{AppId: "arc", Tags: map[string]string{"os": "darwin", "arch": "amd64"}}

	result := isReleaseFrom("", &windowsParams)
	if result != false {
		t.Error("Expected to not be a release file an emtpy string")
	}

	result = isReleaseFrom("arc_20150903.10_windows_amd64.exe", &windowsParams)
	if result != true {
		t.Error("Expected to be a release file arc_20150903.10_windows_amd64.exe")
	}

	result = isReleaseFrom("arc_20150903.10_darwin_amd64", &darwinParams)
	if result != true {
		t.Error("Expected to be a release file arc_20150903.10_darwin_amd64")
	}

	result = isReleaseFrom("arc_darwin_amd64_3.1.0-dev", &darwinParams)
	if result != false {
		t.Error("Expected to be a release file arc_darwin_amd64_3.1.0-dev")
	}

}

//
// ExtractVersionFromRelease
//

func TestExtractVersionFromRelease(t *testing.T) {
	windowsParams := check.Params{AppId: "arc", Tags: map[string]string{"os": "windows", "arch": "amd64"}}
	darwinParams := check.Params{AppId: "arc", Tags: map[string]string{"os": "darwin", "arch": "amd64"}}

	_, err := extractVersionFrom("arc_20150903.10_windows_amd64.exe", &darwinParams)
	if err == nil {
		t.Error("Expected to have an error")
	}

	result, err := extractVersionFrom("arc_20150903.10_darwin_amd64", &darwinParams)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if result != "20150903.10" {
		t.Error("Expected to find version 20150903.10")
	}

	result, err = extractVersionFrom("arc_20150903.10_windows_amd64", &windowsParams)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if result != "20150903.10" {
		t.Error("Expected to find version 20150903.10")
	}

	result, err = extractVersionFrom("arc_20150905.15_darwin_amd64_061430944", &darwinParams)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if result != "20150905.15" {
		t.Error("Expected to find version 20150903.10")
	}
}

//
// ShouldUpdate
//

func TestShouldUpdate(t *testing.T) {
	// file version is greater than the app version and current version
	result, err := shouldUpdate("20150801.10", "20150903.10", "20150101.01")
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if result != true {
		t.Error("Expected to update version")
	}

	// file version is not grater than the current version
	result, err = shouldUpdate("20150801.10", "20150903.10", "20150903.15")
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if result != false {
		t.Error("Expected to not update version")
	}

	// file version is not grater than the current version
	result, err = shouldUpdate("20150801.10", "20150703.01", "20150705.10")
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if result != false {
		t.Error("Expected to not update version")
	}
}

//
// GetHostUrl
//

func TestGetHostUrl(t *testing.T) {
	req, _ := http.NewRequest("POST", "https://localhost:3000/updates", bytes.NewBufferString(""))
	req.TLS = &tls.ConnectionState{}
	url := getHostUrl(req)

	if url.String() != "https://localhost:3000" {
		t.Error("Expected schema https and host localhost:3000")
	}

	req2, _ := http.NewRequest("POST", "http://localhost:3000/updates", bytes.NewBufferString(""))
	url2 := getHostUrl(req2)
	if url2.String() != "http://localhost:3000" {
		t.Error("Expected schema http and host localhost:3000")
	}
}

//
// ParseRequest()
//

func TestParseRequestRequestNil(t *testing.T) {
	params, err := parseRequest(nil)
	if err == nil {
		t.Error("Expected to get an error")
	}
	if params != nil {
		t.Error("Expected to nil params")
	}
}

func TestParseRequestEmptyBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBufferString(""))
	params, err := parseRequest(req)
	if err == nil {
		t.Error("Expected get an error")
	}
	if params != nil {
		t.Error("Expected to nil params")
	}
}

func TestParseRequestEmptyBodyNotJson(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBufferString("not json"))
	params, err := parseRequest(req)
	if err == nil {
		t.Error("Expected get an error")
	}
	if params != nil {
		t.Error("Expected to nil params")
	}
}

func TestParseRequestcheckMissingParams(t *testing.T) {
	// get a success update
	jsonStr := []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	params, err := parseRequest(req)
	if err == nil {
		t.Error("Expected to get an error")
	}
	if params != nil {
		t.Error("Expected to nil params")
	}
}

func TestParseRequestSuccessfulcheckParams(t *testing.T) {
	// get a success update
	jsonStr := []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64","os":"darwin"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	params, err := parseRequest(req)
	if err != nil {
		t.Error("Expected not to get an error")
	}
	if params.AppId != "arc" {
		t.Error("Missing required post attribute 'app_id'")
	}
	if params.AppVersion != "0.1.0-dev" {
		t.Error("Missing required post attribute 'app_version'")
	}
	if params.Tags["os"] != "darwin" {
		t.Error("Missing required post attribute 'tags[os]'")
	}
	if params.Tags["arch"] != "amd64" {
		t.Error("Missing required post attribute 'tags[arch]'")
	}
}

//
// AvailableUpdate()
//

func TestAvailableUpdate(t *testing.T) {
	// get a success update
	jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","tags":{"arch":"amd64","os":"linux"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	releases := []string{"arc_20150903.10_linux_amd64", "arc_20150903.10_windows_amd64.exe", "arc_20150903.5_windows_amd64.exe", "arc_20150904.1_linux_amd64"}

	update, err := AvailableUpdate(req, &releases)
	if err != nil {
		t.Error("Expected not get an error. Got ", err)
	}
	if update == nil {
		t.Error("Expected update NOT to be nil. Got ", update)
	}

	if update.Initiative != "automatically" {
		t.Error("Expected Initiative to be 'automatically'. Got ", update.Initiative)
	}

	if update.Url != "http://0.0.0.0:3000/builds/arc_20150904.1_linux_amd64" {
		t.Error("Expected url to be 'http://0.0.0.0:3000/builds/arc_20150904.1_linux_amd64'. Got ", update.Url)
	}

	if update.Version != "20150904.1" {
		t.Error("Expected version to be '20150904.1'. Got ", update.Version)
	}
}

func TestNoAvailableUpdate(t *testing.T) {
	jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","tags":{"arch":"amd64","os":"linux"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	releases := []string{"arc_20150903.10_linux_amd64", "arc_20150903.10_windows_amd64.exe", "arc_20150903.5_windows_amd64.exe", "arc_20150902.1_linux_amd64"}

	update, err := AvailableUpdate(req, &releases)
	if err != nil {
		t.Error("Expected not get an error. Got ", err)
	}
	if update != nil {
		t.Error("Expected update to be nil")
	}
}

func TestAvailableUpdateArgumentError(t *testing.T) {
	jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","tags":{"arch":"amd64"}}`) // parse error (missing parameter) == UpdateArgumentError == 400
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	releases := []string{}

	_, err := AvailableUpdate(req, &releases)
	if err != UpdateArgumentError {
		t.Error("Expected argument error")
	}
}

func TestAvailableUpdateWithOtherFiles(t *testing.T) {
	// get a success update
	jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","tags":{"arch":"amd64","os":"linux"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	releases := []string{"README.md", "Vagrantfile", "arc_test_linux_amd64", "arc_20150904.1_linux_amd64"}

	update, err := AvailableUpdate(req, &releases)
	if err != nil {
		t.Error("Expected not get an error. Got ", err)
	}
	if update == nil {
		t.Error("Expected update NOT to be nil. Got ", update)
	}

	if update.Initiative != "automatically" {
		t.Error("Expected Initiative to be 'automatically'. Got ", update.Initiative)
	}

	if update.Url != "http://0.0.0.0:3000/builds/arc_20150904.1_linux_amd64" {
		t.Error("Expected url to be 'http://0.0.0.0:3000/builds/arc_20150904.1_linux_amd64'. Got ", update.Url)
	}

	if update.Version != "20150904.1" {
		t.Error("Expected version to be '20150904.1'. Got ", update.Version)
	}
}

//
// Sort by Version
//

func TestSortByVersion(t *testing.T) {
	filenames := []string{"arc_20150903.5_windows_amd64.exe", "arc_20150803.5_linux_amd64", "arc_20151003.7_linux_amd64", "arc_201501003.8_linux_amd64", "arc_20151003.1_windows_amd64.exe"}
	sortedFilenames := []string{"arc_201501003.8_linux_amd64", "arc_20151003.7_linux_amd64", "arc_20151003.1_windows_amd64.exe", "arc_20150903.5_windows_amd64.exe", "arc_20150803.5_linux_amd64"}
	SortByVersion(filenames)

	for i, file := range filenames {
		if file != sortedFilenames[i] {
			t.Error(fmt.Sprint("Expected sorted filenames. Got ", file, " and ", sortedFilenames[i]))
		}
	}
}
