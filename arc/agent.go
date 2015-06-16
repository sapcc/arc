package arc

import (
	"reflect"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type registry struct {
	agents map[string]*agentInfo
}

var (
	agentRegistry = &registry{agents: make(map[string]*agentInfo)}
)

type agentAction func(context.Context, string, func(string)) (string, error)

type agentInfo struct {
	agent   Agent
	actions map[string]agentAction
}

type Agent interface {
	Enabled() bool
	Enable() error
	Disable() error
}

type Registry interface {
	HasAction(agent string, action string) bool
	Agents() []string
	Actions(agent string) []string
}

func AgentRegistry() Registry {
	return agentRegistry
}

func RegisterAgent(name string, agent Agent) {

	agentType := reflect.TypeOf(agent)
	actionMap := make(map[string]agentAction)

	re := regexp.MustCompile("^([A-Z].*)Action$")

	for i := 0; i < agentType.NumMethod(); i++ {
		method := agentType.Method(i)
		if match := re.FindStringSubmatch(method.Name); match != nil {
			action := strings.ToLower(match[1])
			actionFunction, ok := reflect.ValueOf(agent).MethodByName(method.Name).Interface().(func(context.Context, string, func(string)) (string, error))
			if ok {
				actionMap[action] = actionFunction
			} else {
				log.Warnf("Ignoring %s.%s, invalid function signature.", name, method.Name)
			}
		}
	}

	agentRegistry.agents[name] = &agentInfo{agent, actionMap}
}

func ExecuteAction(ctx context.Context, identity string, request *Request, out chan<- *Reply) {
	defer close(out)
	agt := agentRegistry.agents[request.Agent]
	sequence := 0
	reply_number := func() int {
		sequence++
		return sequence
	}

	if agt == nil {
		out <- CreateReply(request, identity, Failed, "Agent not found", reply_number())
		return
	}
	if agt.agent.Enabled() == false {
		out <- CreateReply(request, identity, Failed, "Agent not enabled", reply_number())
		return
	}
	if _, exists := agt.actions[request.Action]; !exists {
		out <- CreateReply(request, identity, Failed, "Action not found", reply_number())
		return
	}
	hearbeat := func(payload string) {
		out <- CreateReply(request, identity, Executing, payload, reply_number())
	}

	result, err := agt.executeAction(ctx, request.Action, request.Payload, hearbeat)
	if err != nil {
		out <- CreateReply(request, identity, Failed, err.Error(), reply_number())
	} else {
		out <- CreateReply(request, identity, Complete, result, reply_number())
	}

}

func (a *agentInfo) executeAction(ctx context.Context, action, payload string, hearbeat func(string)) (string, error) {
	return a.actions[action](ctx, payload, hearbeat)
}

func (r *registry) Agents() []string {

	agents := make([]string, 0, len(r.agents))
	for i, _ := range r.agents {
		agents = append(agents, i)
	}
	return agents

}

func (r *registry) Actions(agent string) []string {
	if agt, found := r.agents[agent]; found {
		actions := make([]string, 0, len(agt.actions))
		for i, _ := range agt.actions {
			actions = append(actions, i)
		}
		return actions
	}
	return nil
}

func (r *registry) HasAction(agent string, action string) bool {
	if a := r.agents[agent]; a != nil {
		_, found := a.actions[action]
		return found
	}
	return false
}
