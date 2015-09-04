package local

import (
	"io/ioutil"
	"errors"
	"net/http"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/inconshreveable/go-update/check"
	"gitHub.***REMOVED***/monsoon/arc/update-server/storage/helpers"
)

type LocalStorage struct {
	BuildsRootPath string
}

func New(c *cli.Context) (*LocalStorage, error) {
	if c.String("path") == ""{
		return nil, errors.New("Build root path is empty")
	}
	
	// check if path exits
	
	return &LocalStorage{
		BuildsRootPath: c.String("path"),
	},nil
}

/* build file pattern			  -> appId, "_", appOs, "_", appArch, "_", appVersion
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
	var fileNames []string
	builds, err := ioutil.ReadDir(l.BuildsRootPath)
	if err != nil {
		return nil, err
	}
	
	for _, f := range builds {
		// filter config file
		if strings.ToLower(f.Name()) != "releases.yml" {
			fileNames = append(fileNames, f.Name())
		}
	}

	if len(fileNames) == 0 {
		fileNames = append(fileNames, "No files found")
	}

	return &fileNames, nil
}

func (l *LocalStorage) GetStoragePath() string{
	return l.BuildsRootPath
}