package onos

import (
	log "github.com/Sirupsen/logrus"
	"reflect"
)

var (
	agentRegistry = make(map[string]Agent)
)

type Agent interface {
	Enabled() bool
	Enable() error
	Disable() error
}

func RegisterAgent(name string, agent Agent) {
	agentType := reflect.TypeOf(agent)
	for i := 0; i < agentType.NumMethod(); i++ {
		method := agentType.Method(i)
		log.Info(method.Name)
	}
	agentRegistry[name] = agent
}
