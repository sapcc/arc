// +build !integration

package main

// import (
// 	"io/ioutil"
// 	"os"
// 	"path"
// 	"testing"
// 	// "fmt"
// )
//
// var release = Release{
// 	Uid: "arcwindowsamd64010dev",
// 	Filename: "arc_windows_amd64_0.1.0-dev",
// 	App: "Arc",
// 	Os: "Windows",
// 	Arch: "amd64",
// 	Version: "0.1.0-dev",
// 	Date: "01.01.2015",
// }
//
// var data = `arc_darwin_amd64_3.1.0-dev:
//   uid: arcdarwinamd64310dev
//   filename: arc_darwin_amd64_3.1.0-dev
//   app: arc
//   os: darwin
//   arch: amd64
//   version: 3.1.0-dev
//   date: 2015.07.15`
//
// func TestGetReleasesYmlEmptyPath(t *testing.T) {
// 	// no builds path given
// 	var releasesTest1 Releases
// 	releasesTest1.Read()
//
// 	if releasesTest1 != nil {
// 		t.Error("Should return an empty map")
// 	}
//
// 	// no file config saved
// 	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
// 	defer func() {
// 		os.RemoveAll(buildsRootPath)
// 		buildsRootPath = ""
// 	}()
//
// 	var releasesTest2 Releases
// 	releasesTest2.Read()
// 	if releasesTest2 != nil {
// 		t.Error("Should return an empty map")
// 	}
// }
//
// func TestGetReleasesYml(t *testing.T) {
// 	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
// 	file, err := os.Create(path.Join(buildsRootPath, "releases.yml"))
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	defer func() {
// 		file.Close()
// 		os.RemoveAll(buildsRootPath)
// 		buildsRootPath = ""
// 	}()
//
// 	file.WriteString(data)
// 	file.Sync()
//
// 	var releases Releases
// 	releases.Read()
// 	if len(releases) != 1 {
// 		t.Error("Expected to get 1 entries in the builds config file")
// 	}
// 	if releases["arc_darwin_amd64_3.1.0-dev"].Uid != "arcdarwinamd64310dev" {
// 		t.Error("Uid no match the one from the test")
// 	}
// 	if releases["arc_darwin_amd64_3.1.0-dev"].Filename != "arc_darwin_amd64_3.1.0-dev" {
// 		t.Error("Filename no match the one from the test")
// 	}
// 	if releases["arc_darwin_amd64_3.1.0-dev"].App != "arc" {
// 		t.Error("App name no match the one from the test")
// 	}
// 	if releases["arc_darwin_amd64_3.1.0-dev"].Os != "darwin" {
// 		t.Error("Os name no match the one from the test")
// 	}
// 	if releases["arc_darwin_amd64_3.1.0-dev"].Arch != "amd64" {
// 		t.Error("Arch name no match the one from the test")
// 	}
// 	if releases["arc_darwin_amd64_3.1.0-dev"].Version != "3.1.0-dev" {
// 		t.Error("Version no match the one from the test")
// 	}
// 	if releases["arc_darwin_amd64_3.1.0-dev"].Date != "2015.07.15" {
// 		t.Error("Date no match the one from the test")
// 	}
// }
//
// func TestUpdateReleasesYmlFileNoExists(t *testing.T) {
// 	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
// 	defer func() {
// 		os.RemoveAll(buildsRootPath)
// 		buildsRootPath = ""
// 	}()
//
// 	releases := Releases{}
// 	releases.Update("arc_windows_amd64_0.1.0-dev", release)
//
// 	releases.Read()
// 	testRelease := releases["arc_windows_amd64_0.1.0-dev"]
// 	compareReleases(testRelease, release, t)
// }
//
// func TestUpdateReleasesYml(t *testing.T) {
// 	buildsRootPath, _ = ioutil.TempDir(os.TempDir(), "arc_builds_")
// 	file, err := os.Create(path.Join(buildsRootPath, "releases.yml"))
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	defer func() {
// 		file.Close()
// 		os.RemoveAll(buildsRootPath)
// 		buildsRootPath = ""
// 	}()
//
// 	file.WriteString(data)
// 	file.Sync()
//
// 	releases := Releases{}
// 	releases.Update("arc_windows_amd64_0.1.0-dev", release)
//
// 	releases.Read()
// 	testRelease := releases["arc_windows_amd64_0.1.0-dev"]
// 	compareReleases(testRelease, release, t)
// }
//
// // private
//
// func compareReleases(release1 Release, release2 Release, t *testing.T) {
// 	if release1.Uid != release2.Uid {
// 		t.Error("Uid no match the one from the test")
// 	}
// 	if release1.Filename != release2.Filename {
// 		t.Error("Filename no match the one from the test")
// 	}
// 	if release1.App != release2.App {
// 		t.Error("App no match the one from the test")
// 	}
// 	if release1.Os != release2.Os {
// 		t.Error("Os no match the one from the test")
// 	}
// 	if release1.Arch != release2.Arch {
// 		t.Error("Arch no match the one from the test")
// 	}
// 	if release1.Version != release2.Version {
// 		t.Error("Version no match the one from the test")
// 	}
// 	if release1.Date != release2.Date {
// 		t.Error("Date no match the one from the test")
// 	}
// }