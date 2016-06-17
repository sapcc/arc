package swift

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/ncw/swift"
	"gitHub.***REMOVED***/monsoon/arc/updater"

	"gitHub.***REMOVED***/monsoon/arc/update-server/storage/helpers"
)

type SwiftStorage struct {
	Connection swift.Connection
	Container  string
}

func New(c *cli.Context) (*SwiftStorage, error) {
	if c.String("username") == "" || c.String("password") == "" || c.String("domain") == "" || c.String("auth-url") == "" || c.String("container") == "" {
		return nil, errors.New("Not enough arguments in call swift new")
	}

	// create object
	swiftStorage := SwiftStorage{
		swift.Connection{
			UserName:  c.String("username"),
			ApiKey:    c.String("password"),
			AuthUrl:   c.String("auth-url"),
			Domain:    c.String("domain"),
			TenantId:  c.String("project-id"),
			Retries:   1,
			UserAgent: fmt.Sprintf("%s (arc update-site; container: %s)", swift.DefaultUserAgent, c.String("container")),
		},
		c.String("container"),
	}

	// authenticate
	err := swiftStorage.Connection.Authenticate()
	if err != nil {
		return nil, err
	}

	// check and create container
	err = swiftStorage.CheckAndCreateContainer()
	if err != nil {
		return nil, err
	}

	return &swiftStorage, nil
}

func (s *SwiftStorage) GetAvailableUpdate(req *http.Request) (*updater.CheckResult, error) {
	releases, err := s.GetAllUpdates()
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
			err = s.GetUpdate(filename, w)
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

func (s *SwiftStorage) GetAllUpdates() (*[]string, error) {
	var filteredNames []string
	names, err := s.Connection.ObjectNames(s.Container, nil)
	if err != nil {
		return nil, err
	}

	// filter build files
	for _, name := range names {
		r := regexp.MustCompile(helpers.FileNameRegex)
		if r.MatchString(name) {
			// TODO get checksum
			filteredNames = append(filteredNames, name)
		}
	}

	// sort releases by version
	helpers.SortByVersion(filteredNames)

	return &filteredNames, nil
}

func (s *SwiftStorage) GetWebUpdates() (*[]string, *[]string, error) {
	// get all updates sorted by verison (latest first)
	updates, err := s.GetAllUpdates()
	if err != nil {
		return nil, nil, err
	}

	// get latest version
	latestVersion, _ := helpers.GetLatestVersion(updates)

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

func (s *SwiftStorage) GetLastestUpdate(params *updater.CheckParams) (string, error) {
	// get all updates sorted by verison (latest first)
	updates, err := s.GetAllUpdates()
	if err != nil {
		return "", err
	}

	latestUpdate := helpers.GetLatestReleaseFrom(updates, params)
	return latestUpdate, nil
}

func (s *SwiftStorage) GetUpdate(name string, writer io.Writer) error {
	_, err := s.Connection.ObjectGet(s.Container, name, writer, true, nil)
	if err == swift.ObjectNotFound {
		return helpers.ObjectNotFoundError
	} else if err != nil {
		return err
	}
	return nil
}

func (s *SwiftStorage) GetStoragePath() string {
	return s.Connection.AuthUrl
}

func (s *SwiftStorage) IsConnected() bool {
	_, _, err := s.Connection.Container(s.Container)
	if err != nil {
		return false
	}
	return true
}

// private

func (s *SwiftStorage) CheckAndCreateContainer() error {
	_, _, err := s.Connection.Container(s.Container)
	if err == swift.ContainerNotFound {
		err = s.Connection.ContainerCreate(s.Container, nil)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
