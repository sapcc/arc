package arc

import (
	"gitHub.***REMOVED***/monsoon/arc/version"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"	
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
	facts["project"] = h.config.Project
	facts["identity"] = h.config.Identity
	facts["organization"] = h.config.Organization
	return facts, nil
}
