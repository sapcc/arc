package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"net/http"
)

func serveJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(jobs)
}

func serveJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["jobId"]

	job := getJob(agentId)
	if job == nil {
		log.Errorf("Job with id %q not found.", agentId)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(job)
}

func executeJob(w http.ResponseWriter, r *http.Request) {
	// create job
	job, err := models.CreateJob(&r.Body)
	if err != nil {
		log.Errorf("Error creating a job. Got %q", err.Error())
		http.Error(w, http.StatusText(400), 400)
	} else {
		job.Status = models.Queued
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(job)
	}

	// create a mqtt request

	// save db
}

func serveAgents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agents)
}

func serveAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	agent := getAgent(agentId)
	if agent == nil {
		log.Errorf("Agent with id %q not found.", agentId)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agent)
}

func serveFacts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	agent := getAgent(agentId)
	if agent == nil {
		log.Errorf("Agent with id %q not found.", agentId)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agent.Facts)
}

func serveFact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["agentId"]
	factId := vars["factId"]

	agent := getAgent(agentId)
	if agent == nil {
		log.Errorf("Agent with id %q not found.", agentId)
		http.NotFound(w, r)
		return
	}

	var fact models.Fact
	for _, f := range agent.Facts {
		if f.Id == factId {
			fact = models.Fact(f)
		}
	}
	if len(fact.Id) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(fact)
}

// private

func getJob(jobId string) *models.Job {
	var job models.Job
	for _, j := range jobs {
		if j.Request.RequestID == jobId {
			job = models.Job(j)
		}
	}
	if len(job.Request.RequestID) == 0 {
		return nil
	}
	return &job
}

func getAgent(agentId string) *models.Agent {
	var agent models.Agent
	for _, a := range agents {
		if a.Id == agentId {
			agent = models.Agent(a)
		}
	}
	if len(agent.Id) == 0 {
		return nil
	}
	return &agent
}
