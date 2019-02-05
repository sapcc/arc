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

type agentAction func(context.Context, *Job) (string, error)

type agentInfo struct {
	agent   Agent
	actions map[string]agentAction
}

type Agent interface {
	Enabled() bool
	Enable(context.Context, *Job) (string, error)
	Disable(context.Context, *Job) (string, error)
}

type Registry interface {
	HasAction(agent string, action string) bool
	Agents() []string
	Actions(agent string) []string
	IsEnabled(agent string) bool
}

func AgentRegistry() Registry {
	return agentRegistry
}

func RegisterAgent(name string, agent Agent) {

	agentType := reflect.TypeOf(agent)
	actionMap := make(map[string]agentAction)
	actionMap["enable"] = agent.Enable
	actionMap["disable"] = agent.Disable

	re := regexp.MustCompile("^([A-Z].*)Action$")

	for i := 0; i < agentType.NumMethod(); i++ {
		method := agentType.Method(i)
		if match := re.FindStringSubmatch(method.Name); match != nil {
			action := strings.ToLower(match[1])
			actionFunction, ok := reflect.ValueOf(agent).MethodByName(method.Name).Interface().(func(context.Context, *Job) (string, error))
			if ok {
				actionMap[action] = actionFunction
			} else {
				log.Warnf("Ignoring %s.%s, invalid function signature.", name, method.Name)
			}
		}
	}

	agentRegistry.agents[name] = &agentInfo{agent, actionMap}
}

func ExecuteAction(ctx context.Context, job *Job) {
	agt := agentRegistry.agents[job.Agent]

	if agt == nil {
		job.Fail("Agent not found\n")
		return
	}
	if job.Action != "enable" && job.Action != "disable" && !agt.agent.Enabled() {
		job.Fail("Agent not enabled\n")
		return
	}
	if _, exists := agt.actions[job.Action]; !exists {
		job.Fail("Action not found\n")
		return
	}

	result, err := agt.executeAction(ctx, job)
	if err != nil {
		job.Fail(err.Error())
	} else {
		job.Complete(result)
	}

}

func (a *agentInfo) executeAction(ctx context.Context, job *Job) (string, error) {
	return a.actions[job.Action](ctx, job)
}

func (r *registry) Agents() []string {

	agents := make([]string, 0, len(r.agents))
	for i := range r.agents {
		agents = append(agents, i)
	}
	return agents

}

func (r *registry) Actions(agent string) []string {
	if agt, found := r.agents[agent]; found {
		actions := make([]string, 0, len(agt.actions))
		for i := range agt.actions {
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

func (r *registry) IsEnabled(agent string) bool {
	if a := r.agents[agent]; a != nil {
		return a.agent.Enabled()
	}
	return false
}
