package arc

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
