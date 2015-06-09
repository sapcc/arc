package models

import (
	"encoding/json"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gopkg.in/gorp.v1"
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
	var job Job
	decoder := json.NewDecoder(*data)
	err := decoder.Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func SaveJob(dbmap *gorp.DbMap, job *Job) error {
	if dbmap == nil {		
		return errors.New("Db mapper is nil")
	}
	
	err := dbmap.Insert(job)
	if err != nil {
		return err
	}
	
	return nil
}

func UpdateJob(dbmap *gorp.DbMap, job *Job) error {
  _, err := dbmap.Update(&job)
	if err != nil {
		return err
	}
	return nil
}

func GetAllJobs(dbmap *gorp.DbMap) (*Jobs, error) {
	var jobs Jobs
	_, err := dbmap.Select(&jobs, "select * from jobs order by requestid")
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func GetJob(dbmap *gorp.DbMap, requestId string) (*Job, error) {
	var job Job
	err := dbmap.SelectOne(&job, "select * from jobs where requestid=$1", requestId)
	if err != nil {
		return nil, err
	}
	return &job, nil
}