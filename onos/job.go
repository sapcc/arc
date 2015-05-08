package onos

import (
	"golang.org/x/net/context"
)

type key int

const jobIDKey key = 0

func NewJobContext(ctx context.Context, jobID string) context.Context {
	return context.WithValue(ctx, jobIDKey, jobID)
}

func JobFromContext(ctx context.Context) (string, bool) {
	jobID, ok := ctx.Value(jobIDKey).(string)
	return jobID, ok
}
