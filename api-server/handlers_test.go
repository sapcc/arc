// +build integration

package main

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Job Handlers", func() {

	var (
		job models.Job
	)

	Describe("serveJobs", func() {

		It("returns a 500 error if something goes wrong", func() {
			// save bad data
			job := models.Job{}
			job.RpcVersionExample()
			job.Status = 6 // not existing status
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// make a request
			req, err := http.NewRequest("GET", "/jobs", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code, header and body
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))
		})

		It("returns empty json arry if no jobs found", func() {
			// make a request
			req, err := http.NewRequest("GET", "/jobs", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			jobs := make(models.Jobs, 0)
			err = json.Unmarshal(w.Body.Bytes(), &jobs)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobs).To(Equal(make(models.Jobs, 0)))
		})

		It("returns all jobs", func() {
			// fill db
			jobs := models.Jobs{}
			jobs.CreateAndSaveRpcVersionExamples(db, 3)

			// make request
			req, err := http.NewRequest("GET", "/jobs", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			dbJobs := make(models.Jobs, 0)
			err = json.Unmarshal(w.Body.Bytes(), &dbJobs)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbJobs[0].RequestID).To(Equal(jobs[0].RequestID))
			Expect(dbJobs[1].RequestID).To(Equal(jobs[1].RequestID))
			Expect(dbJobs[2].RequestID).To(Equal(jobs[2].RequestID))
		})

	})

	Describe("serveJob", func() {

		JustBeforeEach(func() {
			// save a job
			job = models.Job{}
			job.RpcVersionExample()
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a 404 error if job not found", func() {
			// make request
			req, err := http.NewRequest("GET", "/jobs/non_existing_id", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 500 error if something is wrong with the query", func() {
			// wrong query
			tmp_GetJobQuery := GetJobQuery
			GetJobQuery = "SELECT * Wrong_Job_Table jobs WHERE id=$1"

			// make request
			req, err := http.NewRequest("GET", fmt.Sprint("/jobs/", job.RequestID), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			// copy back the jobquery defenition
			GetJobQuery = tmp_GetJobQuery
		})

		It("returns a 500 error if something else goes wrong", func() {
			// introduce bad data
			job.Status = 6 // not existing status should produce an error when encoding job to json
			err := job.Update(db)
			Expect(err).NotTo(HaveOccurred())

			// make a request
			req, err := http.NewRequest("GET", fmt.Sprint("/jobs/", job.RequestID), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))
		})

		It("should return the job", func() {
			// make a request
			req, err := http.NewRequest("GET", fmt.Sprint("/jobs/", job.RequestID), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			var dbJob models.Job
			err = json.Unmarshal(w.Body.Bytes(), &dbJob)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJob.RequestID).To(Equal(job.RequestID))
		})
	})

	Describe("executeJob", func() {

		JustBeforeEach(func() {
			config.Identity = "darwin"
			config.Transport = "fake"
			var err error
			tp, err = arcNewConnection(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a 400 error if request is wrong", func() {
			jsonStr := []byte(`this is not json`)
			// make a request
			req, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jsonStr))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(400))
		})

		It("returns a 500 error if something goes wrong", func() {
			// copy the db pointer
			db_copy := db
			db = nil

			jsonStr := []byte(`{"to":"darwin","timeout":60,"agent":"rpc","action":"version"}`)
			// make a request
			req, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jsonStr))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			// set the pointer back to the db
			db = db_copy
		})

		It("should save the job and return the unique id as JSON", func() {
			jsonStr := []byte(`{"to":"darwin","timeout":60,"agent":"rpc","action":"version"}`)
			// make a request
			req, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jsonStr))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			var dbJobID models.JobID
			err = json.Unmarshal(w.Body.Bytes(), &dbJobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJobID.RequestID).NotTo(BeEmpty())

			// check the job is being saved
			dbJob := models.Job{Request: arc.Request{RequestID: dbJobID.RequestID}}
			dbJob.Get(db)
			Expect(dbJob.RequestID).To(Equal(dbJobID.RequestID))
		})

	})

})

