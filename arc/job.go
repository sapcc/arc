package arc

import (
	"encoding/json"
	"fmt"
	"strings"

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

type JobState byte

const (
	_ JobState = iota
	Queued
	Executing
	Failed
	Complete
)

var jobStateStringMap = map[JobState]string{Queued: "queued", Executing: "executing", Failed: "failed", Complete: "complete"}

func (j *JobState) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("state should be a string, got %s", data)
	}
	str = strings.ToLower(str)
	for key, val := range jobStateStringMap {
		if val == str {
			*j = key
			return nil
		}
	}
	return fmt.Errorf("Invalid job state: %q", str)
}

func (j JobState) MarshalJSON() ([]byte, error) {
	got, ok := jobStateStringMap[j]
	if !ok {
		return nil, fmt.Errorf("Invalid job state: %q", j)
	}
	return json.Marshal(got)
}

func (j JobState) String() string {
	return jobStateStringMap[j]
}
