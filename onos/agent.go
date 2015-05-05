package onos

import (
	log "github.com/Sirupsen/logrus"
	"reflect"
)

var (
	agentRegistry = make(map[string]Agent)
)

type Action func(Agent, payload string)

type Agent interface {
	Enabled() bool
	Enable() error
	Disable() error
}

type agent struct {
	actions map[string]Action
}

func RegisterAgent(name string, agent Agent) {
	agentType := reflect.TypeOf(agent)
	for i := 0; i < agentType.NumMethod(); i++ {
		method := agentType.Method(i)
		log.Info(method.Name)
		log.Info(method.Type)
	}
	agentRegistry[name] = agent
}
