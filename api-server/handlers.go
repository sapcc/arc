package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
  "github.com/gorilla/handlers"	

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
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
	err := job.Get(db)
	if err == sql.ErrNoRows {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Job with id %q not found", jobId), http.StatusNotFound)
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
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Error creating a job", http.StatusBadRequest)
		return
	}

	// create job
	job, err := models.CreateJob(&data, config.Identity)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Error creating a job", http.StatusBadRequest)
		return
	}

	// save db
	err = job.Save(db)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Error saving job", http.StatusInternalServerError)
		return
	}

	// create a mqtt request
	err = arcSendRequest(&job.Request)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Error saving job", http.StatusInternalServerError)
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
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	logEntry := models.Log{JobID: jobId}
	err := logEntry.GetOrCollect(db)
	if err == sql.ErrNoRows {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Logs for Job with id %q not found.", jobId), http.StatusNotFound)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(logEntry.Content))
}

/*
 * Agents
 */

func serveAgents(w http.ResponseWriter, r *http.Request) {
	agents := models.Agents{}
	err := agents.Get(db, r.URL.Query().Get("q"))

	if err == models.FilterError {
		checkErrAndReturnStatus(w, err, "Error serving filtered Agents.", http.StatusBadRequest)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, "Error getting all agents.", http.StatusInternalServerError)
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
	if err == sql.ErrNoRows {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found. Got %q", agentId), http.StatusNotFound)
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
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	agent := models.Agent{AgentID: agentId}
	err := agent.Get(db)
	if err == sql.ErrNoRows {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found", agentId), http.StatusNotFound)
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
	Status   int  `json:"status"`
	Message	 string  `json:"error"`
}

func serveReadiness(w http.ResponseWriter, r *http.Request) {	
	//check db connection
	rows, err := db.Query(ownDb.CheckConnection)
	if err != nil {
		ready := Readiness{
			Status: http.StatusBadGateway,
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
			Status: http.StatusBadGateway,
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

func combineLogHandler(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func loggingHandler(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    t1 := time.Now()
    next.ServeHTTP(w, r)
    t2 := time.Now()
    log.Infof("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
  }

  return http.HandlerFunc(fn)
}

// private

func checkErrAndReturnStatus(w http.ResponseWriter, err error, msg string, status int) {
	if err != nil {
		log.Errorf("Error, returning status %v. %s %s", status, msg, err.Error())
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, http.StatusText(status), status)
	}
}
