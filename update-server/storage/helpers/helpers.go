package helpers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"fmt"
	"regexp"	
	
	"github.com/hashicorp/go-version"
	log "github.com/Sirupsen/logrus"
	"github.com/inconshreveable/go-update/check"
)

var UpdateArgumentError = fmt.Errorf("Update arguments are missing or wrong")

const BuildRelativeUrl = "/builds/"

// type Update struct {
// 	Uid      string `json:"uid"`
// 	Filename string `json:"filename"`
// 	App      string `json:"app"`
// 	Os       string `json:"os"`
// 	Arch     string `json:"arch"`
// 	Version  string `json:"version"`
// 	Date     string `json:"date"`
// }
//
// type Updates map[string]Update


/*
 * Results:
 * nil, Error 						-> There is an error
 * *check.Result, nil			-> There is an available update result to send back
 * nil, nil								-> No updates available
 */

func AvailableUpdate(req *http.Request, releases *[]string) (*check.Result, error) {
	// get check.Params
	reqParams, err := parseRequest(req)
	if err != nil {
		return nil, err
	}
	
	// get host url
	hostUrl := getHostUrl(req)
	if hostUrl == nil {
		return nil, fmt.Errorf("Computed host url is nil. Request %q", req)
	}

	buildFile := ""
	buildVersion := "20150101.01"	
	// loop over the releases and compare versions
	for _, f := range *releases {
		if isReleaseFrom(f, reqParams) {
			fileVersion, err := extractVersionFrom(f, reqParams)
			if err != nil {
				return nil, err
			}			
			result, err := shouldUpdate(reqParams.AppVersion, fileVersion, buildVersion)
			if err != nil {
				return nil, err
			}
			if result == true {
				buildFile = f
				buildVersion = fileVersion
			}
		}
	}

	if len(buildFile) > 0 {		
		return &check.Result{
			Initiative: "automatically",
			Url:        fmt.Sprint(hostUrl, BuildRelativeUrl, buildFile),
			Version:    buildVersion,
		}, nil		
	}
	
	return nil, nil
}


// private

func shouldUpdate(appVersion string, fileVersion string, currentVersion string) (bool, error) {
	av, err := version.NewVersion(appVersion)
	if err != nil {
		return false, err
	}
	fv, err := version.NewVersion(fileVersion)
	if err != nil {
		return false, err
	}
	cv, err := version.NewVersion(currentVersion)
	if err != nil {
		return false, err
	}	
	if fv.GreaterThan(av) && fv.GreaterThan(cv) {
		return true, nil
	}
	return false, nil
}

func isReleaseFrom(filename string, params *check.Params) bool {
	r, _ := regexp.Compile(fmt.Sprint(params.AppId, "_(.+)_", params.Tags["os"], "_", params.Tags["arch"]))
	return r.MatchString(filename)
}

func extractVersionFrom(filename string, params *check.Params) (string, error) {
	r, _ := regexp.Compile(fmt.Sprint(params.AppId, "_(.+)_", params.Tags["os"], "_", params.Tags["arch"]))		
	results := r.FindStringSubmatch(filename)
	if len(results) < 1 {
		return "", fmt.Errorf("Version could not be found.")
	}
	return results[1], nil
}

func parseRequest(req *http.Request) (*check.Params, error) {
	// check arguments
	if req == nil {
		log.Errorf("Request are empty or nil")
		return nil, UpdateArgumentError
	}

	// read body
	if req.Body == nil {
		return nil, fmt.Errorf("Error while reading the request body. Request body is nil")
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the request body. Got %q", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("No request body")
	}

	// convert to check.Params struc
	var reqParams check.Params
	err = json.Unmarshal(body, &reqParams)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling body. Got %q", err)
	}

	// check required post attributes
	if len(reqParams.AppId) == 0 {
		return nil, fmt.Errorf("Missing required post attribute 'app_id'")
	}
	if len(reqParams.AppVersion) == 0 {
		return nil, fmt.Errorf("Missing required post attribute 'app_version'")
	}
	if len(reqParams.Tags["os"]) == 0 {
		return nil, fmt.Errorf("Missing required post attribute 'tags[os]'")
	}
	if len(reqParams.Tags["arch"]) == 0 {
		return nil, fmt.Errorf("Missing required post attribute 'tags[arch]'")
	}

	return &reqParams, nil
}

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