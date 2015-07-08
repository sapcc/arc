package fake

import (
	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type FakeClient struct {
	Name      string
	Done      chan bool
	ReplyChan chan *arc.Reply
	ReqChan   chan *arc.Request
}

func New(config arc.Config) (*FakeClient, error) {
	log.Infof("Using FAKE transport")
	return &FakeClient{
			Name: "fake",
			Done: make(chan bool)},
		nil
}

func (c *FakeClient) Connect() error {
	return nil
}

func (c *FakeClient) Disconnect() {
}

func (c *FakeClient) Subscribe(identity string) (<-chan *arc.Request, func()) {
	log.Infof("Subscribe with the FAKE transport")

	out := make(chan *arc.Request)
	c.ReqChan = out
	cancel := func() {
		log.Infof("FAKE transport closed")
		close(out)
	}
	return out, cancel
}

func (c *FakeClient) Request(msg *arc.Request) {
	go func() {
		log.Infof("Writing Request into the FAKE transport. %q", msg)
		c.ReqChan <- msg
	}()
}

func (c *FakeClient) Reply(msg *arc.Reply) {
	go func() {
		log.Infof("Writing Reply into the FAKE transport. %q", msg)
		c.ReplyChan <- msg
	}()
}

func (c *FakeClient) SubscribeJob(requestId string) (<-chan *arc.Reply, func()) {
	return nil, nil
}

func (c *FakeClient) SubscribeReplies() (<-chan *arc.Reply, func()) {
	log.Infof("SubscribeReplies with the FAKE transport")

	out := make(chan *arc.Reply)
	c.ReplyChan = out
	cancel := func() {
		log.Infof("FAKE transport closed")
		close(out)
	}
	return out, cancel
}

func (c *FakeClient) Registration(msg *arc.Registration) {
}

func (c *FakeClient) SubscribeRegistrations() (<-chan *arc.Registration, func()) {
	out := make(chan *arc.Registration)

	cancel := func() {
		log.Info("FAKE transport closed")
		close(out)
	}

	return out, cancel
}

func (c *FakeClient) DoneSignal() {
	c.Done <- true
}
