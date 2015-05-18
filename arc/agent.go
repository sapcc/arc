package arc

import (
	"reflect"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	agentRegistry = make(map[string]*agentInfo)
)

type Agent interface {
	Enabled() bool
	Enable() error
	Disable() error
}

type AgentAction func(context.Context, string) (string, error)

type agentInfo struct {
	agent   Agent
	actions map[string]AgentAction
}

func RegisterAgent(name string, agent Agent) {

	agentType := reflect.TypeOf(agent)
	actionMap := make(map[string]AgentAction)

	re := regexp.MustCompile("^([A-Z].*)Action$")

	for i := 0; i < agentType.NumMethod(); i++ {
		method := agentType.Method(i)
		if match := re.FindStringSubmatch(method.Name); match != nil {
			action := strings.ToLower(match[1])
			actionFunction, ok := reflect.ValueOf(agent).MethodByName(method.Name).Interface().(func(context.Context, string) (string, error))
			if ok {
				actionMap[action] = actionFunction
			} else {
				log.Warnf("Ignoring %s.%s, invalid function signature.", name, method.Name)
			}
		}
	}

	agentRegistry[name] = &agentInfo{agent, actionMap}
}

func ExecuteAction(ctx context.Context, request *Request, out chan<- *Reply) {
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

	result, err := agt.executeAction(ctx, request.Action, request.Payload)
	if err != nil {
		out <- CreateReply(request, err.Error())
	} else {
		out <- CreateReply(request, result)
	}

}

func (a *agentInfo) executeAction(ctx context.Context, action, payload string) (string, error) {
	return a.actions[action](ctx, payload)
}
