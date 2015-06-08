package models

import (
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"io"
	"encoding/json"
)

type Job struct {
	Request arc.Request `json:"request"`
	Status  Status `json:"status"`
}

type Jobs []Job

type Status string
const (
	Queued Status = "queued"
	Executing Status = "executing"
	Failed Status = "failed"
)

func CreateJob(data *io.ReadCloser) (*Job, error) {
	var job Job
	decoder := json.NewDecoder(*data)
	err := decoder.Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}