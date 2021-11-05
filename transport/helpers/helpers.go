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

type RevokedCertError struct {
	Msg string
}

func (e RevokedCertError) Error() string {
	return e.Msg
}
