package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pagination"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

/*
 * Jobs
 */

func serveJobs(w http.ResponseWriter, r *http.Request) {
	// get authentication
	authorization := auth.GetIdentity(r)

	// filter per agent
	filterAgentId := r.URL.Query().Get("agent_id")

	// pagination
	pagination := pagination.CreatePagination(*r.URL)

	// read jobs
	jobs := models.Jobs{}
	err := jobs.GetAuthorized(db, authorization, filterAgentId, pagination)
	if err != nil {
		if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "Error getting all jobs. ", http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "Error getting all jobs. ", http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, "Error getting all jobs. ", http.StatusInternalServerError, r)
			return
		}
	}

	// set pagination header
	w.Header().Set("Pagination-Elements", fmt.Sprintf("%v", pagination.TotalElements))
	w.Header().Set("Pagination-Pages", fmt.Sprintf("%v", pagination.TotalPages))
	w.Header().Set("Link", pagination.GetLinks())

	// set the header and body
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(jobs)
	checkErrAndReturnStatus(w, err, "Error encoding Jobs to JSON", http.StatusInternalServerError, r)
}

func serveJob(w http.ResponseWriter, r *http.Request) {
	// get authentication
	authorization := auth.GetIdentity(r)

	vars := mux.Vars(r)
	jobId := vars["jobId"]

	job := models.Job{Request: arc.Request{RequestID: jobId}}
	err := job.GetAuthorized(db, authorization)
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Job with id %q not found", jobId), http.StatusNotFound, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Job with id %q.", jobId), http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Job with id %q.", jobId), http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Job with id %q.", jobId), http.StatusInternalServerError, r)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(job)
	checkErrAndReturnStatus(w, err, "Error encoding Jobs to JSON", http.StatusInternalServerError, r)
}

func executeJob(w http.ResponseWriter, r *http.Request) {
	errorText := "Error creating a job. "

	// get authentication parameters
	authorization := auth.GetIdentity(r)

	// read request body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		checkErrAndReturnStatus(w, err, errorText, http.StatusBadRequest, r)
		return
	}
	defer r.Body.Close()

	// create job
	job, err := models.CreateJobAuthorized(db, &data, config.Identity, authorization)
	if err != nil {
		if _, ok := err.(models.JobBadRequestError); ok {
			checkErrAndReturnStatus(w, err, errorText, http.StatusBadRequest, r)
			return
		} else if _, ok := err.(models.JobTargetAgentNotFoundError); ok {
			checkErrAndReturnStatus(w, err, errorText, http.StatusNotFound, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, errorText, http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, errorText, http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, errorText, http.StatusInternalServerError, r)
			return
		}
	}

	// save db
	err = job.Save(db)
	if err != nil {
		checkErrAndReturnStatus(w, err, errorText, http.StatusInternalServerError, r)
		return
	}

	// create a mqtt request
	err = arcSendRequest(&job.Request)
	if err != nil {
		checkErrAndReturnStatus(w, err, errorText, http.StatusInternalServerError, r)
		return
	}

	// create response
	response := models.JobID{job.RequestID}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(response)
	checkErrAndReturnStatus(w, err, "Error encoding Jobs to JSON", http.StatusInternalServerError, r)
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
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Logs for Job with id %q not found.", jobId), http.StatusNotFound, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusInternalServerError, r)
			return
		}
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

	// filter facts
	filterFacts := r.URL.Query().Get("q")

	// show facts
	showFacts := strings.Split(r.URL.Query().Get("facts"), ",")

	// pagination
	pagination := pagination.CreatePagination(*r.URL)

	// get agents
	agents := models.Agents{}
	err := agents.GetAuthorizedAndShowFacts(db, filterFacts, authorization, showFacts, pagination)
	if err != nil {
		if _, ok := err.(models.FilterError); ok {
			logInfoAndReturnHttpErrStatus(w, err, "Error serving filtered Agents.", http.StatusBadRequest, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, "Error getting all agents.", http.StatusInternalServerError, r)
			return
		}
	}

	// set pagination header
	w.Header().Set("Pagination-Elements", fmt.Sprintf("%v", pagination.TotalElements))
	w.Header().Set("Pagination-Pages", fmt.Sprintf("%v", pagination.TotalPages))
	w.Header().Set("Link", pagination.GetLinks())

	// set the header and body
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agents)
	checkErrAndReturnStatus(w, err, "Error encoding Agents to JSON", http.StatusInternalServerError, r)
}

