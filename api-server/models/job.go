package models

import ()

type Job struct {
	ReqID   string `json:"req_id"`
	Payload string `json:"payload"`
	Status  Status `json:"status"`
}

type Jobs []Job

type Status string
const (
	Queued Status = "queued"
	Executing Status = "executing"
	Failed Status = "failed"
)
