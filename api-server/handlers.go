package main

import (
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/databus23/requestutil"
	"github.com/gorilla/mux"

	"github.com/sapcc/arc/api-server/auth"
	ownDb "github.com/sapcc/arc/api-server/db"
	"github.com/sapcc/arc/api-server/models"
	"github.com/sapcc/arc/api-server/pagination"
	"github.com/sapcc/arc/api-server/pki"
	"github.com/sapcc/arc/arc"
	"github.com/sapcc/arc/version"
)

/*
 * Pki
 */

func servePkiToken(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logAndReturnHttpPkiError(w, http.StatusBadRequest, err)
	}
	var tokenRequest pki.TokenRequest
	if len(body) > 0 {
		if err = json.Unmarshal(body, &tokenRequest); err != nil {
			logAndReturnHttpPkiError(w, http.StatusBadRequest, err)
			return
		}
	}

	// get authentication
	authorization := auth.GetIdentity(r)
	// create token
	token, err := pki.CreateToken(db, authorization, tokenRequest)
	if err != nil {
		logAndReturnHttpPkiError(w, http.StatusInternalServerError, err)
		return
	}

	url := fmt.Sprintf("%s://%s/api/v1/agents/init/%s", requestutil.Scheme(r), requestutil.HostWithPort(r), token)

	info := struct {
		Token       string `json:"token"`
		SignURL     string `json:"url"`
		EndpointURL string `json:"endpoint_url"`
		UpdateURL   string `json:"update_url"`
		ApiURL      string `json:"renew_cert_url"`
	}{
		Token:       token,
		SignURL:     url,
		EndpointURL: agentEndpointURL,
		UpdateURL:   agentUpdateURL,
		ApiURL:      agentApiURL,
	}

	switch r.Header.Get("Accept") {
	case "text/cloud-config":
		w.Header().Set("Content-Type", "text/cloud-config")
		err = cloudConfigInstaller.Execute(w, info)
		checkErrAndReturnStatus(w, err, "Error generating cloud-config reponse", http.StatusInternalServerError, r)
	case "text/x-shellscript":
		w.Header().Set("Content-Type", "text/x-shellscript")
		err = shellScriptInstaller.Execute(w, info)
		checkErrAndReturnStatus(w, err, "Error generating shell script", http.StatusInternalServerError, r)
	case "text/x-powershellscript":
		w.Header().Set("Content-Type", "text/x-powershellscript")
		err = powershellScriptInstaller.Execute(w, info)
		checkErrAndReturnStatus(w, err, "Error generating powershell script", http.StatusInternalServerError, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(info)
		checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
	}

}

func renewPkiCert(w http.ResponseWriter, r *http.Request) {
	// Extract CN, OU and O from the certificate used on the tls communication
	sub, err := pki.TLSRequestSubject(r)
	if err != nil {
		logAndReturnHttpPkiError(w, http.StatusUnauthorized, err)
		return
	}

	// create a TokenRequest with the CN
	var tr pki.TokenRequest
	tr.CN = sub.CommonName

	// add OU and O over the authorization because it will override the tokenRequest
	auth := auth.Authorization{}
	auth.ProjectId = sub.OrganizationalUnit
	auth.ProjectDomainId = sub.Organization

	// read request body
	csr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logAndReturnHttpPkiError(w, http.StatusUnauthorized, err)
		return
	}

	// check we have a valid certificate request
	pemBlock, _ := pem.Decode(csr)
	if pemBlock == nil {
		logAndReturnHttpPkiError(w, http.StatusUnauthorized, errors.New(fmt.Sprint(pki.TLS_CRT_REQUEST_INVALID, err)))
		return
	}
	_, err = x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		logAndReturnHttpPkiError(w, http.StatusUnauthorized, errors.New(fmt.Sprint(pki.TLS_CRT_REQUEST_INVALID, err)))
		return
	}

	// create a token
	token, err := pki.CreateToken(db, &auth, tr)
	if err != nil {
		logAndReturnHttpPkiError(w, http.StatusInternalServerError, errors.New(fmt.Sprint("Error creating token. ", err)))
		return
	}

	// sign tocken with csr
	pemCert, ca, err := pki.SignToken(db, token, csr)
	if err != nil {
		if _, ok := err.(pki.SignForbidden); ok {
			logAndReturnHttpPkiError(w, http.StatusForbidden, err)
		} else {
			logAndReturnHttpPkiError(w, http.StatusInternalServerError, err)
		}
		return
	}

	acceptHeader := r.Header.Get("Accept")

	if acceptHeader == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"certificate": string(*pemCert),
			"ca":          ca,
		}
		err = json.NewEncoder(w).Encode(response)
		checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
	} else {
		w.Header().Set("Content-Type", "application/pkix-cert")
		if _, err := w.Write(*pemCert); err != nil {
			checkErrAndReturnStatus(w, err, "Error writing reponse", http.StatusInternalServerError, r)
		}
	}

}

