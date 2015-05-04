package onos

type Config struct {
	Endpoints  []string `toml:"endpoints"`
	ClientCa   string   `toml:"client-cakeys"`
	ClientCert string   `toml:"client-cert"`
	ClientKey  string   `toml:"client-key"`
	ConfigDir  string   `toml:"config-dir"`
	Transport  string   `toml:"transport"`
	Identity   string   `toml:"ideentity"`
	Project    string   `toml:"project"`
}
