package updates

import (
  "net/http"
	"github.com/inconshreveable/go-update/check"
)

func New(r *http.Request) *check.Result {	
	return &check.Result{
			Initiative: "automatically",
			Url: "http://localhost:3000/static/builds/arc_2",
			Version: "2",
		}
}

func getAvailableUpdate(appId string, os string, version string) *string{
	
	return nil
}