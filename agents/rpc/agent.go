package rpc

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gitHub.***REMOVED***/monsoon/onos/onos"
	"golang.org/x/net/context"
)

type RpcAgent struct{}

func init() {
	onos.RegisterAgent("rpc", []string{"ping", "sleep"}, new(RpcAgent))
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

func (a *RpcAgent) Execute(ctx context.Context, action string, payload string) (string, error) {
	switch action {
	case "ping":
		return a.ping(ctx, payload)
	case "sleep":
		return a.sleep(ctx, payload)
	}
	return "", errors.New("Unknown Action")
}

//private

func (a *RpcAgent) ping(ctx context.Context, payload string) (string, error) {
	return "pong", nil
}

func (a *RpcAgent) sleep(ctx context.Context, payload string) (string, error) {

	wait, err := strconv.Atoi(payload)
	if err != nil {
		return "", err
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Duration(wait) * time.Second):

	}
	return fmt.Sprintf("Slept for %d second(s)", wait), nil

}
