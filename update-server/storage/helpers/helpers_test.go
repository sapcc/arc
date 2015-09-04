// +build !integration

package helpers

import (
	"github.com/inconshreveable/go-update/check"	
  "testing"
	
	// "bytes"
	// "github.com/inconshreveable/go-update/check"
	// "io/ioutil"
	// "net/http"
	// "os"
	// "strings"
)

func TestIsRelease(t *testing.T) {
	windowsParams := check.Params{ AppId: "arc", Tags: map[string]string{ "os": "windows", "arch": "amd64", },}
	darwinParams := check.Params{ AppId: "arc", Tags: map[string]string{ "os": "darwin", "arch": "amd64", },}

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

func TestExtractVersionFromRelease(t *testing.T) {
	windowsParams := check.Params{ AppId: "arc", Tags: map[string]string{ "os": "windows", "arch": "amd64", },}	
	darwinParams := check.Params{ AppId: "arc", Tags: map[string]string{ "os": "darwin", "arch": "amd64", },}
	
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
}

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

// func TestUpdatesNewSuccess(t *testing.T) {
// 	file, _ := ioutil.TempFile(os.TempDir(), "arc_darwin_amd64_3.1.0-dev_")
// 	defer os.Remove(file.Name())
//
// 	// get a success update
// 	jsonStr := []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64","os":"darwin"}}`)
// 	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
// 	update, err := New(req, os.TempDir())
// 	if err != nil {
// 		t.Error("Expected not get an error. Got ", err)
// 	}
// 	if update == nil {
// 		t.Error("Expected update NOT to be nil. Got ", update)
// 	}
//
// 	if update.Initiative != "automatically" {
// 		t.Error("Expected Initiative to be 'automatically'. Got ", update.Initiative)
// 	}
//
// 	if !strings.HasPrefix(update.Url, "http://0.0.0.0:3000/builds/arc_darwin_amd64_3.1.0-dev") {
// 		t.Error("Expected url to be 'http://0.0.0.0:3000/builds/arc_darwin_amd64_3.1.0-dev'. Got ", update.Url)
// 	}
//
// 	if update.Version != "3.1.0-dev" {
// 		t.Error("Expected version to be '3.1.0-dev'. Got ", update.Version)
// 	}
// }
//
// func TestUpdatesNewReturnNil(t *testing.T) {
// 	var update *check.Result
// 	var req *http.Request
// 	var err error
//
// 	// post request host is missing or wrong
// 	hosts := []string{"", "miau"}
// 	for _, h := range hosts {
// 		req, _ = http.NewRequest("POST", h, bytes.NewBufferString(""))
// 		update, err = New(req, "/some/build/path/")
// 		if err == nil {
// 			t.Error("Expected err to be nil when testing wrong hosts. Got ", err)
// 		}
// 		if update != nil {
// 			t.Error("Expected update to be nil when testing wrong hosts. Got ", update)
// 		}
// 	}
//
// 	// post request body is empty
// 	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBufferString(""))
// 	update, err = New(req, "/some/build/path/")
// 	if err == nil {
// 		t.Error("Expected err to be nil when testing empty body. Got ", err)
// 	}
// 	if update != nil {
// 		t.Error("Expected update to be nil when testing empty body. Got ", update)
// 	}
//
// 	// check that the body is json
// 	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBufferString("not json"))
// 	update, err = New(req, "/some/build/path/")
// 	if err == nil {
// 		t.Error("Expected err to be nil when testing json body. Got ", err)
// 	}
// 	if update != nil {
// 		t.Error("Expected update to be nil when testing json body. Got ", update)
// 	}
//
// 	// check that the body has not the required params
// 	requiredParams := []string{
// 		`{"param1":"param1"}`,
// 		`{"app_id":"arc"}`,
// 		`{"app_id":"arc","app_version":"0.1.0-dev"}`,
// 		`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"os":"darwin"}}`,
// 	}
// 	var jsonStr []byte
// 	for _, p := range requiredParams {
// 		jsonStr = []byte(p)
// 		req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
// 		update, err = New(req, "/some/build/path/")
// 		if err == nil {
// 			t.Error("Expected err to be nil when testing required params. Got ", err)
// 		}
// 		if update != nil {
// 			t.Error("Expected update to be nil when testing required params. Got ", update)
// 		}
// 	}
//
// 	// check wrong version format
// 	jsonStr = []byte(`{"app_id":"arc","app_version":"miau","tags":{"arch":"amd64","os":"darwin"}}`) //
// 	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
// 	update, err = New(req, "/some/build/path/")
// 	if err == nil {
// 		t.Error("Expected err to be nil when testing wrong version format. Got ", err)
// 	}
// 	if update != nil {
// 		t.Error("Expected update to be nil when wrong version format. Got ", update)
// 	}
//
// 	// check wrong build path
// 	jsonStr = []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64","os":"darwin"}}`) //
// 	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
// 	update, err = New(req, "/some/build/path/")
// 	if err == nil {
// 		t.Error("Expected err to be nil when testing wrong build path. Got ", err)
// 	}
// 	if update != nil {
// 		t.Error("Expected update to be nil when wrong build path. Got ", update)
// 	}
// }
