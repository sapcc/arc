package transport

import (
	"errors"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/transport/fake"
	"gitHub.***REMOVED***/monsoon/arc/transport/mqtt"
)

type Transport interface {
	Connect() error
	Disconnect()
	Request(msg *arc.Request) error
	Registration(msg *arc.Registration) error
	Reply(msg *arc.Reply) error
	Subscribe(identity string) (messages <-chan *arc.Request, cancel func())
	SubscribeJob(requestId string) (messages <-chan *arc.Reply, cancel func())
	SubscribeReplies() (messages <-chan *arc.Reply, cancel func())
	SubscribeRegistrations() (messages <-chan *arc.Registration, cancel func())
}

func New(config arc.Config) (Transport, error) {
	switch config.Transport {
	case "mqtt":
		return mqtt.New(config)
	case "fake":
		return fake.New(config)
	}
	return nil, errors.New("Invalid transport")
}
