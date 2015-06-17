package rpc

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

type rpcAgent struct{}

func init() {
	arc.RegisterAgent("rpc", new(rpcAgent))
}

func (a *rpcAgent) Enabled() bool { return true }

func (a *rpcAgent) Enable(ctx context.Context, job *arc.Job) (string, error) { return "", nil }

func (a *rpcAgent) Disable(ctx context.Context, job *arc.Job) (string, error) { return "", nil }

func (a *rpcAgent) VersionAction(ctx context.Context, job *arc.Job) (string, error) {
	return version.String(), nil
}

func (a *rpcAgent) PingAction(ctx context.Context, job *arc.Job) (string, error) {
	return "pong", nil
}

func (a *rpcAgent) SleepAction(ctx context.Context, job *arc.Job) (string, error) {

	wait, err := strconv.Atoi(job.Payload)
	if err != nil {
		return "", err
	}
	job.Heartbeat("")
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Duration(wait) * time.Second):

	}
	return fmt.Sprintf("Slept for %d second(s)", wait), nil

}
