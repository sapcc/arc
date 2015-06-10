package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"io"
)

type Job struct {
	arc.Request `json:"request"`
	Status      string `json:"status"`
}

type Jobs []Job

type Status string

func CreateJob(data *io.ReadCloser) (*Job, error) {
	// unmarschall body to a request
	var tmpReq arc.Request
	decoder := json.NewDecoder(*data)
	err := decoder.Decode(&tmpReq)
	if err != nil {
		return nil, err
	}

	// validate request
	request, err := arc.CreateRequest(tmpReq.Agent, tmpReq.Action, tmpReq.To, tmpReq.Timeout, tmpReq.Payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		*request,
		arc.Queued.String(),
	}, nil
}

func SaveJob(db *sql.DB, job *Job) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	var lastInsertId string
	err := db.QueryRow(ownDb.InsertJobQuery, job.Version, job.Sender, job.RequestID, job.To, job.Timeout, job.Agent, job.Action, job.Payload, job.Status).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateJob(db *sql.DB, job *Job) error {
	return nil
}

func GetAllJobs(db *sql.DB) (*Jobs, error) {
	var jobs Jobs
	rows, err := db.Query(ownDb.GetAllJobsQuery)
	if err != nil {
		return nil, err
	}

	var job Job
	for rows.Next() {
		err = rows.Scan(&job.Version, &job.Sender, &job.RequestID, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status)
		if err != nil {
			log.Errorf("Error scaning job results. Got ", err.Error())
			continue
		}
		jobs = append(jobs, job)
	}

	return &jobs, nil
}

func GetJob(db *sql.DB, requestId string) (*Job, error) {
	var job Job
	err := db.QueryRow(ownDb.GetJobQuery, requestId).Scan(&job.Version, &job.Sender, &job.RequestID, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
