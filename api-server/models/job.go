package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type Job struct {
	arc.Request `json:"request"`
	Status      arc.JobState `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type JobID struct {
	RequestID string   `json:"request_id"`
}

type Jobs []Job

type Status string

func CreateJob(data *[]byte, identity string) (*Job, error) {
	var tmpReq arc.Request
	// unmarshal
	err := json.Unmarshal(*data, &tmpReq)
	if err != nil {
		return nil, err
	}

	// create a valid request
	request, err := arc.CreateRequest(tmpReq.Agent, tmpReq.Action, identity, tmpReq.To, tmpReq.Timeout, tmpReq.Payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		*request,
		arc.Queued,
		time.Now(),
		time.Now(),
	}, nil
}

func SaveJob(db *sql.DB, job *Job) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	var lastInsertId string
	err := db.QueryRow(ownDb.InsertJobQuery, job.Version, job.Sender, job.RequestID, job.To, job.Timeout, job.Agent, job.Action, job.Payload, job.Status, job.CreatedAt, job.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateJob(db *sql.DB, reply *arc.Reply) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	res, err := db.Exec(ownDb.UpdateJobQuery, reply.State, time.Now(), reply.RequestID)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("%v rows where updated with id %q", affect, reply.RequestID)

	return nil
}

func GetAllJobs(db *sql.DB) (*Jobs, error) {
	if db == nil {
		return nil, errors.New("Db is nil")
	}
		
	jobs := make(Jobs,0)
	rows, err := db.Query(ownDb.GetAllJobsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var job Job
	for rows.Next() {
		err = rows.Scan(&job.Version, &job.Sender, &job.RequestID, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt)
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
	err := db.QueryRow(ownDb.GetJobQuery, requestId).Scan(&job.Version, &job.Sender, &job.RequestID, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
