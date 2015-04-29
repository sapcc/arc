package transport

import (
	"errors"
	"gitHub.***REMOVED***/monsoon/onos/transport/mqtt"
	"gitHub.***REMOVED***/monsoon/onos/types"
)

type Transport interface {
	Connect()
	Disconnect()
	Publish(types.Message)
	Subscribe(func(types.Message))
}

func New(config types.Config) (Transport, error) {
	switch config.Transport {
	case "mqtt":
		return mqtt.New(config)
	}
	return nil, errors.New("Invalid transport")
}
