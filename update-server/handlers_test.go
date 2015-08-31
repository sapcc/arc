// +build !integration

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestServeAvailableUpdates(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "arc_darwin_amd64_3.1.0-dev_")
	defer os.Remove(file.Name())

	buildsRootPath = os.TempDir()
	defer func() { buildsRootPath = "" }()

	jsonStr := []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64","os":"darwin"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	serveAvailableUpdates(w, req)
	if w.Code != 200 {
		t.Error("Expected code to be '200'. Got ", w.Code)
	}
}

func TestServeNonAvailableUpdates(t *testing.T) {
	buildsRootPath = os.TempDir()
	defer func() { buildsRootPath = "" }()

	jsonStr := []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64","os":"darwin"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	serveAvailableUpdates(w, req)
	if w.Code != 204 {
		t.Error("Expected code to be '204'. Got ", w.Code)
	}
}

func TestServeAvailableUpdatesError(t *testing.T) {
	defer func() { buildsRootPath = "" }()

	// Return 500 if the there was a problem with the builds path or request
	jsonStr := []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64","os":"darwin"}}`)
	paths := []string{"", "/non/existing/path"}
	for _, path := range paths {
		buildsRootPath = path
		req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		serveAvailableUpdates(w, req)
		if w.Code != 500 {
			t.Error("Expected code to be '500'. Got ", w.Code)
		}
	}

	// Return 400 if the body format is wrong or any param needed is missing
	body := []string{
		"",
		"not json",
		`{"param1":"param1"}`,
		`{"app_id":"arc"}`,
		`{"app_id":"arc","app_version":"0.1.0-dev"}`,
		`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"os":"darwin"}}`,
	}
	buildsRootPath = "/some/builds/path"
	for _, item := range body {
		jsonStr = []byte(item)
		req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		serveAvailableUpdates(w, req)
		if w.Code != 400 {
			t.Error("Expected code to be '400'. Got ", w.Code)
		}
	}
}

func TestUploadFilenameMissing(t *testing.T) {
	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()
	
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/upload", bytes.NewBuffer([]byte("binary file")))
	w := httptest.NewRecorder()
	uploadHandler(w, req)
	if w.Code != 400 {
		t.Error("Expected code to be '400'. Got ", w.Code)
	}
}

func TestUpload(t *testing.T) {
	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/upload?filename=test", bytes.NewBuffer([]byte("binary file")))
	w := httptest.NewRecorder()
	uploadHandler(w, req)
	
	buildPath := path.Join(buildsRootPath, "test")
	if _, err := os.Stat(buildPath); os.IsNotExist(err) {
		t.Error("Expected to find build file")
	}	
	if _, err := os.Stat(releasesConfigPath()); os.IsNotExist(err) {
		t.Error("Expected to find releases config file")
	}	
}