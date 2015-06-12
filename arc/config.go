package arc

import (
	"fmt"
)

type Config struct {
	Endpoints  []string
	ClientCa   string
	ClientCert string
	ClientKey  string
	Transport  string
	Identity   string
	Project    string
	LogLevel   string
}

func (c Config) String() string {
	return fmt.Sprintf("Endpoints: %s, ClientCa: %s, ClientCert: %s, ClientKey: %s, Transport: %s, Identity: %s, Project: %s, LogLevel: %s", c.Endpoints, c.ClientCa, c.ClientCert, c.ClientKey, c.Transport, c.Identity, c.Project, c.LogLevel)
}
