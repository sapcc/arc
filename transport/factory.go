package transport

import (
	"gitHub.***REMOVED***/monsoon/onos/transport/mqtt"
)

type Transport interface {
	Connect()
	Disconnect()
}

func New() (Transport, error) {
	return mqtt.New()
}
