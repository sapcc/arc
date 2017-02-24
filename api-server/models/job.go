package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"

	"gitHub.***REMOVED***/monsoon/arc/api-server/auth"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pagination"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

var (
	metricJobExecuted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "arc_job_executed",
		Help: "Total number of jobs executed.",
	})
)

type JobTargetAgentNotFoundError struct {
	Msg string
}

func (e JobTargetAgentNotFoundError) Error() string {
	return e.Msg
}

type JobBadRequestError struct {
	Msg string
}

func (e JobBadRequestError) Error() string {
	return e.Msg
}

type Job struct {
	arc.Request
	Status    arc.JobState `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Project   string       `json:"project"`
	User      JSONB        `json:"user"`
}

type JobID struct {
	RequestID string `json:"request_id"`
}

type Jobs []Job

type Status string

func init() {
	prometheus.MustRegister(metricJobExecuted)
}

func CreateJob(db *sql.DB, data *[]byte, identity string, user *auth.User) (*Job, error) {
	if db == nil {
		return nil, errors.New("Db is nil")
	}

	if user.Id == "" {
		return nil, errors.New("User id is blank")
	}

	// unmarshal data
	var tmpReq arc.Request
	err := json.Unmarshal(*data, &tmpReq)
	if err != nil {
		return nil, JobBadRequestError{Msg: err.Error()}
	}

	// create a validate request
	request, err := arc.CreateRequest(tmpReq.Agent, tmpReq.Action, identity, tmpReq.To, tmpReq.Timeout, tmpReq.Payload)
	if err != nil {
		return nil, JobBadRequestError{Msg: err.Error()}
	}

	// check and get the target project id
	agent := Agent{AgentID: tmpReq.To}
	err = agent.Get(db)
	if err == sql.ErrNoRows {
		return nil, JobTargetAgentNotFoundError{Msg: err.Error()}
	} else if err != nil {
		return nil, err
	}

	job := Job{
		*request,
		arc.Queued,
		time.Now(),
		time.Now(),
		agent.Project,
		JSONB{},
	}

	// add the user
	userJsonb, err := JobUserToJSONB(*user)
	if err != nil {
		return nil, err
	}
	job.User = *userJsonb

	// increment metric
	metricJobExecuted.Inc()

	return &job, nil
}

func CreateJobAuthorized(db *sql.DB, data *[]byte, identity string, authorization *auth.Authorization) (*Job, error) {
	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return nil, err
	}

	job, err := CreateJob(db, data, identity, &authorization.User)
	if err != nil {
		return nil, err
	}

	// check project
	if job.Project != authorization.ProjectId {
		return nil, auth.NotAuthorized{Msg: fmt.Sprintf("%s is not project %s", job.Project, authorization.ProjectId)}
	}

	return job, nil
}

func (jobs *Jobs) Get(db *sql.DB) error {
	return jobs.getAllJobs(db, buildJobsQuery(ownDb.GetAllJobsQuery, "", "", nil))
}

func (jobs *Jobs) GetAuthorized(db *sql.DB, authorization *auth.Authorization, agentId string, pag *pagination.Pagination) error {
	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// count jobs and set total pages
	countJobs, err := countJobs(db, buildJobsQuery(ownDb.CountAllJobsQuery, authorization.ProjectId, agentId, nil))
	if err != nil {
		return err
	}
	err = pag.SetTotalElements(countJobs)
	if err != nil {
		return err
	}

	return jobs.getAllJobs(db, buildJobsQuery(ownDb.GetAllJobsQuery, authorization.ProjectId, agentId, pag))
}

func (job *Job) Get(db *sql.DB) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	err := db.QueryRow(ownDb.GetJobQuery, job.RequestID).Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.Project, &job.User)
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
		return auth.NotAuthorized{Msg: fmt.Sprintf("%s is not project %s", job.Project, authorization.ProjectId)}
	}

	return nil
}

func (job *Job) Save(db *sql.DB) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	jobUser, err := job.User.Value()
	if err != nil {
		return err
	}

	var lastInsertId string
	err = db.QueryRow(ownDb.InsertJobQuery, job.RequestID, job.Version, job.Sender, job.To, job.Timeout, job.Agent, job.Action, job.Payload, job.Status, job.CreatedAt, job.UpdatedAt, job.Project, jobUser).Scan(&lastInsertId)
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
		return err
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
	if err = tx.QueryRow(ownDb.GetJobQuery, job.RequestID).Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.Project, &job.User); err != nil {
		return
	}

	log.Infof("%v rows where updated with id %q", affect, job.RequestID)

	return
}

func CleanJobs(db *sql.DB) (affectHeartbeatJobs int64, affectTimeOutJobs int64, affectOldJobs int64, err error) {
	if db == nil {
		return 0, 0, 0, errors.New("Clean job routine: Db connection is nil")
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, 0, 0, err
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
	affectOldJobs = 0

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

	// clean jobs which are older than 30 days
	res, err = tx.Exec(ownDb.CleanJobsOldQuery, 30)
	if err != nil {
		return
	}

	affectOldJobs, err = res.RowsAffected()
	if err != nil {
		return
	}

	return
}

// private

func buildJobsQuery(baseQuery string, authProjectId, agentId string, pag *pagination.Pagination) string {
	resultQuery := fmt.Sprintf(baseQuery, "", "")
	authQuery := ""
	paginationQuery := ""

	// check pagination
	if pag != nil {
		paginationQuery = fmt.Sprintf(`OFFSET %v LIMIT %v`, pag.Offset, pag.Limit)
		resultQuery = fmt.Sprintf(baseQuery, "", paginationQuery)
	}

	// check authority
	if authProjectId != "" {
		authQuery = fmt.Sprintf(`project = '%s'`, authProjectId)
	}

	if authQuery != "" {
		resultQuery = fmt.Sprintf(baseQuery, fmt.Sprint("WHERE ", authQuery), paginationQuery)
		if agentId != "" {
			resultQuery = fmt.Sprintf(baseQuery, fmt.Sprintf(`WHERE %s AND ( "to" = '%s')`, authQuery, agentId), paginationQuery)
		}
	} else {
		if agentId != "" {
			resultQuery = fmt.Sprintf(baseQuery, fmt.Sprintf(`WHERE "to" = '%s'`, agentId), paginationQuery)
		}
	}

	return resultQuery
}

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
		err = rows.Scan(&job.RequestID, &job.Version, &job.Sender, &job.To, &job.Timeout, &job.Agent, &job.Action, &job.Payload, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.Project, &job.User)
		if err != nil {
			log.Errorf("Error scaning job results. Got ", err.Error())
			continue
		}
		*jobs = append(*jobs, job)
	}

	rows.Close()
	return nil
}

func countJobs(db *sql.DB, query string) (int, error) {
	if db == nil {
		return 0, errors.New("Db is nil")
	}

	var countJob int
	err := db.QueryRow(query).Scan(&countJob)
	if err != nil {
		return 0, err
	}

	return countJob, nil
}
