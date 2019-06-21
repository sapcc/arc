package transport

import (
	"errors"

	"github.com/sapcc/arc/arc"
	arc_config "github.com/sapcc/arc/config"
	"github.com/sapcc/arc/transport/fake"
	"github.com/sapcc/arc/transport/helpers"
	"github.com/sapcc/arc/transport/mqtt"
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
	return nil, errors.New("invalid transport")
}
