// +build !integration

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

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
