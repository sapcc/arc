package arc

import (
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

type Source struct {
	config arc.Config
}

func New(config arc.Config) Source {
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
