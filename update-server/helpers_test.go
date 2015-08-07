// +build !integration

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"path"
)

func TestGetReleasesYmlEmptyPath(t *testing.T) {
	// no path given
	releases, _ := getBuildExtraInfo()
	if releases != nil {
		t.Error("Should return an empty map")
	}
	
	// no file config given
	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")	
	defer func() {
		buildsRootPath = ""		
	}()
		
	releases, _ = getBuildExtraInfo()
	if releases != nil {
		t.Error("Should return an empty map")
	}
}

func TestGetReleasesYml(t *testing.T) {
	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
	file, err := os.Create(path.Join(buildsRootPath, "releases.yml"))
	if err != nil {
		t.Error(err)
	}
	
	defer func() {
		file.Close()
		os.Remove(file.Name())
		os.Remove(buildsRootPath)
		buildsRootPath = ""		
	}()

	data := `arc_darwin_amd64_3.1.0-dev:
  uid: arcdarwinamd64310dev
  filename: arc_darwin_amd64_3.1.0-dev
  app: arc
  os: darwin
  arch: amd64
  version: 3.1.0-dev
  date: 2015.07.15`

	file.WriteString(data)
	file.Sync()

	releases, _ := getBuildExtraInfo()
	if len(releases) != 1 {
		t.Error("Expected to get 1 entries in the builds config file")
	}
	if releases["arc_darwin_amd64_3.1.0-dev"].Uid != "arcdarwinamd64310dev" {
		t.Error("Uid no match the one from the test")
	}	
	if releases["arc_darwin_amd64_3.1.0-dev"].Filename != "arc_darwin_amd64_3.1.0-dev" {
		t.Error("Filename no match the one from the test")
	}
	if releases["arc_darwin_amd64_3.1.0-dev"].App != "arc" {
		t.Error("App name no match the one from the test")
	}
	if releases["arc_darwin_amd64_3.1.0-dev"].Os != "darwin" {
		t.Error("Os name no match the one from the test")
	}
	if releases["arc_darwin_amd64_3.1.0-dev"].Arch != "amd64" {
		t.Error("Arch name no match the one from the test")
	}
	if releases["arc_darwin_amd64_3.1.0-dev"].Version != "3.1.0-dev" {
		t.Error("Version no match the one from the test")
	}
	if releases["arc_darwin_amd64_3.1.0-dev"].Date != "2015.07.15" {
		t.Error("Date no match the one from the test")
	}	
}

func TestGetAllBuildsEmpty(t *testing.T) {
	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
	defer func() {
		os.Remove(buildsRootPath)
		buildsRootPath = ""
	}()

	builds := getAllBuilds()
	if len(*builds) != 1 {
		t.Error("Expected to get 1 builds")
	}

	if (*builds)[0] != "No files found" {
		t.Error("Expected get 'No files found' text")
	}
}

func TestGetAllBuilds(t *testing.T) {
	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
	file1, _ := ioutil.TempFile(buildsRootPath, "arc_darwin_amd64_3.1.0-dev_")
	file2, _ := ioutil.TempFile(buildsRootPath, "arc_windows_amd64_3.1.0-dev_")
	defer func() {
		os.Remove(file1.Name())
		os.Remove(file2.Name())
		os.Remove(buildsRootPath)
		buildsRootPath = ""
	}()

	builds := getAllBuilds()
	if len(*builds) != 2 {
		t.Error("Expected to get 2 builds")
	}
	for _, f := range *builds {
		if !strings.HasSuffix(file1.Name(), f) && !strings.HasSuffix(file2.Name(), f) {
			t.Error("Expected to find build file")
		}
	}
}

func TestGetAllBuildsFilter(t *testing.T) {
	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
	file1, err := os.Create(path.Join(buildsRootPath, "releases.yml"))
	if err != nil {
		t.Error(err)
	}	
	file2, _ := ioutil.TempFile(buildsRootPath, "arc_windows_amd64_3.1.0-dev_")
	defer func() {
		file1.Close()
		os.Remove(file1.Name())
		os.Remove(file2.Name())
		os.Remove(buildsRootPath)
		buildsRootPath = ""
	}()
	
	builds := getAllBuilds()
	if len(*builds) != 1 {
		t.Error("Expected to get 1 builds")
	}
	
	for _, f := range *builds {
		if !strings.HasSuffix(file2.Name(), f) {
			t.Error("Expected to find build file")
		}
	}
}
