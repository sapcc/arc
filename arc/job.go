package arc

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

type contextKey int

const jobIDKey contextKey = 0

func NewJobContext(ctx context.Context, jobID string) context.Context {
	return context.WithValue(ctx, jobIDKey, jobID)
}

func JobFromContext(ctx context.Context) (string, bool) {
	jobID, ok := ctx.Value(jobIDKey).(string)
	return jobID, ok
}

type Job struct {
	Jid            string
	Payload        string
	Agent          string
	Action         string
	out            chan<- *Reply
	reply_sequence uint
	request        *Request
	identity       string
}

func NewJob(identity string, request *Request, out chan<- *Reply) *Job {

	return &Job{
		Jid:            request.RequestID,
		Payload:        request.Payload,
		Agent:          request.Agent,
		Action:         request.Action,
		request:        request,
		out:            out,
		reply_sequence: 0,
		identity:       identity,
	}
}

func (j *Job) Heartbeat(payload string) {
	j.out <- CreateReply(j.request, j.identity, Executing, payload, j.reply_number())
}

func (j *Job) Fail(payload string) {
	j.out <- CreateReply(j.request, j.identity, Failed, payload, j.reply_number())
	close(j.out)
}

func (j *Job) Complete(payload string) {
	j.out <- CreateReply(j.request, j.identity, Complete, payload, j.reply_number())
	close(j.out)
}

func (j *Job) reply_number() uint {
	j.reply_sequence++
	return j.reply_sequence
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
