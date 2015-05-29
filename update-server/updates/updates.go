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

type AvailableUpdate struct {
	buildName	string
	version string
}

var BuildsRootPath string
var BuildRelativeUrl string

/*
 * return nil if no update available
 */
func New(req *http.Request, buildsPath string, buildUrl string) *check.Result {
	// save statics path
	BuildsRootPath = buildsPath
	BuildRelativeUrl = buildUrl
	
	// get host url
	hostUrl := getHostUrl(req)
		
	// read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf(err.Error())
	}

	// convert to check.Params struc
	var reqParams check.Params
	err = json.Unmarshal(body, &reqParams)
	if err != nil {
		log.Errorf(err.Error())
	}

	// get build url
	availableUpdate := getAvailableUpdate(reqParams.AppId, reqParams.AppVersion, reqParams.Tags["os"], reqParams.Tags["arch"])

	if availableUpdate != nil {
		return &check.Result{
			Initiative: "automatically",
			Url:        fmt.Sprint(hostUrl, BuildRelativeUrl, availableUpdate.buildName),
			Version:    availableUpdate.version,
		}		
	}

	return nil
}

// private

func getHostUrl(req *http.Request) *url.URL{
	host := req.Host
	scheme := ""
	
	// set the scheme
	if req.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"		
	}
		
	return &url.URL{Scheme: scheme, Host: host}
}

func getAvailableUpdate(appId string, appVersion string, appOs string, appArch string) *AvailableUpdate {
	av, err := semver.Make(appVersion)
	if err != nil {
		log.Errorf("Error parsing app version. Got ", err)
		return nil
	}
	
	buildFile := ""
	buildVersion := "0.0.0"
	builds, _ := ioutil.ReadDir(BuildsRootPath)
	for _, f := range builds {		
		if strings.HasPrefix(f.Name(), fmt.Sprint(appId, "_", appOs, "_", appArch, "_")) {
			fileVersion := strings.Split(f.Name(), "_")[3]
			bv, err := semver.Make(fileVersion)
			if err != nil {
				log.Errorf("Error parsing build version. Got %q. With error %q", buildVersion, err)
			}
			
			nbv, _ := semver.Make(buildVersion)
			if bv.Compare(av) == 1 && bv.Compare(nbv) == 1{
				buildFile = f.Name()
				buildVersion = fileVersion	
			}
		}
	}
	
	if len(buildFile) > 0 {
		return &AvailableUpdate{
			buildName: buildFile, 
			version: buildVersion,
		} 
	}
	
	return nil
}
