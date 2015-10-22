package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

var JobTargetAgentNotFoundError = fmt.Errorf("Target agent where the job has to be executed not found.")
var JobBadRequestError = fmt.Errorf("Error unmarschaling or creating/validating the arc request.")

type Job struct {
	arc.Request
	Status    arc.JobState `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Project   string       `json:"project"`
}

type JobID struct {
	RequestID string `json:"request_id"`
}

type Jobs []Job

type Status string

func CreateJob(db *sql.DB, data *[]byte, identity string) (*Job, error) {
	if db == nil {
		return nil, errors.New("Db is nil")
	}

	// unmarshal data
	var tmpReq arc.Request
	err := json.Unmarshal(*data, &tmpReq)
	if err != nil {
		return nil, JobBadRequestError
	}

	// create a validate request
	request, err := arc.CreateRequest(tmpReq.Agent, tmpReq.Action, identity, tmpReq.To, tmpReq.Timeout, tmpReq.Payload)
	if err != nil {
		return nil, JobBadRequestError
	}

	// check and get the target project id
	agent := Agent{AgentID: tmpReq.To}
	err = agent.Get(db)
	if err == sql.ErrNoRows {
		return nil, JobTargetAgentNotFoundError
	} else if err != nil {
		return nil, err
	}

	return &Job{
		*request,
		arc.Queued,
		time.Now(),
		time.Now(),
		agent.Project,
	}, nil
}

func CreateJobAuthorized(db *sql.DB, data *[]byte, identity string, authorization *auth.Authorization) (*Job, error) {
	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return nil, err
	}

	job, err := CreateJob(db, data, identity)
	if err != nil {
		return nil, err
	}

	// check project
	if job.Project != authorization.ProjectId {
		return nil, auth.NotAuthorized
	}

	return job, nil
}

func (jobs *Jobs) Get(db *sql.DB) error {
	return jobs.getAllJobs(db, fmt.Sprintf(ownDb.GetAllJobsQuery, ""))
}

func (jobs *Jobs) GetAuthorized(db *sql.DB, authorization *auth.Authorization) error {
	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	return jobs.getAllJobs(db, fmt.Sprintf(ownDb.GetAllJobsQuery, fmt.Sprintf(`WHERE project='%s'`, authorization.ProjectId)))
}

func (job *Job) Get(db *sql.DB) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	err := db.QueryRow(ownDb.GetJobQuery, job.RequestID).Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.Project)
	if err != nil {
		return err
	}
	return nil
}

func (job *Job) GetAuthorized(db *sql.DB, authorization *auth.Authorization) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// get the job
	err = job.Get(db)
	if err != nil {
		return err
	}

	// check project
	if job.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	return nil
}

func (job *Job) Save(db *sql.DB) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	var lastInsertId string
	err := db.QueryRow(ownDb.InsertJobQuery, job.RequestID, job.Version, job.Sender, job.To, job.Timeout, job.Agent, job.Action, job.Payload, job.Status, job.CreatedAt, job.UpdatedAt, job.Project).Scan(&lastInsertId)
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
	if err = tx.QueryRow(ownDb.GetJobQuery, job.RequestID).Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.Project); err != nil {
		return
	}

	log.Infof("%v rows where updated with id %q", affect, job.RequestID)

	return
}

func CleanJobs(db *sql.DB) (affectHeartbeatJobs int64, affectTimeOutJobs int64, err error) {
	if db == nil {
		return 0, 0, errors.New("Clean job routine: Db connection is nil")
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

	affectHeartbeatJobs = 0
	affectTimeOutJobs = 0

	// clean jobs which no heartbeat was send back after created_at + 60 sec
	res, err := tx.Exec(ownDb.CleanJobsNonHeartbeatQuery, 60)
	if err != nil {
		return
	}

	affectHeartbeatJobs, err = res.RowsAffected()
	if err != nil {
		return
	}

	// clean jobs which the timeout + 60 sec has exceeded and still in queued or executing status
	res, err = tx.Exec(ownDb.CleanJobsTimeoutQuery, 60)
	if err != nil {
		return
	}

	affectTimeOutJobs, err = res.RowsAffected()
	if err != nil {
		return
	}

	return
}

// private

func (jobs *Jobs) getAllJobs(db *sql.DB, query string) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	*jobs = make(Jobs, 0)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var job Job
	for rows.Next() {
		err = rows.Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.Project)
		if err != nil {
			log.Errorf("Error scaning job results. Got ", err.Error())
			continue
		}
		*jobs = append(*jobs, job)
	}

	rows.Close()
	return nil
}