var _ = Describe("Facts Handlers", func() {

	Describe("serveAgents", func() {

		var (
			agents models.Agents
		)

		It("returns a 500 error if something goes wrong", func() {
			tmp_GetAgentsQuery := GetAgentsQuery
			GetAgentsQuery = "SELECT DISTINCT agent_id,created_at,updated_at FROM Wrong_Facts order by updated_at"

			// make a request
			req, err := http.NewRequest("GET", "/agents", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			GetAgentsQuery = tmp_GetAgentsQuery
		})

		It("returns empty json arry if no agents found", func() {
			// make a request
			req, err := http.NewRequest("GET", "/agents", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			agents = make(models.Agents, 0)
			err = json.Unmarshal(w.Body.Bytes(), &agents)
			Expect(err).NotTo(HaveOccurred())
			Expect(agents).To(Equal(make(models.Agents, 0)))
		})

		It("returns all agents", func() {
			agents := models.Agents{}
			agents.CreateAndSaveRegistryExamples(db, 3)

			// make a request
			req, err := http.NewRequest("GET", "/agents", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			dbAgents := make(models.Agents, 0)
			err = json.Unmarshal(w.Body.Bytes(), &dbAgents)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))
			Expect(dbAgents[2].AgentID).To(Equal(agents[2].AgentID))
		})

	})

	Describe("serveAgent", func() {

		var (
			agent models.Agent
		)

		JustBeforeEach(func() {
			// build a request
			req := models.Request{}
			req.RegistryExample()

			// save a fact
			fact := models.Fact{}
			err := fact.Update(db, &req.Request)
			Expect(err).NotTo(HaveOccurred())
			agent = fact.Agent
		})

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := http.NewRequest("GET", "/agents/non_exisitng_id", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 500 error if something is wrong with the query", func() {
			// wrong query
			tmp_GetAgentQuery := GetAgentQuery
			GetAgentQuery = "SELECT agent_id,created_at,updated_at FROM Wrong_facts_table WHERE agent_id=$1"

			// make request
			req, err := http.NewRequest("GET", fmt.Sprint("/agents/", agent.AgentID), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			// copy back the query defenition
			GetAgentQuery = tmp_GetAgentQuery
		})

		It("return an angent", func() {
			// make request
			req, err := http.NewRequest("GET", fmt.Sprint("/agents/", agent.AgentID), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			var dbAgent models.Agent
			err = json.Unmarshal(w.Body.Bytes(), &dbAgent)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgent.AgentID).To(Equal(agent.AgentID))
		})

	})

	Describe("serveFacts", func() {

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := http.NewRequest("GET", "/agents/non_existing_id/facts", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 500 error if something is wrong with the query", func() {
			tmp_GetFactQuery := GetFactQuery
			GetFactQuery = "SELECT * FROM Wrong_facts_table WHERE agent_id=$1"

			// make request
			req, err := http.NewRequest("GET", "/agents/non_existing_id/facts", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			GetFactQuery = tmp_GetFactQuery
		})

		It("returns the facts from an agent", func() {
			// build a request
			regReq := models.Request{}
			regReq.RegistryExample()

			// save a fact
			fact := models.Fact{}
			err := fact.Update(db, &regReq.Request)
			Expect(err).NotTo(HaveOccurred())
			agent := fact.Agent

			// make request
			req, err := http.NewRequest("GET", fmt.Sprint("/agents/", agent.AgentID, "/facts"), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			Expect(w.Body.String()).To(Equal(fact.Facts))
		})

	})

})

var _ = Describe("Log Handlers", func() {

	var (
		job models.Job
	)

	JustBeforeEach(func() {
		// save a job
		job = models.Job{}
		job.RpcVersionExample()
		err := job.Save(db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("serveJobLog", func() {

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := http.NewRequest("GET", "/jobs/non_existing_id/log", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 500 error if something is wrong with the query", func() {
			// wrong query
			tmp_GetLogQuery := GetLogQuery
			GetLogQuery = "SELECT * Wrong_Log_Table logs WHERE job_id=$1"

			// make request
			req, err := http.NewRequest("GET", fmt.Sprint("/jobs/", job.RequestID, "/log"), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			// copy back the query defenition
			GetLogQuery = tmp_GetLogQuery
		})

		It("returns the log from the log table", func() {
			// save log for the job
			reply := models.Reply{}
			reply.ExecuteScriptExample(job.RequestID, true, "Log text", 1)
			err := models.ProcessLogReply(db, &reply.Reply)
			Expect(err).NotTo(HaveOccurred())

			// make request
			req, err := http.NewRequest("GET", fmt.Sprint("/jobs/", job.RequestID, "/log"), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			Expect(w.Body.String()).To(Equal(reply.Payload))
		})

		It("returns log chuncks", func() {
			// save log chunks
			logpart := models.LogPart{}
			content := logpart.SaveLogPartExamples(db, job.RequestID)

			// make request
			req, err := http.NewRequest("GET", fmt.Sprint("/jobs/", job.RequestID, "/log"), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(200))

			// check json body response
			Expect(w.Body.String()).To(Equal(content))
		})

	})

})

var _ = Describe("root", func() {

})
