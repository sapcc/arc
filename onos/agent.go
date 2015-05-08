package onos

import (
	"golang.org/x/net/context"
)

var (
	agentRegistry = make(map[string]*agentInfo)
)

type Agent interface {
	Enabled() bool
	Enable() error
	Disable() error
	Execute(ctx context.Context, action string, payload string) (string, error)
}

type agentInfo struct {
	agent   Agent
	actions map[string]bool
}

func RegisterAgent(name string, actions []string, agent Agent) {
	actionMap := make(map[string]bool)
	for _, a := range actions {
		actionMap[a] = true
	}
	agentRegistry[name] = &agentInfo{agent, actionMap}
}

func ExecuteAction(ctx context.Context, request *Message, out chan<- *Message) {
	defer close(out)
	agt := agentRegistry[request.Agent]
	if agt == nil {
		out <- CreateReply(request, "Agent not found")
		return
	}
	if agt.agent.Enabled() == false {
		out <- CreateReply(request, "Agent not enabled")
		return
	}
	if _, exists := agt.actions[request.Action]; !exists {
		out <- CreateReply(request, "Action not found")
		return
	}

	result, err := agt.agent.Execute(ctx, request.Action, request.Payload)
	if err != nil {
		out <- CreateReply(request, err.Error())
	} else {
		out <- CreateReply(request, result)
	}

}
