package helpers

type TransportType string

const (
	MQTT TransportType = "mqtt"
	Fake TransportType = "fake"
)

type TransportIdentity struct {
	Identity     string
	Project      string
	Organization string
	Transport    TransportType
}
