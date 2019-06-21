package agents

import (
	"github.com/sapcc/arc/arc"
)

type Source struct{}

func New() Source {
	return Source{}
}

func (h Source) Name() string {
	return "agents"
}

func (h Source) Facts() (map[string]interface{}, error) {

	agents := make(map[string]string)
	registry := arc.AgentRegistry()
	for _, agent := range registry.Agents() {
		if registry.IsEnabled(agent) {
			agents[agent] = "enabled"
		} else {
			agents[agent] = "disabled"
		}
	}
	facts := make(map[string]interface{})
	facts["agents"] = agents

	return facts, nil
}
