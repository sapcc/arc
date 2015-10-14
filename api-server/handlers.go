package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

/*
 * Jobs
 */

func serveJobs(w http.ResponseWriter, r *http.Request) {
	// get authentication
	authorization := auth.GetIdentity(r)

	// read jobs
	jobs := models.Jobs{}
	err := jobs.GetAuthorized(db, authorization)
	if err == auth.IdentityStatusInvalid || err == auth.NotAuthorized {
		logInfoAndReturnHttpErrStatus(w, err, "Error getting all jobs. ", http.StatusUnauthorized)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, "Error getting all jobs. ", http.StatusInternalServerError)
		return
	}

	// set the header and body
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(jobs)
	checkErrAndReturnStatus(w, err, "Error encoding Jobs to JSON", http.StatusInternalServerError)
}

func serveJob(w http.ResponseWriter, r *http.Request) {
	// get authentication
	authorization := auth.GetIdentity(r)

	vars := mux.Vars(r)
	jobId := vars["jobId"]

	job := models.Job{Request: arc.Request{RequestID: jobId}}
	err := job.GetAuthorized(db, authorization)
	if err == sql.ErrNoRows {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Job with id %q not found", jobId), http.StatusNotFound)
		return
	} else if err == auth.IdentityStatusInvalid || err == auth.NotAuthorized {
		logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Job with id %q.", jobId), http.StatusUnauthorized)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Job with id %q.", jobId), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(job)
	checkErrAndReturnStatus(w, err, "Error encoding Jobs to JSON", http.StatusInternalServerError)
}

func executeJob(w http.ResponseWriter, r *http.Request) {
	// get authentication
	authorization := auth.GetIdentity(r)

	// read request body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Error creating a job. ", http.StatusBadRequest)
		return
	}

	// create job
	job, err := models.CreateJobAuthorized(db, &data, config.Identity, authorization)
	if err == models.JobTargetAgentNotFoundError {
		checkErrAndReturnStatus(w, err, "Error creating a job. ", http.StatusNotFound)
		return
	} else if err == auth.IdentityStatusInvalid || err == auth.NotAuthorized {
		logInfoAndReturnHttpErrStatus(w, err, "Error creating a job. ", http.StatusUnauthorized)
		return
	} else if err == models.JobBadRequestError {
		checkErrAndReturnStatus(w, err, "Error creating a job. ", http.StatusBadRequest)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, "Error creating a job. ", http.StatusInternalServerError)
		return
	}

	// save db
	err = job.Save(db)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Error creating a job. ", http.StatusInternalServerError)
		return
	}

	// create a mqtt request
	err = arcSendRequest(&job.Request)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Error creating a job. ", http.StatusInternalServerError)
		return
	}

	// create response
	response := models.JobID{job.RequestID}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(response)
	checkErrAndReturnStatus(w, err, "Error encoding Jobs to JSON", http.StatusInternalServerError)
}

/*
 * Logs
 */

func serveJobLog(w http.ResponseWriter, r *http.Request) {
	// get authentication
	authorization := auth.GetIdentity(r)
	
	// get the job id
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	// get the log
	logEntry := models.Log{JobID: jobId}
	err := logEntry.GetOrCollectAuthorized(db, authorization)
	if err == sql.ErrNoRows {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Logs for Job with id %q not found.", jobId), http.StatusNotFound)
		return
	} else if err == auth.IdentityStatusInvalid || err == auth.NotAuthorized {
		logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusUnauthorized)
		return		
	} else if err != nil {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(logEntry.Content))
}

/*
 * Agents with authorization
 */

func serveAgents(w http.ResponseWriter, r *http.Request) {
	// get authentication
	authorization := auth.GetIdentity(r)

	// get agents
	agents := models.Agents{}
	err := agents.GetAuthorized(db, r.URL.Query().Get("q"), authorization)

	if err == models.FilterError {
		checkErrAndReturnStatus(w, err, "Error serving filtered Agents.", http.StatusBadRequest)
		return
	} else if err == auth.IdentityStatusInvalid || err == auth.NotAuthorized {
		logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, "Error getting all agents.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agents)
}

func serveAgent(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get agent
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	agent := models.Agent{AgentID: agentId}
	err := agent.GetAuthorized(db, authorization)
	if err == sql.ErrNoRows {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found. Got %q", agentId), http.StatusNotFound)
		return
	} else if err == auth.IdentityStatusInvalid || err == auth.NotAuthorized {
		logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q", agentId), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agent)
	checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError)
}

/*
 * Facts
 */

func serveFacts(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get the agent id
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	// get the agent
	agent := models.Agent{AgentID: agentId}
	err := agent.GetAuthorized(db, authorization)
	if err == sql.ErrNoRows {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found", agentId), http.StatusNotFound)
		return
	} else if err == auth.IdentityStatusInvalid || err == auth.NotAuthorized {
		logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q", agentId), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(agent.Facts))
}

/*
 * Root and Healthcheck
 */

func serveVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Arc api-server " + version.String()))
}

/*
 * ServeReadiness
 */

type Readiness struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

func serveReadiness(w http.ResponseWriter, r *http.Request) {
	//check db connection
	rows, err := db.Query(ownDb.CheckConnection)
	if err != nil {
		ready := Readiness{
			Status:  http.StatusBadGateway,
			Message: "Ping to the DB failed",
		}

		// convert struct to json
		body, err := json.Marshal(ready)
		checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError)

		// return the error with json body
		http.Error(w, string(body), ready.Status)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		log.Errorf("Error, returning status %v. %s", ready.Status, ready.Message)
		return
	}
	defer rows.Close()

	//check mosquitto transport connection
	if !tp.IsConnected() {
		ready := Readiness{
			Status:  http.StatusBadGateway,
			Message: "Ping to the transport failed",
		}

		// convert struct to json
		body, err := json.Marshal(ready)
		checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError)

		// return the error with json body
		http.Error(w, string(body), ready.Status)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		log.Errorf("Error, returning status %v. %s", ready.Status, ready.Message)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Ready!!!"))
}

// private

func logInfoAndReturnHttpErrStatus(w http.ResponseWriter, err error, msg string, status int) {
	if err != nil {
		log.Infof("Error, returning status %v. %s %s", status, msg, err.Error())
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, http.StatusText(status), status)
}

func checkErrAndReturnStatus(w http.ResponseWriter, err error, msg string, status int) {
	if err != nil {
		log.Errorf("Error, returning status %v. %s %s", status, msg, err.Error())
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, http.StatusText(status), status)
	}
}
