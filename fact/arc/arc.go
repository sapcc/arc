package arc

import (
	arc_config "github.com/sapcc/arc/config"
	"github.com/sapcc/arc/version"
)

type Source struct {
	config arc_config.Config
}

func New(config arc_config.Config) Source {
	return Source{config: config}
}

func (h Source) Name() string {
	return "arc"
}

func (h Source) Facts() (map[string]interface{}, error) {
	facts := make(map[string]interface{})

	facts["arc_version"] = version.String()

	// set online true when updating the facts
	// this should fix the problem when deploying the API and the broker at the same time and the agents
	// send the "online" message before the API is ready to accept incoming broker messages.
	facts["online"] = true

	if h.config.Project != "" {
		facts["project"] = h.config.Project
	}
	if h.config.Identity != "" {
		facts["identity"] = h.config.Identity
	}
	if h.config.Organization != "" {
		facts["organization"] = h.config.Organization
	}
	return facts, nil
}
