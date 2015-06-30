package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

/*
 * Jobs
 */

func serveJobs(w http.ResponseWriter, r *http.Request) {
	jobs := models.Jobs{}
	if err := jobs.Get(db); err != nil {
		checkErrAndReturnStatus(w, err, "Error getting all jobs", http.StatusInternalServerError)
		return
	}

	// set the header and body
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(jobs)
	checkErrAndReturnStatus(w, err, "Error encoding Jobs to JSON", http.StatusInternalServerError)
}

func serveJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	job := models.Job{Request: arc.Request{RequestID: jobId}}
	if err := job.Get(db); err != nil {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Job with id %q not found", jobId), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(job)
	checkErrAndReturnStatus(w, err, "Error encoding Jobs to JSON", 500)
}

func executeJob(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error creating a job. Got %q", err.Error())
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// create job
	job, err := models.CreateJob(&data, config.Identity)
	if err != nil {
		log.Errorf("Error creating a job. Got %q", err.Error())
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// save db
	err = job.Save(db)
	if err != nil {
		log.Errorf("Error saving job. Got %q", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// create a mqtt request
	arcSendRequest(&job.Request)

	// create response
	response := models.JobID{job.RequestID}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

/*
 * Logs
 */

func serveJobLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	logEntry := models.Log{JobID: jobId}
	err := logEntry.Get(db)
	if err != nil {
		log.Errorf("Logs for Job with id %q not found.", jobId)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(logEntry.Content))
}

/*
 * Agents
 */

func serveAgents(w http.ResponseWriter, r *http.Request) {
	agents := models.Agents{}
	err := agents.Get(db)
	if err != nil {
		log.Errorf("Error getting all agents. Got %q", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agents)
}

func serveAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	agent := models.Agent{AgentID: agentId}
	err := agent.Get(db)
	if err != nil {
		log.Errorf("Agent with id %q not found. Got %q", agentId, err.Error())
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agent)
}

/*
 * Facts
 */

func serveFacts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	fact := models.Fact{Agent: models.Agent{AgentID: agentId}}
	err := fact.Get(db)
	if err != nil {
		log.Errorf("Agent with id %q not found. Got %q", agentId, err.Error())
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(fact.Facts))
}

/*
 * Root
 */

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Arc api-server " + version.String()))
}

// private

func checkErrAndReturnStatus(w http.ResponseWriter, err error, msg string, status int) {
	if err != nil {
		log.Errorf("Error, returning status %v. %s %s", status, msg, err.Error())
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, http.StatusText(status), status)
	}
}
