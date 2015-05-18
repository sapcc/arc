package transport

import (
	"errors"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/transport/mqtt"
)

type Transport interface {
	Connect()
	Disconnect()
	Request(*arc.Request)
	Reply(*arc.Reply)
	Subscribe() <-chan *arc.Request
	SubscribeJob(requestId string) <-chan *arc.Reply
}

func New(config arc.Config) (Transport, error) {
	switch config.Transport {
	case "mqtt":
		return mqtt.New(config)
	}
	return nil, errors.New("Invalid transport")
}
