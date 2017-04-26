package arc

import (
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/version"
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

	if len(h.config.Project) > 0 {
		facts["project"] = h.config.Project
	}
	if len(h.config.Identity) > 0 {
		facts["identity"] = h.config.Identity
	}
	if len(h.config.Organization) > 0 {
		facts["organization"] = h.config.Organization
	}
	return facts, nil
}
