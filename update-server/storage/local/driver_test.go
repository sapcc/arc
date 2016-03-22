// +build !integration

package local

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/inconshreveable/go-update"
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
	checksum_data := "checksum data"
	filename := "arc_20150905.15_linux_amd64"
	err := createTestBuildFile(buildsRootPath, filename, "")
	err = createChecksumFile(buildsRootPath, filename, checksum_data)
	if err != nil {
		t.Error(fmt.Sprint("Expected to not have an error. ", err))
	}

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
		t.Error(fmt.Sprint("Expected to not have an error. ", err))
	}
	if update == nil {
		t.Error("Expected not nil")
	}
}

//
// Checksum
//

func TestChecksumFail(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	createTestBuildFile(buildsRootPath, "arc_20150905.15_linux_amd64", "")
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
	if err == nil {
		t.Error("Expected to have an error.")
	}
	if update != nil {
		t.Error("Expected update to be nil")
	}
}

func TestChecksumSuccess(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	createTestBuildFile(buildsRootPath, "arc_20150905.15_linux_amd64", "test test test")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	//Checksum pattern "486c9e5b987027990865ed3109554cb6d9d6469397ea2ee0745999649defd203 *arc_20160321.2_linux_amd64"
	expectedChecksum, err := update.ChecksumForFile(path.Join(buildsRootPath, "arc_20150905.15_linux_amd64"))
	err = createChecksumFile(buildsRootPath, "arc_20150905.15_linux_amd64", fmt.Sprintf("%x *arc_20150905.15_linux_amd64", expectedChecksum))
	if err != nil {
		t.Error(fmt.Sprint("Expected to not have an error. ", err))
	}

	jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","tags":{"arch":"amd64","os":"linux"}}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))

	ls := LocalStorage{
		BuildsRootPath: buildsRootPath,
	}
	update, err := ls.GetAvailableUpdate(req)
	if err != nil {
		t.Error(fmt.Sprint("Expected to not have an error. ", err))
	}
	if update == nil {
		t.Error("Expected not nil")
	}

	// check checksum is being added right
	if update.Checksum != fmt.Sprintf("%x", expectedChecksum) {
		t.Error("Expected to find checksum. Got ", update.Checksum, " but should be ", fmt.Sprintf("%x", expectedChecksum))
	}

	// compare checksum from result
	decChecksum, err := hex.DecodeString(update.Checksum)
	if err != nil {
		t.Error(fmt.Sprint("Expected to not have an error. ", err))
	}
	if !bytes.Equal(expectedChecksum, decChecksum) {
		t.Errorf("Updated file has wrong checksum. Expected: %x, got: %x", expectedChecksum, decChecksum)
	}
}

//
// Get updates
//

func TestGetUpdate(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	file, _ := ioutil.TempFile(buildsRootPath, "arc_20150905.15_linux_amd64_")
	target, _ := ioutil.TempFile(buildsRootPath, "target_file_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	data := "Some interesting data"
	file.WriteString(data)
	w := bufio.NewWriter(target)
	ls := LocalStorage{
		BuildsRootPath: buildsRootPath,
	}

	_, filename := path.Split(file.Name())
	err := ls.GetUpdate(filename, w)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	w.Flush() //ensure all buffered operations have been applied to the underlying writer

	content, _ := ioutil.ReadFile(target.Name())
	if string(content) != data {
		t.Error("Expected to get the source data in the target file")
	}
}

func TestGetAllUpdates(t *testing.T) {
	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	filename := "arc_20150905.15_linux_amd64"
	createTestBuildFile(buildsRootPath, filename, "")
	checksum_data := "checksum data"
	createChecksumFile(buildsRootPath, filename, checksum_data)
	filename2 := "arc_20150904.10_windows_amd64"
	createTestBuildFile(buildsRootPath, filename2, "")
	createChecksumFile(buildsRootPath, filename2, checksum_data)
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
		t.Error(fmt.Sprint("Expected to find two releases. Found Updates: ", len(*updates)))
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
	createTestBuildFile(buildsRootPath, "arc_20150905.10_linux_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150905.10_windows_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150906.07_linux_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150906.07_windows_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150805.15_linux_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150805.15_windows_amd64", "")
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
	createTestBuildFile(buildsRootPath, "arc_20150905.10_linux_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150905.10_windows_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150906.07_linux_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150906.07_windows_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150805.15_linux_amd64", "")
	createTestBuildFile(buildsRootPath, "arc_20150805.15_windows_amd64", "")
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

//
// helpers
//

func createTestBuildFile(buildsRootPath, name, data string) error { //*os.File,
	file, err := os.Create(path.Join(buildsRootPath, name))
	if err != nil {
		return err
	}
	defer file.Close()
	if len(data) > 0 {
		_, err := file.WriteString(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func createChecksumFile(buildsRootPath, referenceFileName, checksumData string) error {
	// extract the temp file name
	i := strings.LastIndex(referenceFileName, "/")
	filename_ext := referenceFileName[i+1:]

	// create a checksum file without extra random data in the name
	checksum, err := os.Create(path.Join(buildsRootPath, fmt.Sprint(filename_ext, ".sha256")))
	if err != nil {
		return err
	}
	defer checksum.Close()
	if len(checksumData) > 0 {
		_, err := checksum.WriteString(checksumData)
		if err != nil {
			return err
		}
	}
	return nil
}
