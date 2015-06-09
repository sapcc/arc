package models

import (
	"encoding/json"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"database/sql"
	log "github.com/Sirupsen/logrus"	
	"io"
	"errors"
)

type Job struct {
	arc.Request `json:"request"`
	Status      string `json:"status"`
}

type Jobs []Job

type Status string

const (
	Queued    Status = "queued"
	Executing Status = "executing"
	Failed    Status = "failed"
)

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
		string(Queued),
	}, nil
}

func SaveJob(db *sql.DB, job *Job) error {
	if db == nil {
		return errors.New("Db is nil")
	}
	
	var lastInsertId string
	err := db.QueryRow(`INSERT INTO jobs(version,sender,requestid,"to",timeout,agent,action,payload,status) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) returning requestid;`, 
		job.Version, job.Sender, job.RequestID, job.To, job.Timeout, job.Agent, job.Action, job.Payload, job.Status).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	
	return nil
}

func UpdateJob(db *sql.DB, job *Job) error {
  /*_, err := dbmap.Update(&job)
	if err != nil {
		return err
	}*/
	return nil
}

func GetAllJobs(db *sql.DB) (*Jobs, error) {
	var jobs Jobs
	rows, err := db.Query("SELECT * FROM jobs order by requestid")
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
	err := db.QueryRow("SELECT * FROM jobs WHERE requestid=$1", requestId).Scan(&job.Version, &job.Sender, &job.RequestID, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status)
	if err != nil {
		return nil, err
	}
	return &job, nil
}