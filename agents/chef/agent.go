package chef

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"

	"gitHub.***REMOVED***/monsoon/arc/arc"

	log "github.com/Sirupsen/logrus"
	"github.com/blang/semver"
	"golang.org/x/net/context"
)

type chefAgent struct{}

var (
	defaultOmnitruckUrl = "https://www.chef.io/chef/metadata"
	chef_version        = "12.3.0"
)

type omnitruckResponse struct {
	Url     string
	Sha256  string
	Md5     string
	Relpath string
}

type chefZeroPayload struct {
	RunList    []string               `json:"run_list"`
	RecipeURL  string                 `json:"recipe_url"`
	Attributes map[string]interface{} `json:"attributes"`
}

func init() {
	arc.RegisterAgent("chef", new(chefAgent))
}

func (a *chefAgent) Enabled() bool {
	cmd := exec.Command(chef_binary, "-v")
	out, err := cmd.Output()
	if err != nil {
		log.Warn("Chef not installed: ", err)
		return false
	}
	if v := regexp.MustCompile(`\d+\.\d+\.\d+`).Find(out); v != nil {
		version := string(v)
		log.Info("Detected chef version ", version)
		installed_version, _ := semver.Make(version)
		wanted_version, _ := semver.Make(chef_version)
		if installed_version.GTE(wanted_version) {
			return true
		} else {
			log.Infof("Installed version (%s) of chef is outdated", installed_version)
		}
	}

	return false
}

func (a *chefAgent) Enable(ctx context.Context, job *arc.Job) (string, error) {
	//if a.Enabled() {
	//  return "Already installed", nil
	//}

	facts, _ := arc.FactsFromContext(ctx)
	if facts["platform"] != nil && facts["platform_version"] != nil {
		omnitruckUrl := defaultOmnitruckUrl
		if job.Payload != "" {
			omnitruckUrl = job.Payload
		}
		job.Heartbeat("")
		platform := facts["platform"].(string)
		platform_version := facts["platform_version"].(string)
		if installer, err := downloadInstaller(omnitruckUrl, platform, platform_version); err == nil {
			if err := install(installer); err != nil {
				return "", err
			}
			return "Agent enabled", nil
		} else {
			return "", err
		}

	}

	return "", fmt.Errorf("Enabling chef agent failed")
}

func (a *chefAgent) Disable(ctx context.Context, job *arc.Job) (string, error) { return "", nil }

func (a *chefAgent) ZeroAction(ctx context.Context, job *arc.Job) (string, error) {
	var data chefZeroPayload
	if err := json.Unmarshal([]byte(job.Payload), &data); err != nil {
		return "", fmt.Errorf("Invalid json payload for zero action of chef agent: %s", err)
	}
	if data.RecipeURL == "" {
		return "", fmt.Errorf("recipe_url not given or invalid")
	}
	if data.RunList == nil {
		return "", fmt.Errorf("run_list not given or invalid")
	}
	if data.Attributes == nil {
		//if no attributes given init the map
		data.Attributes = make(map[string]interface{})
	}
	data.Attributes["run_list"] = data.RunList
	dna, err := json.Marshal(data.Attributes)
	if err != nil {
		return "", fmt.Errorf("Failed to serialize dna.json: %s", err)
	}
	tempfile, _ := arc.TempFile("", "dna", ".json")
	log.Infof("Writing dna.json to %s", tempfile.Name())
	if _, err := tempfile.Write(dna); err != nil {
		return "", fmt.Errorf("Failed to write dna.json to disk: %s", err)
	}
	tempfile.Close()
	//defer os.Remove(tempfile.Name())

	process := arc.NewSubprocess(chef_binary, "--local-mode", "--recipe-url="+data.RecipeURL, "-j", tempfile.Name())

	output, err := process.Start()
	if err != nil {
		return "", err
	}
	//send empty heartbeat so that the caller knows the command is executing
	job.Heartbeat("")

	for {
		select {
		case <-ctx.Done():
			//The context was cancelled, stop the process
			process.Kill()
		case <-process.Done():
			//drain the output channel before quitting
			for {
				select {
				case line := <-output:
					job.Heartbeat(line)
				default:
					return "", process.Error()
				}
			}
		case line := <-output:
			job.Heartbeat(line)
		}
	}

	return "", nil
}

func downloadInstaller(omnitruckUrl, platform, platform_version string) (string, error) {
	metadata_url := fmt.Sprintf("%s?v=%s&p=%s&pv=%s&m=x86_64", omnitruckUrl, chef_version, platform, platform_version)
	log.Infof("Fetching Omnitruck metadata from %s", metadata_url)
	var client http.Client
	req, err := http.NewRequest("GET", metadata_url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Requesting omnitruck metdata failed: %s", err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("Failed to read response body", err)
		}
		var omnitruck omnitruckResponse
		if err := json.Unmarshal(body, &omnitruck); err == nil {
			extension := regexp.MustCompile(`(rpm|deb|msi|dmg)$`).FindString(omnitruck.Url)
			if extension == "" {
				return "", fmt.Errorf("Unknown package type: %s", omnitruck.Url)
			}
			if resp, err := http.Get(omnitruck.Url); err == nil {
				defer resp.Body.Close()
				file, _ := arc.TempFile("", "chef-installer", "."+extension)
				log.Infof("Downloading omnibus installer from %s to %s", omnitruck.Url, file.Name())
				if _, err := io.Copy(file, resp.Body); err != nil {
					return "", fmt.Errorf("Failed to save download to file: %s", err)
				}
				file.Close()
				log.Info("Download succeeded")
				return file.Name(), nil
			} else {
				return "", fmt.Errorf("Failed to download installer: %s", err)
			}
		} else {
			return "", fmt.Errorf("Failed to unmarshal omnitruck response", err)
		}

	}

}
