package execute

import (
	"encoding/json"
	"fmt"
	"os/exec"
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

	path, err := exec.LookPath(command[0])
	if err != nil {
		return "", fmt.Errorf("Command %s not found.", command[0])
	}

	exec.Command(command[0], command[1:]...)

	return path, nil

}

func splitArgs(cmd string) []string {
	var args []string
	err := json.Unmarshal([]byte(cmd), &args)
	if err != nil {
		return strings.Fields(cmd)
	}
	return args
}
