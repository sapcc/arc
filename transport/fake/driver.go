package fake

import (
	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
)

type FakeClient struct {
	Name      string
	Done      chan bool
	ReplyChan chan *arc.Reply
	RegChan   chan *arc.Registration
	ReqChan   chan *arc.Request
	Connected bool
}

func New(config arc_config.Config) (*FakeClient, error) {
	log.Infof("Using FAKE transport")
	
	// used to fake the connectivity
	isConnected := true
	if config.Organization == "no-connected" {
		isConnected = false
	}
	
	return &FakeClient{
			Name: "fake",
			Done: make(chan bool),
			Connected: isConnected,
		},
		nil
}

func (c *FakeClient) Connect() error {
	return nil
}

func (c *FakeClient) Disconnect() {
}

func (c *FakeClient) IsConnected() bool {
	return c.Connected
}

func (c *FakeClient) Subscribe(identity string) (<-chan *arc.Request, func()) {
	log.Infof("Subscribe with the FAKE transport")

	out := make(chan *arc.Request)
	c.ReqChan = out
	cancel := func() {
		log.Infof("FAKE request transport closed")
		close(out)
	}
	return out, cancel
}

func (c *FakeClient) Request(msg *arc.Request) error {
	go func() {
		log.Infof("Writing Request into the FAKE transport. %q", msg)
		c.ReqChan <- msg
	}()
	return nil
}

func (c *FakeClient) Reply(msg *arc.Reply) error {
	go func() {
		log.Infof("Writing Reply into the FAKE transport. %q", msg)
		c.ReplyChan <- msg
	}()
	return nil
}

func (c *FakeClient) SubscribeJob(requestId string) (<-chan *arc.Reply, func()) {
	return nil, nil
}

func (c *FakeClient) SubscribeReplies() (<-chan *arc.Reply, func()) {
	log.Infof("SubscribeReplies with the FAKE transport")

	out := make(chan *arc.Reply)
	c.ReplyChan = out
	cancel := func() {
		log.Infof("FAKE reply transport closed")
		close(out)
	}
	return out, cancel
}

func (c *FakeClient) Registration(msg *arc.Registration) error {
	go func() {
		log.Infof("Writing Request into the FAKE transport. %q", msg)
		c.RegChan <- msg
	}()
	return nil
}

func (c *FakeClient) SubscribeRegistrations() (<-chan *arc.Registration, func()) {
	log.Infof("SubscribeRegistrations with the FAKE transport")

	out := make(chan *arc.Registration)
	c.RegChan = out
	cancel := func() {
		log.Info("FAKE registration transport closed")
		close(out)
	}

	return out, cancel
}

func (c *FakeClient) DoneSignal() {
	c.Done <- true
}
