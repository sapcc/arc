// +build !integration

package local

import (
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/codegangsta/cli"
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
