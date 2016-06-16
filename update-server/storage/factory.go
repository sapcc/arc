package storage

import (
	"errors"
	"io"
	"net/http"

	"github.com/codegangsta/cli"
	"gitHub.***REMOVED***/monsoon/arc/updater"

	"gitHub.***REMOVED***/monsoon/arc/update-server/storage/local"
	"gitHub.***REMOVED***/monsoon/arc/update-server/storage/swift"
)

type StorageType int

const (
	_ StorageType = iota
	Local
	Swift
)

type Storage interface {
	GetAvailableUpdate(req *http.Request) (*updater.CheckResult, error)
	GetAllUpdates() (*[]string, error)
	GetWebUpdates() (*[]string, *[]string, error)
	GetUpdate(name string, writer io.Writer) error
	GetStoragePath() string
	GetLastestUpdate(params *updater.CheckParams) (string, error)
	IsConnected() bool
}

func New(storage StorageType, c *cli.Context) (Storage, error) {
	switch storage {
	case Local:
		return local.New(c)
	case Swift:
		return swift.New(c)
	}
	return nil, errors.New("Invalid storage")
}
