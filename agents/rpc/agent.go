package rpc

import (
	"gitHub.***REMOVED***/monsoon/onos/onos"
)

type RpcAgent struct{}

func init() {
	onos.RegisterAgent("rpc", new(RpcAgent))
}

func (a *RpcAgent) Enabled() bool {
	return true
}
func (a *RpcAgent) Enable() error {
	return nil
}
func (a *RpcAgent) Disable() error {
	return nil
}

func (a *RpcAgent) PingAction(payload string) {
}
