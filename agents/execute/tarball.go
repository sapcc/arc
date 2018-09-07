package execute

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/c4milo/unzipit"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"golang.org/x/net/context"
)

type tarballPayload struct {
	URL         string            `json:"url"`
	Path        string            `json:"path"`
	Arguments   []string          `json:"arguments"`
	Environment map[string]string `json:"environment"`
}

func (a *executeAgent) TarballAction(ctx context.Context, job *arc.Job) (string, error) {
	var data tarballPayload
	if err := json.Unmarshal([]byte(job.Payload), &data); err != nil {
		return "", fmt.Errorf("Invalid json payload for tarball action: %s", err)
	}

	//send empty heartbeat so that the caller knows the command is executing
	job.Heartbeat("")

	tmpDir, err := ioutil.TempDir("", "execute-tarball")
	if err != nil {
		return "", fmt.Errorf("Failed to create temporary directory: %s", tmpDir)
	}
	defer os.Remove(tmpDir)

	log.Info("Fetching ", data.URL)
	res, err := http.Get(data.URL)
	if err != nil {
		return "", err
	}
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("Failed to retrieve %s: %s", data.URL, res.Status)
	}
	defer res.Body.Close()
	_, err = unzipit.UnpackStream(res.Body, tmpDir)
	if err != nil {
		return "", err
	}

	// powershell scripts cannot run directy on the win instancen.
	var process *arc.Subprocess
	if runtime.GOOS == "windows" {
		process = arc.NewSubprocess("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "RemoteSigned", "-Command", "$ErrorActionPreference = 'Stop'; & "+path.Join(tmpDir, data.Path))
	} else {
		process = arc.NewSubprocess(path.Join(tmpDir, data.Path), data.Arguments...)
	}

	log.Info("Running ", strings.Join(process.Command, " "))
	process.Dir = tmpDir
	if data.Environment != nil && len(data.Environment) > 0 {
		envVars := os.Environ()
		for key, val := range data.Environment {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, val))
		}
		process.Env = envVars
	}
	output, err := process.Start()
	if err != nil {
		return "", err
	}
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
}
