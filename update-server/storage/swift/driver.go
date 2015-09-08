package swift

import (
	"errors"
	"net/http"
	"io"
		
	"github.com/codegangsta/cli"
	"github.com/ncw/swift"
	"github.com/inconshreveable/go-update/check"
	"gitHub.***REMOVED***/monsoon/arc/update-server/storage/helpers"
)

type SwiftStorage struct {
	Connection  swift.Connection 
	Container   string
}

func New(c *cli.Context) (*SwiftStorage, error) {
	if c.String("username") == "" || c.String("password") == "" || c.String("domain") == "" || c.String("auth_url") == "" || c.String("container") == "" {
		return nil, errors.New("Not enough arguments in call swift new")
	}

	// create object
	swiftStorage := SwiftStorage{
			swift.Connection{
				UserName: c.String("username"),
				ApiKey:   c.String("password"),
				AuthUrl:  c.String("auth_url"),
				Domain:   c.String("domain"),
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

func (s *SwiftStorage) GetAvailableUpdate(req *http.Request) (*check.Result, error) {
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
		return result, nil
	}

	return nil, nil
}

func (s *SwiftStorage) GetAllUpdates() (*[]string, error) {
	names, err := s.Connection.ObjectNames(s.Container, nil)
	if err != nil {
		return nil, err
	}
	
	return &names, nil
}

func (s *SwiftStorage) GetUpdate(name string, writer io.Writer) (error) {
	_, err := s.Connection.ObjectGet(s.Container, name, writer, false, nil)
	if err != nil {
		return err
	}
	return nil		
}

func (s *SwiftStorage) GetStoragePath() string {
	return s.Connection.AuthUrl
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