func serveAgent(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get agent
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	// Show facts
	showFacts := strings.Split(r.URL.Query().Get("facts"), ",")

	agent := models.Agent{AgentID: agentId}
	err := agent.GetAuthorizedAndShowFacts(db, authorization, showFacts)
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found. ", agentId), http.StatusNotFound, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Error getting agent with id %q. ", agentId), http.StatusInternalServerError, r)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agent)
	checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError, r)
}

func deleteAgent(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get agent
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	agent := models.Agent{AgentID: agentId}
	err := agent.DeleteAuthorized(db, authorization)
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found. Got %q", agentId), http.StatusNotFound, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Error deleten agent with id %q. ", agentId), http.StatusInternalServerError, r)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Agent with id %q deleted. ", agentId)))
}

/*
 * Facts with authorization
 */

func serveFacts(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get the agent id
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	// get the agent
	agent := models.Agent{AgentID: agentId}
	err := agent.GetAuthorizedAndShowFacts(db, authorization, []string{"all"})
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found", agentId), http.StatusNotFound, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q", agentId), http.StatusInternalServerError, r)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agent.Facts)
	checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError, r)
}

/*
 * Tags
 */

func serveAgentTags(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get the agent id
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	// get the agent
	agent := models.Agent{AgentID: agentId}
	err := agent.GetAuthorizedAndShowFacts(db, authorization, []string{"all"})
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found", agentId), http.StatusNotFound, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q", agentId), http.StatusInternalServerError, r)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agent.Tags)
	checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError, r)
}

func saveAgentTags(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get the agent id
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	// parse form
	r.ParseForm()

	// process data
	err := models.ProcessTags(db, authorization, agentId, r.Form)
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found", agentId), http.StatusNotFound, r)
			return
		} else if serr, ok := err.(*models.TagError); ok {
			jsonString, err := serr.MessagesToJson()
			if err != nil {
				checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError, r)
			}
			http.Error(w, jsonString, http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q. ", agentId), http.StatusInternalServerError, r)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("All tags saved!"))
}

func deleteAgentTag(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get the agent id
	vars := mux.Vars(r)
	agentId := vars["agentId"]
	value := vars["value"]

	// get the agent
	agent := models.Agent{AgentID: agentId}
	err := agent.DeleteTagAuthorized(db, authorization, value)
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found", agentId), http.StatusNotFound, r)
			return
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
			return
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Error removing tag %q from Agent id %q. ", value, agentId), http.StatusInternalServerError, r)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Tag from agent with id %q and value %q is removed!", agentId, value)))
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
		checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError, r)

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
		checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError, r)

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

func logInfoAndReturnHttpErrStatus(w http.ResponseWriter, err error, msg string, status int, r *http.Request) {
	// status string, code int, title string, detail string
	apiError := NewApiError(http.StatusText(status), status, msg, err, r)

	if err != nil {
		log.Infof("Error, returning status %v.\nDetails: %+v", apiError.Status, apiError)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, apiError.String(), status)
}

func checkErrAndReturnStatus(w http.ResponseWriter, err error, msg string, status int, r *http.Request) {
	if err != nil {
		apiError := NewApiError(http.StatusText(status), status, msg, err, r)
		log.Errorf("Error, returning status %v.\nDetails: %+v", apiError.Status, apiError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, apiError.String(), status)
	}
}
