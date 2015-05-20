package rpc

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type rpcAgent struct{}

func init() {
	arc.RegisterAgent("rpc", new(rpcAgent))
}

func (a *rpcAgent) Enabled() bool { return true }

func (a *rpcAgent) Enable() error { return nil }

func (a *rpcAgent) Disable() error { return nil }

func (a *rpcAgent) PingAction(ctx context.Context, payload string, heartbeat func(string)) (string, error) {
	return "pong", nil
}

func (a *rpcAgent) SleepAction(ctx context.Context, payload string, heartbeat func(string)) (string, error) {

	wait, err := strconv.Atoi(payload)
	if err != nil {
		return "", err
	}
	heartbeat("")
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Duration(wait) * time.Second):

	}
	return fmt.Sprintf("Slept for %d second(s)", wait), nil

}
