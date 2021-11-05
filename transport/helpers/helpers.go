package helpers

import "time"

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

type DriverError struct {
	Err       error
	TimeStamp time.Time
}

type RevokedCertError struct {
	Msg string
}

func (e RevokedCertError) Error() string {
	return e.Msg
}
