package local

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/inconshreveable/go-update/check"

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
func (l *LocalStorage) GetAvailableUpdate(req *http.Request) (*check.Result, error) {
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

func (l *LocalStorage) GetUpdate(name string, writer io.Writer) error {
	return nil
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
