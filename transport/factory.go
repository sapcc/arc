package transport

import (
	"errors"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/transport/fake"
	"gitHub.***REMOVED***/monsoon/arc/transport/helpers"
	"gitHub.***REMOVED***/monsoon/arc/transport/mqtt"
)

type Transport interface {
	Connect() error
	Disconnect()
	IsConnected() bool
	Request(msg *arc.Request) error
	Registration(msg *arc.Registration) error
	Reply(msg *arc.Reply) error
	Subscribe(identity string) (messages <-chan *arc.Request, cancel func())
	SubscribeJob(requestId string) (messages <-chan *arc.Reply, cancel func())
	SubscribeReplies() (messages <-chan *arc.Reply, cancel func())
	SubscribeRegistrations() (messages <-chan *arc.Registration, cancel func())
	IdentityInformation() helpers.TransportIdentity
}

func New(config arc_config.Config, server bool) (Transport, error) {
	switch helpers.TransportType(config.Transport) {
	case helpers.MQTT:
		return mqtt.New(config, server)
	case helpers.Fake:
		return fake.New(config, server)
	}
	return nil, errors.New("Invalid transport")
}
