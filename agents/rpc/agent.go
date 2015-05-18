package rpc

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type RpcAgent struct{}

func init() {
	arc.RegisterAgent("rpc", new(RpcAgent))
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

//private

func (a *RpcAgent) PingAction(ctx context.Context, payload string) (string, error) {
	return "pong", nil
}

func (a *RpcAgent) SleepAction(ctx context.Context, payload string) (string, error) {

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
