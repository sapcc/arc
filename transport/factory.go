package transport

import (
	"errors"
	"gitHub.***REMOVED***/monsoon/onos/onos"
	"gitHub.***REMOVED***/monsoon/onos/transport/mqtt"
)

type Transport interface {
	Connect()
	Disconnect()
	Publish(onos.Message)
	Subscribe(func(onos.Message))
}

func New(config onos.Config) (Transport, error) {
	switch config.Transport {
	case "mqtt":
		return mqtt.New(config)
	}
	return nil, errors.New("Invalid transport")
}
