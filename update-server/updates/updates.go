package updates

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/blang/semver"
	"github.com/inconshreveable/go-update/check"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type availableUpdate struct {
	buildName string
	version   string
}

const buildRelativeUrl = "/builds/"

/*
 * return nil if no update available
 */
func New(req *http.Request, buildsRootPath string) *check.Result {
	// check required build path
	if len(buildsRootPath) == 0 {
		log.Errorf("Build path is missing. Got %q", buildsRootPath)
		return nil
	}

	// get host url
	hostUrl := getHostUrl(req)
	if hostUrl == nil {
		log.Errorf("Computed host url is nil. Request %q", req)
		return nil
	}

	// read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Error while reading the request body. Got %q", err)
		return nil
	}
	if len(body) == 0 {
		log.Errorf("No request body. Got %q", body)
		return nil
	}

	// convert to check.Params struc
	var reqParams check.Params
	err = json.Unmarshal(body, &reqParams)
	if err != nil {
		log.Errorf("Error unmarshaling body. Got %q", err)
		return nil
	}

	// check required post attributes
	if len(reqParams.AppId) == 0 {
		log.Errorf("Missing required post attribute 'app_id'. Got %q", reqParams.AppId)
		return nil
	}
	if len(reqParams.AppVersion) == 0 {
		log.Errorf("Missing required post attribute 'app_version'. Got %q", reqParams.AppVersion)
		return nil
	}
	if len(reqParams.Tags["os"]) == 0 {
		log.Errorf("Missing required post attribute 'tags[os]'. Got %q", reqParams.Tags["os"])
		return nil
	}
	if len(reqParams.Tags["arch"]) == 0 {
		log.Errorf("Missing required post attribute 'tags[arch]'. Got %q", reqParams.Tags["arch"])
		return nil
	}

	// get build url
	au := getAvailableUpdate(buildsRootPath, reqParams.AppId, reqParams.AppVersion, reqParams.Tags["os"], reqParams.Tags["arch"])
	if au != nil {
		return &check.Result{
			Initiative: "automatically",
			Url:        fmt.Sprint(hostUrl, buildRelativeUrl, au.buildName),
			Version:    au.version,
		}
	}

	return nil
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

func getAvailableUpdate(buildsRootPath string, appId string, appVersion string, appOs string, appArch string) *availableUpdate {
	av, err := semver.Make(appVersion)
	if err != nil {
		log.Errorf("Error parsing app version. Got %q", err.Error())
		return nil
	}

	buildFile := ""
	buildVersion := "0.0.0"
	builds, err := ioutil.ReadDir(buildsRootPath)
	if err != nil {
		log.Errorf("Error reading builds dir. Got %q", err.Error())
		return nil
	}

	for _, f := range builds {
		if strings.HasPrefix(f.Name(), fmt.Sprint(appId, "_", appOs, "_", appArch, "_")) {
			fileVersion := strings.Split(f.Name(), "_")[3]
			bv, err := semver.Make(fileVersion)
			if err != nil {
				log.Errorf("Error parsing build version. Got %q. With error %q", buildVersion, err.Error())
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
		}
	}

	return nil
}
