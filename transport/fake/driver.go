package fake

import (
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type FakeClient struct {
	Name string
}

func New(config arc.Config) (*FakeClient, error) {
	return &FakeClient{Name: "fake"}, nil
}

func (c *FakeClient) Connect() error {
	return nil
}

func (c *FakeClient) Disconnect() {
}

func (c *FakeClient) Subscribe(identity string) (<-chan *arc.Request, func()) {
	return nil, nil
}

func (c *FakeClient) Request(msg *arc.Request) {
}

func (c *FakeClient) Reply(msg *arc.Reply) {
}

func (c *FakeClient) SubscribeJob(requestId string) (<-chan *arc.Reply, func()) {
	return nil, nil
}

func (c *FakeClient) SubscribeReplies() (<-chan *arc.Reply, func()) {
	return nil, nil
}
