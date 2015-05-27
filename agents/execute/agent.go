package execute

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	"golang.org/x/net/context"
)

type executeAgent struct{}

func init() {
	arc.RegisterAgent("execute", new(executeAgent))
}

func (a *executeAgent) Enabled() bool { return true }

func (a *executeAgent) Enable() error { return nil }

func (a *executeAgent) Disable() error { return nil }

func (a *executeAgent) CommandAction(ctx context.Context, payload string, heartbeat func(string)) (string, error) {

	command := splitArgs(payload)
	if len(command) == 0 {
		return "", fmt.Errorf("Invalid payload. Command should by a string or array.")
	}

	process := arc.NewSubprocess(command[0], command[1:]...)

	output, err := process.Start()
	if err != nil {
		return "", err
	}
	//send empty heartbeat so that the caller knows the command is executing
	heartbeat("")

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
					heartbeat(line)
				default:
					return "", process.Error()
				}
			}
		case line := <-output:
			heartbeat(line)
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
