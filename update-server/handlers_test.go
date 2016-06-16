// +build !integration

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/codegangsta/cli"

	"gitHub.***REMOVED***/monsoon/arc/update-server/storage"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

//
// Local storage - AvailableUpdate
//

func TestServeAvailableUpdates(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	checksum_data := "checksum data"
	filename := "arc_20150905.15_linux_amd64"
	err := createTestBuildFile(buildsRootPath, filename)
	if err != nil {
		t.Error(fmt.Sprint("Expected to not have an error. ", err))
	}
	err = createChecksumFile(buildsRootPath, filename, checksum_data)
	if err != nil {
		t.Error(fmt.Sprint("Expected to not have an error. ", err))
	}
	defer func() {
		os.RemoveAll(buildsRootPath)
	}()

	set := flag.NewFlagSet("test", 0)
	set.String("path", buildsRootPath, "local")
	c := cli.NewContext(nil, set, nil)
	st, _ = storage.New(storage.Local, c)

	jsonStr := []byte(`{"app_id":"arc","app_version":"20150901.01","arch":"amd64","os":"linux"}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	serveAvailableUpdates(w, req)
	if w.Code != 200 {
		t.Error("Expected code to be '200'. Got ", w.Code)
	}
}

func TestServeNonAvailableUpdates(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	defer func() {
		os.RemoveAll(buildsRootPath)
	}()

	set := flag.NewFlagSet("test", 0)
	set.String("path", buildsRootPath, "local")
	c := cli.NewContext(nil, set, nil)
	st, _ = storage.New(storage.Local, c)

	jsonStr := []byte(`{"app_id":"arc","app_version":"0.1.0-dev","arch":"amd64","os":"darwin"}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	serveAvailableUpdates(w, req)
	if w.Code != 204 {
		t.Error("Expected code to be '204'. Got ", w.Code)
	}
}

func TestServeAvailableUpdatesError(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	defer func() {
		os.RemoveAll(buildsRootPath)
	}()

	set := flag.NewFlagSet("test", 0)
	set.String("path", buildsRootPath, "local")
	c := cli.NewContext(nil, set, nil)
	st, _ = storage.New(storage.Local, c)

	jsonStr := []byte(`{"app_id":"arc","app_version":"0.1.0-dev","arch":"amd64"}`) // missing tag
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	serveAvailableUpdates(w, req)
	if w.Code != 400 {
		t.Error("Expected code to be '400'. Got ", w.Code)
	}
}

//
// Healthcheck
//

func TestHealthcheck(t *testing.T) {
	// make request
	req, err := http.NewRequest("GET", "/healthcheck", bytes.NewBufferString(""))
	if err != nil {
		t.Error("Expected not get an error")
	}
	w := httptest.NewRecorder()
	serveVersion(w, req)

	if w.Code != 200 {
		t.Error("Expected code to be '200'. Got ", w.Code)
	}
	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Error("Expected to get text/plain; charset=utf-8")
	}
	if w.Body.String() != fmt.Sprint("Arc update-server ", version.String()) {
		t.Error("Expected to get the health page")
	}
}

//
// Local storage - Upload
//

func TestUploadFilenameMissing(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	defer func() {
		os.RemoveAll(buildsRootPath)
	}()

	set := flag.NewFlagSet("test", 0)
	set.String("path", buildsRootPath, "local")
	c := cli.NewContext(nil, set, nil)
	st, _ = storage.New(storage.Local, c)

	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/upload", bytes.NewBuffer([]byte("binary file")))
	w := httptest.NewRecorder()
	uploadHandler(w, req)
	if w.Code != 400 {
		t.Error("Expected code to be '400'. Got ", w.Code)
	}
}

func TestUpload(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	defer func() {
		os.RemoveAll(buildsRootPath)
	}()

	set := flag.NewFlagSet("test", 0)
	set.String("path", buildsRootPath, "local")
	c := cli.NewContext(nil, set, nil)
	st, _ = storage.New(storage.Local, c)

	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/upload?filename=test", bytes.NewBuffer([]byte("binary file")))
	w := httptest.NewRecorder()
	uploadHandler(w, req)

	buildPath := path.Join(buildsRootPath, "test")
	if _, err := os.Stat(buildPath); os.IsNotExist(err) {
		t.Error("Expected to find build file")
	}
}

//
// helpers
//

func createTestBuildFile(buildsRootPath, name string) error {
	file, err := ioutil.TempFile(buildsRootPath, name)
	if err != nil {
		return err
	}
	err = os.Rename(file.Name(), path.Join(buildsRootPath, name))
	if err != nil {
		return err
	}

	return nil
}

func createChecksumFile(buildsRootPath, referenceFileName, checksumData string) error {
	// extract the temp file name
	i := strings.LastIndex(referenceFileName, "/")
	filename_ext := referenceFileName[i+1:]
	// create a checksum file without extra random data in the name
	checksum, _ := ioutil.TempFile(buildsRootPath, fmt.Sprint(filename_ext, ".sha256"))
	checksum.WriteString(checksumData)
	err := os.Rename(checksum.Name(), path.Join(buildsRootPath, fmt.Sprint(filename_ext, ".sha256")))
	if err != nil {
		return err
	}
	return nil
}
