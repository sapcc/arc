// +build !integration

package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

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
