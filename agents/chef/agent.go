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
	clientRbTemplate    = template.Must(template.New("client.rb").Parse(`chef_repo_path '{{ .chefRepoPath }}'{{.eol -}}
recipe_url '{{.recipeURL}}'{{ .eol -}}
{{ if .nodeName -}}
node_name '{{.nodeName}}'{{ .eol -}}
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
	cmd := exec.Command(chefClientBinary, "-v") // #nosec
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

	return "", fmt.Errorf("enabling chef agent failed")
}

func (a *chefAgent) Disable(ctx context.Context, job *arc.Job) (string, error) { return "", nil }

func (a *chefAgent) ZeroAction(ctx context.Context, job *arc.Job) (string, error) {
	var data chefZeroPayload
	if err := json.Unmarshal([]byte(job.Payload), &data); err != nil {
		return "", fmt.Errorf("invalid json payload for zero action of chef agent: %s", err)
	}
	if data.RecipeURL == "" {
		return "", fmt.Errorf("recipe_url not given or invalid")
	}
	if data.RunList == nil {
		return "", fmt.Errorf("run_list not given or invalid")
	}
	if data.NodeName == "" {
		data.NodeName = job.Identity()
	}
	if data.Attributes == nil {
		//if no attributes given init the map
		data.Attributes = make(map[string]interface{})
	}
	data.Attributes["run_list"] = data.RunList
	dna, err := json.MarshalIndent(data.Attributes, " ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize dna.json: %s", err)
	}
	tmpDir, err := ioutil.TempDir("", "chef-zero")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %s", tmpDir)
	}
	//on windows ioutil.TempDir returns backslashes
	tmpDir = strings.Replace(tmpDir, `\`, "/", -1)

	log_level := "info"
	if data.Debug {
		//this might overwhelm the broker atm, therefore we stay at info for the moment
		//log_level = "debug"
	} else {
		//by default if we are not in debug mode we remove the temporary Directory
		defer os.RemoveAll(tmpDir)
	}

	dnaFile, err := os.Create(path.Join(tmpDir, "dna.json"))
	if err != nil {
		return "", fmt.Errorf("failed to create dna.json: %s", err)
	}
	defer dnaFile.Close()

	log.Infof("Writing dna.json to %s", dnaFile.Name())
	if _, err := dnaFile.Write(dna); err != nil {
		return "", fmt.Errorf("failed to write dna.json to disk: %s", err)
	}

	configFile, err := os.Create(path.Join(tmpDir, "client.rb"))
	if err != nil {
		return "", fmt.Errorf("failed to create client.rb: %s", err)
	}
	defer configFile.Close()

	configVars := map[string]string{
		"nodeName":     data.NodeName,
		"chefRepoPath": tmpDir,
		"recipeURL":    data.RecipeURL,
		"eol":          eol,
	}

	if err = clientRbTemplate.Execute(configFile, configVars); err != nil {
		return "", fmt.Errorf("failed to write client.rb to disk: %s", err)
	}

	installedVersion, err := version.NewVersion(chefVersion())
	if err != nil {
		return "", err
	}
	chefZeroMinimalVersion, err := version.NewVersion("12.1.0")
	if err != nil {
		return "", err
	}

	if data.Nodes != nil {
		var nodesDir = path.Join(tmpDir, "nodes")
		if err = os.Mkdir(nodesDir, 0755); /* #nosec */ err != nil {
			return "", fmt.Errorf("failed to create %s: %s", nodesDir, err)
		}
		for i, node := range data.Nodes {
			nodeJson, err := json.MarshalIndent(node, "", "  ")
			if err != nil {
				return "", fmt.Errorf("failed to marshal node %d: %s", i, err)
			}
			nodeName := strconv.Itoa(i)
			if name, ok := node["name"]; ok {
				if s, ok := name.(string); ok {
					nodeName = s
				}
			}
			nodeFile := path.Join(nodesDir, fmt.Sprintf("%s.json", nodeName))
			if err = ioutil.WriteFile(nodeFile, nodeJson, 0644); err != nil {
				return "", fmt.Errorf("failed to write %s: %s", nodeFile, err)
			}
		}
	}

	var process *arc.Subprocess
	if installedVersion.LessThan(chefZeroMinimalVersion) {
		log.Warnf("Detected chef version < %s, falling back to chef-solo", chefZeroMinimalVersion)
		process = arc.NewSubprocess(chefSoloBinary, "--no-fork", "-c", configFile.Name(), "-j", dnaFile.Name(), "--log_level", log_level)
	} else {
		process = arc.NewSubprocess(chefClientBinary, "--local-mode", "--no-fork", "-c", configFile.Name(), "-j", dnaFile.Name(), "--log_level", log_level)
	}
	log.Info("Running ", strings.Join(process.Command, " "))

	output, err := process.Start()
	if err != nil {
		return "", err
	}
	//send heartbeat so that the caller knows the command is executing
	job.Heartbeat(fmt.Sprintf("Running %s\nusing recipe tarball at %s\n", strings.Join(process.Command, " "), data.RecipeURL))

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

	// request metadata
	req, err := http.NewRequest("GET", metadata_url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %s", err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("requesting omnitruck metdata failed: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)
	}

	// omnitruck response
	var omnitruck omnitruckResponse
	err = json.Unmarshal(body, &omnitruck)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal omnitruck response: %s", err)
	}
	extension := regexp.MustCompile(`(rpm|deb|msi|dmg)$`).FindString(omnitruck.Url)
	if extension == "" {
		return "", fmt.Errorf("unknown package type: %s", omnitruck.Url)
	}

	// Download installer
	installerResp, err := http.Get(omnitruck.Url)
	if err != nil {
		return "", fmt.Errorf("failed to download installer: %s", err)
	}
	defer installerResp.Body.Close()
	file, err := arc.TempFile("", "chef-installer", "."+extension)
	if err != nil {
		return "", fmt.Errorf("failed to create tmpfile: %s", err)
	}
	defer file.Close()
	log.Infof("Downloading omnibus installer from %s to %s", omnitruck.Url, file.Name())
	if _, err := io.Copy(file, installerResp.Body); err != nil {
		return "", fmt.Errorf("failed to save download to file: %s", err)
	}
	log.Info("Download succeeded")
	return file.Name(), nil
}
