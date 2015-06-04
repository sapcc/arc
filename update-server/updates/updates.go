package updates

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/inconshreveable/go-update/check"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var ArgumentError = fmt.Errorf("Build path is missing")

type availableUpdate struct {
	buildName string
	version   string
}

const buildRelativeUrl = "/builds/"

/*
 * Results:
 * nil, Error 						-> There is an error
 * *check.Result, nil			-> There is an available update result to send back
 * nil, nil								-> No updates available
 */
func New(req *http.Request, buildsRootPath string) (*check.Result, error) {
	// check arguments
	if len(buildsRootPath) == 0 || req == nil {
		return nil, ArgumentError
	}
	
	// get host url
	hostUrl := getHostUrl(req)
	if hostUrl == nil {
		return nil, errors.New(fmt.Sprintf("Computed host url is nil. Request %q", req))
	}

	// read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading the request body. Got %q", err))
	}
	if len(body) == 0 {
		return nil, errors.New(fmt.Sprintf("No request body. Got %q", body))
	}

	// convert to check.Params struc
	var reqParams check.Params
	err = json.Unmarshal(body, &reqParams)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshaling body. Got %q", err))
	}

	// check required post attributes
	if len(reqParams.AppId) == 0 {
		return nil, errors.New(fmt.Sprintf("Missing required post attribute 'app_id'. Got %q", reqParams.AppId))
	}
	if len(reqParams.AppVersion) == 0 {
		return nil, errors.New(fmt.Sprintf("Missing required post attribute 'app_version'. Got %q", reqParams.AppVersion))
	}
	if len(reqParams.Tags["os"]) == 0 {
		return nil, errors.New(fmt.Sprintf("Missing required post attribute 'tags[os]'. Got %q", reqParams.Tags["os"]))
	}
	if len(reqParams.Tags["arch"]) == 0 {
		return nil, errors.New(fmt.Sprintf("Missing required post attribute 'tags[arch]'. Got %q", reqParams.Tags["arch"]))
	}

	// get build url
	au, err := getAvailableUpdate(buildsRootPath, reqParams.AppId, reqParams.AppVersion, reqParams.Tags["os"], reqParams.Tags["arch"])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting available update %q", err))
	}
	if au != nil {
		return &check.Result{
			Initiative: "automatically",
			Url:        fmt.Sprint(hostUrl, buildRelativeUrl, au.buildName),
			Version:    au.version,
		}, nil
	}

	return nil, nil
}

// private

func getHostUrl(req *http.Request) *url.URL {
	// get the host
	host := req.Host
	if len(host) == 0 {
		return nil
	}

	// get the scheme
	scheme := ""
	// set the scheme
	if req.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}

	return &url.URL{Scheme: scheme, Host: host}
}

/* build file pattern			-> appId, "_", appOs, "_", appArch, "_", version
 * Results:
 * nil, Error 						-> There is an error
 * *availableUpdate, nil	-> There is an available update
 * nil, nil								-> No updates available
 */
func getAvailableUpdate(buildsRootPath string, appId string, appVersion string, appOs string, appArch string) (*availableUpdate, error) {
	av, err := semver.Make(appVersion)
	if err != nil {
		return nil, err
	}

	buildFile := ""
	buildVersion := "0.0.0"
	builds, err := ioutil.ReadDir(buildsRootPath)
	if err != nil {
		return nil, err
	}

	// loop over the builds and compare versions
	for _, f := range builds {
		if strings.HasPrefix(f.Name(), fmt.Sprint(appId, "_", appOs, "_", appArch, "_")) {
			fileVersion := strings.Split(f.Name(), "_")[3]
			bv, err := semver.Make(fileVersion)
			if err != nil {
				return nil, err
			}

			nbv, _ := semver.Make(buildVersion)
			if bv.Compare(av) == 1 && bv.Compare(nbv) == 1 {
				buildFile = f.Name()
				buildVersion = fileVersion
			}
		}
	}

	if len(buildFile) > 0 {
		return &availableUpdate{
			buildName: buildFile,
			version:   buildVersion,
		}, nil
	}

	return nil, nil
}