func signPkiToken(w http.ResponseWriter, r *http.Request) {
	// get token from the request
	vars := mux.Vars(r)
	token := vars["token"]

	csr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logAndReturnHttpPkiError(w, http.StatusBadRequest, err)
	}

	pemCert, ca, err := pki.SignToken(db, token, csr)
	if err != nil {
		if _, ok := err.(pki.SignForbidden); ok {
			logAndReturnHttpPkiError(w, http.StatusForbidden, err)
		} else {
			logAndReturnHttpPkiError(w, http.StatusInternalServerError, err)
		}
		return
	}

	acceptHeader := r.Header.Get("Accept")

	if acceptHeader == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"certificate": string(*pemCert),
			"ca":          ca,
		}
		err = json.NewEncoder(w).Encode(response)
		checkErrAndReturnStatus(w, err, "Error encoding response to JSON", http.StatusInternalServerError, r)
	} else {
		w.Header().Set("Content-Type", "application/pkix-cert")
		if _, err = w.Write(*pemCert); err != nil {
			checkErrAndReturnStatus(w, err, "Error writing response: ", http.StatusInternalServerError, r)
		}
	}
}

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
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "Error getting all jobs. ", http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, "Error getting all jobs. ", http.StatusInternalServerError, r)
		}
		return
	}

	// set pagination header
	pagination.SetHeaders(w)

	// set the header and body
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(jobs)
	checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
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
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Job with id %q.", jobId), http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Job with id %q.", jobId), http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Job with id %q.", jobId), http.StatusInternalServerError, r)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(job)
	checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
}

func executeJob(w http.ResponseWriter, r *http.Request) {
	errorText := "Error creating a job. "

	// read request body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		checkErrAndReturnStatus(w, err, errorText, http.StatusBadRequest, r)
		return
	}
	defer r.Body.Close()

	// get authentication parameters
	authorization := auth.GetIdentity(r)

	// create job
	job, err := models.CreateJobAuthorized(db, &data, config.Identity, authorization)
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
		if err != nil {
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
	checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
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
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Logs for Job with id  %q.", jobId), http.StatusInternalServerError, r)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(logEntry.Content)); err != nil {
		checkErrAndReturnStatus(w, err, "Error writing response: ", http.StatusInternalServerError, r)
	}
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
			logInfoAndReturnHttpErrStatus(w, err, "Invalid selector expression.", http.StatusBadRequest, r)
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, "Error getting all agents.", http.StatusInternalServerError, r)
		}
		return
	}

	// set pagination header
	pagination.SetHeaders(w)

	// set the header and body
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agents)
	checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
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
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Error getting agent with id %q. ", agentId), http.StatusInternalServerError, r)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agent)
	checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
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
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found. Got: ", agentId), http.StatusNotFound, r)
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Error deleten agent with id %q. ", agentId), http.StatusInternalServerError, r)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q", agentId), http.StatusInternalServerError, r)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agent.Facts)
	checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
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
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q", agentId), http.StatusInternalServerError, r)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(agent.Tags)
	checkErrAndReturnStatus(w, err, "Error generating JSON reponse", http.StatusInternalServerError, r)
}

func saveAgentTags(w http.ResponseWriter, r *http.Request) {
	// check authentication
	authorization := auth.GetIdentity(r)

	// get the agent id
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		checkErrAndReturnStatus(w, err, fmt.Sprintf("Error saving agent tags."), http.StatusBadRequest, r)
	}
	var dataStruct map[string]string
	if len(body) > 0 {
		if err = json.Unmarshal(body, &dataStruct); err != nil {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Error saving agent tags."), http.StatusBadRequest, r)
			return
		}
	}

	// process data
	err = models.ProcessTags(db, authorization, agentId, dataStruct)
	if err != nil {
		if err == sql.ErrNoRows {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q not found", agentId), http.StatusNotFound, r)
		} else if serr, ok := err.(*models.TagError); ok {
			jsonString, err := serr.MessagesToJson()
			if err != nil {
				checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError, r)
			}
			http.Error(w, jsonString, http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Agent with id %q. ", agentId), http.StatusInternalServerError, r)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
		} else if _, ok := err.(auth.IdentityStatusInvalid); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else if _, ok := err.(auth.NotAuthorized); ok {
			logInfoAndReturnHttpErrStatus(w, err, "", http.StatusUnauthorized, r)
		} else {
			checkErrAndReturnStatus(w, err, fmt.Sprintf("Error removing tag %q from Agent id %q. ", value, agentId), http.StatusInternalServerError, r)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

/*
 * Root and Healthcheck
 */

func serveVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("Arc api-server " + version.String())); err != nil {
		checkErrAndReturnStatus(w, err, "Error writing response: ", http.StatusInternalServerError, r)
	}
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
	if _, err = w.Write([]byte("Ready!!!")); err != nil {
		checkErrAndReturnStatus(w, err, "Error writing response: ", http.StatusInternalServerError, r)
	}
}

// private

func logAndReturnHttpPkiError(w http.ResponseWriter, status int, err error) {
	_, file, line, _ := runtime.Caller(1)
	log.Errorf("PKI request error. status=%d location=%s:%d error=%v", status, path.Base(file), line, err)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}
	w.WriteHeader(status)
	if err != nil {
		// json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		if errJSONEncode := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJSONEncode != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, err.Error(), status)
		} else {
			w.WriteHeader(status)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		}
	}
}

func logInfoAndReturnHttpErrStatus(w http.ResponseWriter, err error, msg string, status int, r *http.Request) {
	// status string, code int, title string, detail string
	apiError := NewApiError(http.StatusText(status), status, msg, err, r)

	if err != nil {
		log.Infof("Error, returning status %v.\nDetails: %+v", apiError.Status, apiError)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, apiError.toString(), status)
}

func checkErrAndReturnStatus(w http.ResponseWriter, err error, msg string, status int, r *http.Request) {
	if err != nil {
		apiError := NewApiError(http.StatusText(status), status, msg, err, r)
		log.Errorf("Error, returning status %v.\nDetails: %+v", apiError.Status, apiError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, apiError.toString(), status)
	}
}
