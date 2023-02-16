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
