package execute

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	"golang.org/x/net/context"
)

type executeAgent struct{}

func init() {
	arc.RegisterAgent("execute", new(executeAgent))
}

func (a *executeAgent) Enabled() bool { return true }

func (a *executeAgent) Enable(ctx context.Context, job *arc.Job) (string, error) { return "", nil }

func (a *executeAgent) Disable(ctx context.Context, job *arc.Job) (string, error) { return "", nil }

func (a *executeAgent) CommandAction(ctx context.Context, job *arc.Job) (string, error) {

	command := splitArgs(job.Payload)
	if len(command) == 0 {
		return "", fmt.Errorf("invalid payload. Command should by a string or array")
	}

	process := arc.NewSubprocess(command[0], command[1:]...)

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

}

func (a *executeAgent) ScriptAction(ctx context.Context, job *arc.Job) (string, error) {
	if job.Payload == "" {
		return "", errors.New("empty payload")
	}

	file, err := ioutil.TempFile(os.TempDir(), "execute")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %s", err)
	}
	if _, err := file.WriteString(job.Payload); err != nil {
		return "", fmt.Errorf("failed to write script to temporary file: %s", err)
	}

	if err := file.Close(); err != nil {
		return "", fmt.Errorf("failed to close script file: %s", err)
	}

	script_name := file.Name() + scriptSuffix
	if err := os.Rename(file.Name(), script_name); err != nil {
		if removeErr := os.Remove(file.Name()); removeErr != nil {
			return "", fmt.Errorf("error removing file: %s. After getting error renaming file: %s", removeErr, err)
		}
		return "", err
	}
	defer os.Remove(script_name)

	process := scriptCommand(script_name)
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

}

func splitArgs(cmd string) []string {
	var args []string
	err := json.Unmarshal([]byte(cmd), &args)
	if err != nil {
		return strings.Fields(cmd)
	}
	return args
}
