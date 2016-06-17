// +build !integration

package swift

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	//"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/ncw/swift/swifttest"
	"gitHub.***REMOVED***/monsoon/arc/updater"
)

var (
	testSrv *swifttest.SwiftServer
)

const (
	CONTAINER = "arc_releases_test"
)

//
// New()
//

func TestNewMissingParams(t *testing.T) {
	localSet := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(nil, localSet, nil)

	storage, err := New(ctx)
	if err == nil {
		t.Error("Expected to have an error")
	}
	if storage != nil {
		t.Error("Expected to have nil swift storage")
	}
}

func TestNew(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	if storage == nil {
		t.Error("Expected to have nil swift storage")
	}
}

//
// GetAvailableUpdate()
//

func TestGetAvailableUpdateSuccess(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	// save a file
	saveExamples(storage, t)

	// add checksum file
	err = storage.Connection.ObjectPutString(CONTAINER, "arc_20150906.07_linux_amd64.sha256", "checksum for arc_20150906.07_linux_amd64", "")
	if err != nil {
		t.Fatal(err)
	}

	jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","arch":"amd64","os":"linux"}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))

	update, err := storage.GetAvailableUpdate(req)
	if err != nil {
		t.Error(fmt.Sprint("Expected to not have an error. Got ", err))
		return
	}
	if update == nil {
		t.Error("Expected not nil")
		return
	}
	if !strings.Contains(update.Url, "arc_20150906.07_linux_amd64") {
		t.Error("Expected to get the file name in the update url")
	}
}

//
// Checksum
//

// func TestChecksumSuccess(t *testing.T) {
//   storage, err := getTestSwiftStorage()
//   if err != nil {
//     t.Error("Expected to have an error")
//     return
//   }
//   defer func() {
//     shutDownConnection()
//   }()
//
//   files := []string{"arc_20150905.15_linux_amd64", "arc_20150906.07_windows_amd64.exe"}
//   for _, filename := range files {
//     // save a file
//     err = storage.Connection.ObjectPutString(CONTAINER, filename, "123", "")
//     if err != nil {
//       t.Error(fmt.Sprint("Expected to not have an error. Got ", err))
//       return
//     }
//
//     //Checksum pattern "486c9e5b987027990865ed3109554cb6d9d6469397ea2ee0745999649defd203 *arc_20160321.2_linux_amd64"
//     objectBytes, _ := storage.Connection.ObjectGetBytes(CONTAINER, filename)
//     expectedChecksum, _ := checksumForBytes(objectBytes)
//     checksumData := fmt.Sprintf("%x *%s", expectedChecksum, filename)
//     storage.Connection.ObjectPutString(CONTAINER, fmt.Sprint(filename, ".sha256"), checksumData, "")
//
//     jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","arch":"amd64","os":"linux"}`)
//     req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
//
//     update, err := storage.GetAvailableUpdate(req)
//     if err != nil {
//       t.Error(fmt.Sprint("Expected to not have an error. Got ", err))
//       return
//     }
//     if update == nil {
//       t.Error("Expected not nil")
//       return
//     }
//
//     // check checksum is being added right
//     if update.Checksum != fmt.Sprintf("%x", expectedChecksum) {
//       t.Error("Expected to find checksum")
//     }
//
//     // compare checksum from result
//     decChecksum, err := hex.DecodeString(update.Checksum)
//     if err != nil {
//       t.Error(fmt.Sprint("Expected to not have an error. ", err))
//     }
//     if !bytes.Equal(expectedChecksum, decChecksum) {
//       t.Errorf("Updated file %s has wrong checksum. Expected: %x, got: %x", filename, expectedChecksum, decChecksum)
//     }
//   }
// }

