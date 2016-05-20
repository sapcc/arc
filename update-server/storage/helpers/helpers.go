package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/databus23/requestutil"
	version "github.com/hashicorp/go-version"
	"github.com/inconshreveable/go-update/check"
)

var UpdateArgumentError = fmt.Errorf("Update arguments are missing or wrong")
var ObjectNotFoundError = fmt.Errorf("Object not found.")

const (
	BuildRelativeUrl = "/builds/"
	FileNameRegex    = `^(?P<app>[^_]+)_(?P<version>[.0-9]+)_(?P<platform>windows|linux|darwin)_(?P<arch>amd64|386)(.exe)?$`
)

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
		return nil, UpdateArgumentError
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
				log.Warn(err)
				continue
			}
			result, err := shouldUpdate(reqParams.AppVersion, fileVersion, buildVersion)
			if err != nil {
				log.Warn(err)
				continue
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

func SortByVersion(filenames []string) {
	sort.Sort(ByVersion(filenames))
}

func ExtractVersion(filename string) (string, error) {
	r := regexp.MustCompile(FileNameRegex)
	results := r.FindStringSubmatch(filename)
	if len(results) < 3 {
		return "", fmt.Errorf("Version could not be found.")
	}
	return results[2], nil
}

func GetLatestVersion(releases *[]string) (string, error) {
	// sort releases by version
	SortByVersion(*releases)

	// get las version
	latestVersion := ""
	var err error
	if len(*releases) > 0 {
		latestVersion, err = ExtractVersion((*releases)[0])
		if err != nil {
			return "", err
		}
	}
	return latestVersion, nil
}

func GetLatestReleaseFrom(releases *[]string, params *check.Params) string {
	// sort releases by version
	SortByVersion(*releases)

	// take the first relases that match the params
	lastRelease := ""
	for _, release := range *releases {
		found := isReleaseFrom(release, params)
		if found == true {
			lastRelease = release
			break
		}
	}
	return lastRelease
}

func GetChecksumFileName(resultUrl string) string {
	i := strings.LastIndex(resultUrl, "/")
	filename := resultUrl[i+1:]
	filename = fmt.Sprint(filename, ".sha256")
	return filename
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
	r, _ := regexp.Compile(fmt.Sprint(params.AppId, "_([.0-9]+)_", params.Tags["os"], "_", params.Tags["arch"]))
	return r.MatchString(filename)
}

func extractVersionFrom(filename string, params *check.Params) (string, error) {
	r, _ := regexp.Compile(fmt.Sprint(params.AppId, "_([.0-9]+)_", params.Tags["os"], "_", params.Tags["arch"]))
	results := r.FindStringSubmatch(filename)
	if len(results) < 1 {
		return "", fmt.Errorf("Version could not be found.")
	}
	return results[1], nil
}

func parseRequest(req *http.Request) (*check.Params, error) {
	// check arguments
	if req == nil {
		return nil, fmt.Errorf("Request are empty or nil")
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
	return &url.URL{Scheme: requestutil.Scheme(req), Host: requestutil.HostWithPort(req)}
}

type ByVersion []string

func (s ByVersion) Len() int {
	return len(s)
}

func (s ByVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByVersion) Less(i, j int) bool {
	vStr1 := ""
	vStr2 := ""
	split1 := strings.Split(s[i], "_")
	if len(split1) > 1 {
		vStr1 = split1[1]
	}
	split2 := strings.Split(s[j], "_")
	if len(split2) > 1 {
		vStr2 = split2[1]
	}
	v1, err := version.NewVersion(vStr1)
	if err != nil {
		v1, _ = version.NewVersion("20150101.01")
	}
	v2, err := version.NewVersion(vStr2)
	if err != nil {
		v2, _ = version.NewVersion("20150101.01")
	}
	return v1.GreaterThan(v2)
}
