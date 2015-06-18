package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

/*
 * Jobs
 */

func serveJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := models.GetAllJobs(db)
	if err != nil {
		log.Errorf("Error getting all jobs. Got %q", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(jobs)
}

func serveJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	job, err := models.GetJob(db, jobId)
	if err != nil {
		log.Errorf("Job with id %q not found. Got %q", jobId, err.Error())
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(job)
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
	err = models.SaveJob(db, job)
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

	resLogs, err := models.GetLog(db, jobId)
	if err != nil {
		log.Errorf("Logs for Job with id %q not found.", jobId)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(*resLogs))
}

/*
 * Agents
 */

func serveAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := models.GetAgents(db)
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

	agent, err := models.GetAgent(db, agentId)
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

	agent, err := models.GetFact(db, agentId)
	if err != nil {
		log.Errorf("Agent with id %q not found. Got %q", agentId, err.Error())
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//json.NewEncoder(w).Encode(agent.Facts)
	w.Write([]byte(agent.Facts))
}