func TestChecksumFail(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
		return
	}
	defer func() {
		shutDownConnection()
	}()

	// save a file
	filename := "arc_20150905.15_linux_amd64"
	err = storage.Connection.ObjectPutString(CONTAINER, filename, "123", "")
	if err != nil {
		t.Error(fmt.Sprint("Expected to not have an error. Got ", err))
		return
	}

	jsonStr := []byte(`{"app_id":"arc","app_version":"20150903.10","arch":"amd64","os":"linux"}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))

	update, err := storage.GetAvailableUpdate(req)
	if err == nil {
		t.Error("Expected to have an error.")
	}
	if update != nil {
		t.Error("Expected update to be nil")
	}
}

//
// Get updates
//

func TestGetAllUpdates(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	// save files
	saveExamples(storage, t)

	updates, err := storage.GetAllUpdates()
	if err != nil {
		t.Error("Expected to not have an error")
	}

	if len(*updates) != 6 {
		t.Error("Expected to find six releases")
	}
}

func TestGetWebUpdates(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	// save files
	saveExamples(storage, t)

	lastUpdates, allUpdates, err := storage.GetWebUpdates()
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

func TestGetAllUpdatesFilteredFiles(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	// save files
	saveExamples(storage, t)
	err = storage.Connection.ObjectPutString(CONTAINER, "readme.rm", "maiu", "")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Connection.ObjectPutString(CONTAINER, "releases.yaml", "bup", "")
	if err != nil {
		t.Fatal(err)
	}

	updates, err := storage.GetAllUpdates()
	if err != nil {
		t.Error("Expected to not have an error")
	}

	if len(*updates) != 6 {
		t.Error("Expected to find two releases")
	}
}

func TestGetUpdate(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	buildsRootPath, _ := ioutil.TempDir(os.TempDir(), "arc_builds_")
	target, _ := ioutil.TempFile(buildsRootPath, "target_file_")
	defer func() {
		os.RemoveAll(buildsRootPath)
		buildsRootPath = ""
	}()

	// save files
	saveExamples(storage, t)

	w := bufio.NewWriter(target)
	err = storage.GetUpdate("arc_20150905.10_linux_amd64", w)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	w.Flush() //ensure all buffered operations have been applied to the underlying writer

	content, _ := ioutil.ReadFile(target.Name())
	if string(content) != "123" {
		t.Error("Expected to get the source data in the target file")
	}
}

func TestGetLastestUpdate(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	// save files
	saveExamples(storage, t)

	windowsParams := updater.CheckParams{AppId: "arc", OS: "windows", Arch: "amd64"}
	latestUpdate, err := storage.GetLastestUpdate(&windowsParams)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if latestUpdate != "arc_20150906.07_windows_amd64.exe" {
		t.Error(fmt.Sprint("Expected to get last arc_20150906.07_windows_amd64. Got ", latestUpdate))
	}

	linuxParams := updater.CheckParams{AppId: "arc", OS: "linux", Arch: "amd64"}
	latestUpdate, err = storage.GetLastestUpdate(&linuxParams)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if latestUpdate != "arc_20150906.07_linux_amd64" {
		t.Error(fmt.Sprint("Expected to get last arc_20150906.07_linux_amd64. Got ", latestUpdate))
	}
}

func TestGetStoragePath(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	if testSrv.AuthURL != storage.Connection.AuthUrl {
		t.Error(fmt.Sprintf("Expected to get auth url %s. Got %s", testSrv.AuthURL, storage.Connection.AuthUrl))
	}
}

func TestIsConnectedSuccess(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	if storage.IsConnected() == false {
		t.Error("Expected to be connected")
	}
}

func TestIsConnectedFail(t *testing.T) {
	storage, err := getTestSwiftStorage()
	if err != nil {
		t.Error("Expected to have an error")
	}
	defer func() {
		shutDownConnection()
	}()

	tmpStorageUrl := storage.Connection.StorageUrl
	storage.Connection.StorageUrl = "http://miau.com"
	defer func() {
		storage.Connection.StorageUrl = tmpStorageUrl
	}()

	if storage.IsConnected() == true {
		t.Error("Expected to be not connected")
	}
}

// private

func shutDownConnection() {
	if testSrv != nil {
		testSrv.Close()
	}
}

func getTestSwiftStorage() (*SwiftStorage, error) {
	var err error

	// create a test server
	testSrv, err = swifttest.NewSwiftServer("localhost")
	if err != nil {
		return nil, err
	}

	// prepare flags
	localSet := flag.NewFlagSet("test", 0)
	localSet.String("username", "swifttest", "test")
	localSet.String("password", "swifttest", "test")
	localSet.String("domain", "test", "test")
	localSet.String("auth-url", testSrv.AuthURL, "test")
	localSet.String("container", CONTAINER, "test")
	ctx := cli.NewContext(nil, localSet, nil)

	//create storage
	storage, err := New(ctx)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func saveExamples(storage *SwiftStorage, t *testing.T) {
	var err error

	// save files
	err = storage.Connection.ObjectPutString(CONTAINER, "arc_20150905.10_linux_amd64", "123", "")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Connection.ObjectPutString(CONTAINER, "arc_20150905.10_windows_amd64.exe", "123", "")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Connection.ObjectPutString(CONTAINER, "arc_20150906.07_linux_amd64", "456", "")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Connection.ObjectPutString(CONTAINER, "arc_20150906.07_windows_amd64.exe", "456", "")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Connection.ObjectPutString(CONTAINER, "arc_20150805.15_linux_amd64", "789", "")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Connection.ObjectPutString(CONTAINER, "arc_20150805.15_windows_amd64.exe", "789", "")
	if err != nil {
		t.Fatal(err)
	}
}

// ChecksumForBytes returns the sha256 checksum for the given bytes
func checksumForBytes(source []byte) ([]byte, error) {
	return checksumForReader(bytes.NewReader(source))
}

func checksumForReader(rd io.Reader) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, rd); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
