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
	RequestID string `json:"request_id"`
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

func (jobs *Jobs) Get(db *sql.DB) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	*jobs = make(Jobs, 0)
	rows, err := db.Query(ownDb.GetAllJobsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var job Job
	for rows.Next() {
		err = rows.Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			log.Errorf("Error scaning job results. Got ", err.Error())
			continue
		}
		*jobs = append(*jobs, job)
	}

	rows.Close()
	return nil
}

func (job *Job) Get(db *sql.DB) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	err := db.QueryRow(ownDb.GetJobQuery, job.RequestID).Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (job *Job) Save(db *sql.DB) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	var lastInsertId string
	err := db.QueryRow(ownDb.InsertJobQuery, job.RequestID, job.Version, job.Sender, job.To, job.Timeout, job.Agent, job.Action, job.Payload, job.Status, job.CreatedAt, job.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return nil
}

func (job *Job) Update(db *sql.DB) (err error) {
	if db == nil {
		return errors.New("Db is nil")
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// update job
	res, err := tx.Exec(ownDb.UpdateJobQuery, job.Status, job.UpdatedAt, job.RequestID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	// update object data
	if err = tx.QueryRow(ownDb.GetJobQuery, job.RequestID).Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt); err != nil {
		return
	}

	log.Infof("%v rows where updated with id %q", affect, job.RequestID)

	return
}

func CleanJobs(db *sql.DB) error {
	if db == nil {
		return errors.New("Clean job: Db connection is nil")
	}

	// clean jobs which no heartbeat was send back after created_at + 60 sec
	res, err := db.Exec(ownDb.CleanJobsNonHeartbeatQuery, 60)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("Clean job: %v jobs without heartbeat answer where updated", affect)

	// clean jobs which the timeout + 60 sec has exceeded and still in queued or executing status
	res, err = db.Exec(ownDb.CleanJobsTimeoutQuery, 60)
	if err != nil {
		return err
	}

	affect, err = res.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("Clean job: %v timeout jobs where updated", affect)

	return nil
}
