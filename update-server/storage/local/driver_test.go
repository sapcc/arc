// +build !integration

package local

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/inconshreveable/go-update/check"
)

//
// New()
//

func TestNewEmptyPath(t *testing.T) {
	localSet := flag.NewFlagSet("test", 0)
	localSet.String("path", "", "test")
	ctx := cli.NewContext(nil, localSet, nil)

	ls, err := New(ctx)
	if err.Error() != emptyPathError {
		t.Error("Expected to have an empty path error")
	}
	if ls != nil {
		t.Error("Expected to have nil local storage")
	}
}

func TestNewPathNotExists(t *testing.T) {
	localSet := flag.NewFlagSet("test", 0)
	localSet.String("path", "some/non/existing/path", "test")
	ctx := cli.NewContext(nil, localSet, nil)

	ls, err := New(ctx)
	if err == nil || err.Error() == emptyPathError {
		t.Error("Expected to have an error")
	}
	if ls != nil {
		t.Error("Expected to have nil local storage")
	}
}

func TestNew(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	localSet := flag.NewFlagSet("test", 0)
	localSet.String("path", buildsRootPath, "test")
	ctx := cli.NewContext(nil, localSet, nil)

	ls, err := New(ctx)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if ls.BuildsRootPath != buildsRootPath {
		t.Error("Expected to find the buildsRootPath")
	}
}

//
// GetAvailableUpdate()
//

func TestGetAvailableUpdateSuccess(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	ioutil.TempFile(buildsRootPath, "arc_20150905.15_linux_amd64_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","tags":{"arch":"amd64","os":"linux"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))

	ls := LocalStorage{
		BuildsRootPath: buildsRootPath,
	}
	update, err := ls.GetAvailableUpdate(req)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if update == nil {
		t.Error("Expected not nil")
	}
}

//
// GetAllUpdates()
//

func TestGetAllUpdates(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	ioutil.TempFile(buildsRootPath, "arc_20150905.15_linux_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150904.10_windows_amd64_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	ls := LocalStorage{
		BuildsRootPath: buildsRootPath,
	}
	updates, err := ls.GetAllUpdates()
	if err != nil {
		t.Error("Expected to not have an error")
	}

	if len(*updates) != 2 {
		t.Error("Expected to find two releases")
	}
}

func TestGetAllUpdatesFilteredFiles(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	ioutil.TempFile(buildsRootPath, "arc_20150905.15_linux_amd64_")
	ioutil.TempFile(buildsRootPath, "readme.rm")
	ioutil.TempFile(buildsRootPath, "releases.yaml")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	ls := LocalStorage{
		BuildsRootPath: buildsRootPath,
	}
	updates, err := ls.GetAllUpdates()
	if err != nil {
		t.Error("Expected to not have an error")
	}

	if len(*updates) != 1 {
		t.Error("Expected to find two releases")
	}
}

func TestGetWebUpdates(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	ioutil.TempFile(buildsRootPath, "arc_20150905.10_linux_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150905.10_windows_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150906.07_linux_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150906.07_windows_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150805.15_linux_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150805.15_windows_amd64_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	ls := LocalStorage{
		BuildsRootPath: buildsRootPath,
	}
	lastUpdates, allUpdates, err := ls.GetWebUpdates()
	if err != nil {
		t.Error("Expected to not have an error")
	}

	if len(*lastUpdates) != 2 {
		t.Error("Expected to find two releases")
	}
	if len(*allUpdates) != 4 {
		t.Error("Expected to find two releases")
	}
}

func TestGetLastestUpdate(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	ioutil.TempFile(buildsRootPath, "arc_20150905.10_linux_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150905.10_windows_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150906.07_linux_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150906.07_windows_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150805.15_linux_amd64_")
	ioutil.TempFile(buildsRootPath, "arc_20150805.15_windows_amd64_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	ls := LocalStorage{
		BuildsRootPath: buildsRootPath,
	}

	windowsParams := check.Params{AppId: "arc", Tags: map[string]string{"os": "windows", "arch": "amd64"}}
	latestUpdate, err := ls.GetLastestUpdate(&windowsParams)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if !strings.Contains(latestUpdate, "arc_20150906.07_windows_amd64") {
		t.Error(fmt.Sprint("Expected to get last arc_20150906.07_windows_amd64. Got ", latestUpdate))
	}

	linuxParams := check.Params{AppId: "arc", Tags: map[string]string{"os": "linux", "arch": "amd64"}}
	latestUpdate, err = ls.GetLastestUpdate(&linuxParams)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if !strings.Contains(latestUpdate, "arc_20150906.07_linux_amd64") {
		t.Error(fmt.Sprint("Expected to get last arc_20150906.07_linux_amd64. Got ", latestUpdate))
	}
}
