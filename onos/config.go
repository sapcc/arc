package onos

type Config struct {
	Endpoints  []string
	ClientCa   string
	ClientCert string
	ClientKey  string
	ConfigDir  string
	Transport  string
	Identity   string
	Project    string
	LogLevel   string
}
