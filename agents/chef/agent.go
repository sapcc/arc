package chef

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"gitHub.***REMOVED***/monsoon/arc/arc"

	log "github.com/Sirupsen/logrus"
	version "github.com/hashicorp/go-version"
	"golang.org/x/net/context"
)

type chefAgent struct{}

var (
	defaultOmnitruckUrl = "https://www.chef.io/chef/metadata"
	clientRbTemplate    = template.Must(template.New("client.rb").Parse(`chef_repo_path '{{ .chefRepoPath }}'
{{ if .nodeName -}}
node_name '{{.nodeName}}'
{{ end -}}`))
)

type omnitruckResponse struct {
	Url     string
	Sha256  string
	Md5     string
	Relpath string
}

type chefZeroPayload struct {
	RunList    []string                 `json:"run_list"`
	RecipeURL  string                   `json:"recipe_url"`
	Attributes map[string]interface{}   `json:"attributes"`
	Debug      bool                     `json:"debug"`
	Nodes      []map[string]interface{} `json:"nodes"`
	NodeName   string                   `json:"name"`
}

type enableOptions struct {
	OmnitruckUrl string `json:"omnitruck_url"`
	ChefVersion  string `json:"chef_version"`
}

func init() {
	arc.RegisterAgent("chef", new(chefAgent))
}

func chefVersion() string {
	cmd := exec.Command(chefClientBinary, "-v")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	if v := regexp.MustCompile(`\d+\.\d+\.\d+`).Find(out); v != nil {
		return string(v)
	}
	return ""
}

func (a *chefAgent) Enabled() bool {
	if ver := chefVersion(); ver != "" {
		log.Debug("Detected chef version ", ver)
		return true
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

		var opts enableOptions
		if job.Payload != "" {
			if err := json.Unmarshal([]byte(job.Payload), &opts); err != nil {
				return "", err
			}
		}
		if opts.OmnitruckUrl == "" {
			opts.OmnitruckUrl = defaultOmnitruckUrl
		}
		if opts.ChefVersion == "" {
			opts.ChefVersion = "latest"
		}
		job.Heartbeat("Downloading chef installer via " + omnitruckUrl + "\n")
		platform := facts["platform"].(string)
		platform_version := facts["platform_version"].(string)
		if installer, err := downloadInstaller(opts.OmnitruckUrl, platform, platform_version, opts.ChefVersion); err == nil {
			job.Heartbeat("Installing " + installer + "\n")
			if err := install(installer); err != nil {
				return "", err
			}
			return "Agent enabled\n", nil
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
	dna, err := json.MarshalIndent(data.Attributes, " ", "  ")
	if err != nil {
		return "", fmt.Errorf("Failed to serialize dna.json: %s", err)
	}
	tmpDir, err := ioutil.TempDir("", "chef-zero")
	if err != nil {
		return "", fmt.Errorf("Failed to create temporary directory: %s", tmpDir)
	}

	log_level := "info"
	if data.Debug {
		log_level = "debug"
	} else {
		//by default if we are not in debug mode we remove the temporary Directory
		defer os.RemoveAll(tmpDir)
	}

	dnaFile, err := os.Create(path.Join(tmpDir, "dna.json"))
	if err != nil {
		return "", fmt.Errorf("Failed to create dna.json: %s", err)
	}
	log.Infof("Writing dna.json to %s", dnaFile.Name())
	if _, err := dnaFile.Write(dna); err != nil {
		dnaFile.Close()
		return "", fmt.Errorf("Failed to write dna.json to disk: %s", err)
	}
	dnaFile.Close()

	configFile, err := os.Create(path.Join(tmpDir, "client.rb"))
	if err != nil {
		return "", fmt.Errorf("Failed to create client.rb: %s", err)
	}

	configVars := map[string]string{
		"nodeName":     data.NodeName,
		"chefRepoPath": tmpDir,
	}

	if err = clientRbTemplate.Execute(configFile, configVars); err != nil {
		configFile.Close()
		return "", fmt.Errorf("Failed to write client.rb to disk: %s", err)
	}
	configFile.Close()

	installedVersion, err := version.NewVersion(chefVersion())
	chefZeroMinimalVersion, err := version.NewVersion("12.1.0")
	if err != nil {
		return "", err
	}

	if data.Nodes != nil {
		var nodesDir = path.Join(tmpDir, "nodes")
		if err = os.Mkdir(nodesDir, 0755); err != nil {
			return "", fmt.Errorf("Failed to create %s: %s", nodesDir, err)
		}
		for i, node := range data.Nodes {
			nodeJson, err := json.MarshalIndent(node, "", "  ")
			if err != nil {
				return "", fmt.Errorf("Failed to marshal node %d: %s", i, err)
			}
			nodeName := strconv.Itoa(i)
			if name, ok := node["name"]; ok {
				if s, ok := name.(string); ok {
					nodeName = s
				}
			}
			nodeFile := path.Join(nodesDir, fmt.Sprintf("%s.json", nodeName))
			if err = ioutil.WriteFile(nodeFile, nodeJson, 0644); err != nil {
				return "", fmt.Errorf("Failed to write %s: %s", nodeFile, err)
			}
		}
	}

	var process *arc.Subprocess
	if installedVersion.LessThan(chefZeroMinimalVersion) {
		log.Warnf("Detected chef version < %s, falling back to chef-solo", chefZeroMinimalVersion)
		process = arc.NewSubprocess(chefSoloBinary, "--no-fork", "--recipe-url", data.RecipeURL, "-c", configFile.Name(), "-j", dnaFile.Name(), "--log_level", log_level)
	} else {
		process = arc.NewSubprocess(chefClientBinary, "--local-mode", "--no-fork", "--recipe-url", data.RecipeURL, "-c", configFile.Name(), "-j", dnaFile.Name(), "--log_level", log_level)
	}
	log.Info("Running ", strings.Join(process.Command, " "))

	output, err := process.Start()
	if err != nil {
		return "", err
	}
	//send heartbeat so that the caller knows the command is executing
	job.Heartbeat("Running " + strings.Join(process.Command, " "))

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
					//fmt.Print(line)
					job.Heartbeat(line)
				default:
					return "", process.Error()
				}
			}
		case line := <-output:
			//fmt.Print(line)
			job.Heartbeat(line)
		}
	}

	return "", nil
}

func downloadInstaller(omnitruckUrl, platform, platform_version, chef_version string) (string, error) {
	if chef_version == "" {
		chef_version = "latest"
	}
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
