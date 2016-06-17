package local

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"gitHub.***REMOVED***/monsoon/arc/updater"

	"gitHub.***REMOVED***/monsoon/arc/update-server/storage/helpers"
)

var emptyPathError = "Builds root path is empty"

type LocalStorage struct {
	BuildsRootPath string
}

func New(c *cli.Context) (*LocalStorage, error) {
	if c.String("path") == "" {
		return nil, errors.New(emptyPathError)
	}

	// check if path exits
	if _, err := os.Stat(c.String("path")); os.IsNotExist(err) {
		return nil, err
	}

	return &LocalStorage{
		BuildsRootPath: c.String("path"),
	}, nil
}

/* build file pattern			  -> appId, "_", appVersion, "_", appOs, "_", appArch
 * Results:
 * nil, Error 					   	-> There is an error
 * Result, nil             	-> There is an available update
 * nil, nil						      -> No updates available
 */
func (l *LocalStorage) GetAvailableUpdate(req *http.Request) (*updater.CheckResult, error) {
	releases, err := l.GetAllUpdates()
	if err != nil {
		return nil, err
	}

	// get check.Params
	result, err := helpers.AvailableUpdate(req, releases)
	if err != nil {
		return nil, err
	}

	if result != nil {
		// get the filename from the url
		filename := helpers.GetChecksumFileName(result.Url)
		if len(filename) > 1 {
			// get the content of the checksum file
			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			err = l.GetUpdate(filename, w)
			if err != nil {
				return nil, errors.New(fmt.Sprint("Checksum file ", filename, " not found."))
			}
			checksum := strings.Split(b.String(), " ")

			if len(checksum) > 1 {
				result.Checksum = checksum[0]
			} else {
				return nil, errors.New("Checksum file pattern wrong")
			}
		}

		return result, nil
	}

	return nil, nil
}

func (l *LocalStorage) GetAllUpdates() (*[]string, error) {
	var filteredNames []string
	builds, err := ioutil.ReadDir(l.BuildsRootPath)
	if err != nil {
		return nil, err
	}

	// filter files
	for _, f := range builds {
		r := regexp.MustCompile(helpers.FileNameRegex)
		if r.MatchString(f.Name()) {
			filteredNames = append(filteredNames, f.Name())
		}
	}

	// sort releases by version
	helpers.SortByVersion(filteredNames)

	if len(filteredNames) == 0 {
		filteredNames = append(filteredNames, "No files found")
	}

	return &filteredNames, nil
}

func (l *LocalStorage) GetWebUpdates() (*[]string, *[]string, error) {
	// get all updates sorted by verison (latest first)
	updates, err := l.GetAllUpdates()
	if err != nil {
		return nil, nil, err
	}

	// get latest version
	latestVersion, err := helpers.GetLatestVersion(updates)
	if err != nil {
		return nil, nil, err
	}

	var latestUpdates []string
	var allUpdates []string

	for _, update := range *updates {
		version, err := helpers.ExtractVersion(update)
		if err != nil {
			continue
		}
		if latestVersion == version {
			latestUpdates = append(latestUpdates, update)
		} else {
			allUpdates = append(allUpdates, update)
		}
	}

	return &latestUpdates, &allUpdates, nil
}

func (l *LocalStorage) GetLastestUpdate(params *updater.CheckParams) (string, error) {
	// get all updates sorted by verison (latest first)
	updates, err := l.GetAllUpdates()
	if err != nil {
		return "", err
	}

	latestUpdate := helpers.GetLatestReleaseFrom(updates, params)
	return latestUpdate, nil
}

func (l *LocalStorage) GetUpdate(name string, writer io.Writer) error {
	// check if file exists
	path := path.Join(l.GetStoragePath(), name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return helpers.ObjectNotFoundError
	}

	// reader
	reader, err := os.Open(path)
	if err != nil {
		return err
	}

	// copy data from reader to writer
	_, err = io.Copy(writer, reader)
	return err
}

func (l *LocalStorage) GetStoragePath() string {
	return l.BuildsRootPath
}

// check if the path still exists
func (s *LocalStorage) IsConnected() bool {
	_, err := os.Stat(s.BuildsRootPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
