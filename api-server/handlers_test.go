// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/version"
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
			req, err := newAuthorizedRequest("GET", getUrl("/jobs"), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code, header and body
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("GET", getUrl("/jobs"), "")

			// make a request with X-Identity-Status to Confirmed but not X-Project-Id
			req, err := http.NewRequest("GET", getUrl("/jobs"), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Add("X-Identity-Status", `Confirmed`)
			req.Header.Add("X-Project-Id", `some_different_project`)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))
			// check no jobs returned
			dbJobs := make(models.Jobs, 0)
			err = json.Unmarshal(w.Body.Bytes(), &dbJobs)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbJobs)).To(Equal(0))
		})

		It("returns empty arry if no jobs found", func() {
			// make a request
			req, err := newAuthorizedRequest("GET", getUrl("/jobs"), bytes.NewBufferString(""), map[string]string{})
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
			req, err := newAuthorizedRequest("GET", getUrl("/jobs"), bytes.NewBufferString(""), map[string]string{})
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

			Expect(dbJobs[0].RequestID).To(Equal(jobs[2].RequestID))
			Expect(dbJobs[1].RequestID).To(Equal(jobs[1].RequestID))
			Expect(dbJobs[2].RequestID).To(Equal(jobs[0].RequestID))
		})

		It("returns jobs filtered by agent_id", func() {
			// fill db
			jobs := models.Jobs{}
			jobs.CreateAndSaveRpcVersionExamples(db, 3)

			// add job with special target
			job := models.Job{}
			job.ExecuteScriptExample()
			job.To = "my_test_laptop"
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// make request
			req, err := newAuthorizedRequest("GET", getUrl("/jobs?agent_id=my_test_laptop"), bytes.NewBufferString(""), map[string]string{})
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
			Expect(len(dbJobs)).To(Equal(1))
			Expect(dbJobs[0].RequestID).To(Equal(job.RequestID))
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

		It("returns a 500 error if something is wrong with the query", func() {
			// wrong query
			tmp_GetJobQuery := GetJobQuery
			GetJobQuery = "SELECT * Wrong_Job_Table jobs WHERE id=$1"

			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID)), bytes.NewBufferString(""), map[string]string{})
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
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID)), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))
		})

		It("returns a 404 error if job not found", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl("/jobs/non_existing_id"), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID)), "")
			checkNonAuthorizeProjectRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID)), "")
		})

		It("should return the job", func() {
			// make a request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID)), bytes.NewBufferString(""), map[string]string{})
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
			// setting up the transport
			config.Identity = "darwin"
			config.Transport = "fake"
			var err error
			tp, err = arcNewConnection(config)
			Expect(err).NotTo(HaveOccurred())
			// create an agent
			agent := models.Agent{}
			agent.Example()
			agent.AgentID = "darwin"
			err = agent.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a 400 error if request is wrong", func() {
			jsonStr := []byte(`this is not json`)
			// make a request
			req, err := newAuthorizedRequest("POST", getUrl("/jobs"), bytes.NewBuffer(jsonStr), map[string]string{})
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
			req, err := newAuthorizedRequest("POST", getUrl("/jobs"), bytes.NewBuffer(jsonStr), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			// set the pointer back to the db
			db = db_copy
		})

		It("returns a 404 error if the targe agent does not exist", func() {
			jsonStr := []byte(`{"to":"non_existing_agent","timeout":60,"agent":"rpc","action":"version"}`)
			// make a request
			req, err := newAuthorizedRequest("POST", getUrl("/jobs"), bytes.NewBuffer(jsonStr), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("POST", getUrl("/jobs"), `{"to":"darwin","timeout":60,"agent":"rpc","action":"version"}`)
			checkNonAuthorizeProjectRequest("POST", getUrl("/jobs"), `{"to":"darwin","timeout":60,"agent":"rpc","action":"version"}`)
		})

		It("should save the job and return the unique id as JSON", func() {
			jsonStr := []byte(`{"to":"darwin","timeout":60,"agent":"rpc","action":"version"}`)
			// make a request
			req, err := newAuthorizedRequest("POST", getUrl("/jobs"), bytes.NewBuffer(jsonStr), map[string]string{})
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

var _ = Describe("Agent Handlers", func() {

	Describe("serveAgents", func() {

		var (
			agents models.Agents
		)

		It("returns a 500 error if something goes wrong", func() {
			tmp_GetAgentsQuery := GetAgentsQuery
			GetAgentsQuery = "SELECT DISTINCT agent_id,created_at,updated_at FROM Wrong_Facts order by updated_at"

			// make a request
			req, err := newAuthorizedRequest("GET", getUrl("/agents"), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			GetAgentsQuery = tmp_GetAgentsQuery
		})

		It("returns a 400 if the filter query is wrong", func() {
			// make a request
			req, err := newAuthorizedRequest("GET", getUrl(`/agents?q=os+%3D`), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(400))
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("GET", getUrl("/agents"), "")

			// make a request with X-Identity-Status to Confirmed but not X-Project-Id
			agents := models.Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)
			req, err := http.NewRequest("GET", getUrl("/agents"), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Add("X-Identity-Status", `Confirmed`)
			req.Header.Add("X-Project-Id", `some_different_project`)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))
			// check no agents returned
			dbAgents := make(models.Agents, 0)
			err = json.Unmarshal(w.Body.Bytes(), &dbAgents)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(0))
		})

		It("returns empty json arry if no agents found", func() {
			// make a request
			req, err := newAuthorizedRequest("GET", getUrl("/agents"), bytes.NewBufferString(""), map[string]string{})
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
			agents.CreateAndSaveAgentExamples(db, 3)

			// make a request
			req, err := newAuthorizedRequest("GET", getUrl("/agents"), bytes.NewBufferString(""), map[string]string{})
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

			Expect(dbAgents[0].AgentID).To(Equal(agents[2].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))
			Expect(dbAgents[2].AgentID).To(Equal(agents[0].AgentID))
		})

		It("returns all agents filtered", func() {
			var (
				facts      = `{"os": "%s", "online": "%s", "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
				os         = []string{"darwin", "windows", "windows"}
				online     = []string{"true", "false", "true"}
				tagsKey    = []string{"landscape", "landscape", "landscape"}
				tagsValue  = []string{"development", "staging", "production"}
				tagsKey2   = []string{"pool", "pool", "pool"}
				tagsValue2 = []string{"green", "green", "blue"}
			)
			agents := models.Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)
			for i := 0; i < len(agents); i++ {
				currentAgent := agents[i]
				// change facts
				err := json.Unmarshal([]byte(fmt.Sprintf(facts, os[i], online[i])), &currentAgent.Facts)
				Expect(err).NotTo(HaveOccurred())
				err = currentAgent.Update(db)
				Expect(err).NotTo(HaveOccurred())
				// add tags
				authorization := auth.Authorization{IdentityStatus: "Confirmed", UserId: "userID", ProjectId: currentAgent.Project}
				err = currentAgent.AddTagAuthorized(db, &authorization, tagsKey[i], tagsValue[i])
				Expect(err).NotTo(HaveOccurred())
				err = currentAgent.AddTagAuthorized(db, &authorization, tagsKey2[i], tagsValue2[i])
				Expect(err).NotTo(HaveOccurred())
			}

			// make a request
			req, err := newAuthorizedRequest("GET", getUrl(`/agents?q=@os+%3D+%22darwin%22+OR+%28landscape+%3D+%22staging%22+AND+pool+%3D+%22green%22%29`), bytes.NewBufferString(""), map[string]string{})
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
			Expect(len(dbAgents)).To(Equal(2))
			Expect(dbAgents[0].AgentID).To(Equal(agents[1].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[0].AgentID))
		})

		It("should show facts", func() {
			facts := `{"os": "%s", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
			os := []string{"darwin", "windows", "windows"}
			agents := models.Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)
			for i := 0; i < len(agents); i++ {
				currentAgent := agents[i]
				if err := json.Unmarshal([]byte(fmt.Sprintf(facts, os[i])), &currentAgent.Facts); err != nil {
					Expect(err).NotTo(HaveOccurred())
				}
				err := currentAgent.Update(db)
				Expect(err).NotTo(HaveOccurred())
			}

			// make a request
			req, err := newAuthorizedRequest("GET", getUrl(`/agents?q=@os+%3D+"windows"&facts=os,online`), bytes.NewBufferString(""), map[string]string{})
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
			Expect(len(dbAgents)).To(Equal(2))
			for i := 0; i < len(dbAgents); i++ {
				currentAgent := dbAgents[i]
				Expect(len(currentAgent.Facts)).To(Equal(2))
				_, ok := currentAgent.Facts["os"]
				Expect(ok).To(Equal(true))
				_, ok = currentAgent.Facts["online"]
				Expect(ok).To(Equal(true))
			}
		})

	})

	Describe("serveAgent", func() {

		var (
			agent models.Agent
		)

		JustBeforeEach(func() {
			agent = models.Agent{}
			agent.Example()
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl("/agents/non_exisitng_id"), bytes.NewBufferString(""), map[string]string{})
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
			GetAgentQuery = "SELECT * FROM Wrong_facts_table WHERE agent_id=$1"

			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID)), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			// copy back the query defenition
			GetAgentQuery = tmp_GetAgentQuery
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID)), "")
			checkNonAuthorizeProjectRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID)), "")
		})

		It("return an angent", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID)), bytes.NewBufferString(""), map[string]string{})
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

			// check no facts are shown in the json response
			var objmap map[string]*json.RawMessage
			err = json.Unmarshal(w.Body.Bytes(), &objmap)
			Expect(err).NotTo(HaveOccurred())
			var nilJson *json.RawMessage
			Expect(objmap["facts"]).To(Equal(nilJson))
		})

		It("should show facts", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "?facts=os,online")), bytes.NewBufferString(""), map[string]string{})
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
			Expect(len(dbAgent.Facts)).To(Equal(2))
			_, ok := dbAgent.Facts["os"]
			Expect(ok).To(Equal(true))
			_, ok = dbAgent.Facts["online"]
			Expect(ok).To(Equal(true))
		})

	})

	Describe("deleteAgent", func() {

		var (
			agent models.Agent
		)

		JustBeforeEach(func() {
			agent = models.Agent{}
			agent.Example()
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := newAuthorizedRequest("DELETE", getUrl("/agents/non_exisitng_id"), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 500 error if something is wrong happens", func() {
			// wrong query
			tmp_DeleteAgentQuery := DeleteAgentQuery
			DeleteAgentQuery = `DELETE FROM miaus WHERE agent_id=$1`

			// make request
			req, err := newAuthorizedRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID)), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			// copy back the query defenition
			DeleteAgentQuery = tmp_DeleteAgentQuery
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID)), "")
			checkNonAuthorizeProjectRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID)), "")
		})

		It("Delete an angent", func() {
			// make request
			req, err := newAuthorizedRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID)), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(200))
		})

	})

	Describe("serveFacts", func() {

		var (
			agent models.Agent
		)

		JustBeforeEach(func() {
			agent = models.Agent{}
			agent.Example()
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a 500 error if something is wrong with the query", func() {
			tmp_GetAgentQuery := GetAgentQuery
			GetAgentQuery = "SELECT * FROM Wrong_facts_table WHERE agent_id=$1"

			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/facts")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			GetAgentQuery = tmp_GetAgentQuery
		})

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl("/agents/non_existing_id/facts"), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/facts")), "")
			checkNonAuthorizeProjectRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/facts")), "")
		})

		It("returns the facts from an agent", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/facts")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check facts
			checkFacts := models.JSONBfromString(w.Body.String())
			eq := reflect.DeepEqual(agent.Facts, checkFacts)
			Expect(eq).To(Equal(true))
		})

	})

})

var _ = Describe("Tags", func() {

	var (
		agent models.Agent
	)

	JustBeforeEach(func() {
		agent = models.Agent{}
		agent.Example()
		err := agent.Save(db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("serveAgentTags", func() {

		It("returns a 500 error if something is wrong with the query", func() {
			tmp_GetAgentQuery := GetAgentQuery
			GetAgentQuery = `SELECT DISTINCT * FROM wrong_table WHERE agent_id=$1 order by created_at DESC`

			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			GetAgentQuery = tmp_GetAgentQuery
		})

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", "non_exisitng_agent", "/tags")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns an empty json if no tags for the agent", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check body is json empty
			var objmap *json.RawMessage
			err = json.Unmarshal(w.Body.Bytes(), &objmap)
			Expect(err).NotTo(HaveOccurred())
			data := []byte(`{}`)
			nilJson := (*json.RawMessage)(&data)
			Expect(objmap).To(Equal(nilJson))
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), "")
			checkNonAuthorizeProjectRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), "")
		})

		It("returns the tags from an agent", func() {
			authorization := auth.Authorization{IdentityStatus: "Confirmed", UserId: "userID", ProjectId: agent.Project}
			err := agent.AddTagAuthorized(db, &authorization, "cat", "miau")
			Expect(err).NotTo(HaveOccurred())
			err = agent.AddTagAuthorized(db, &authorization, "dog", "bup")
			Expect(err).NotTo(HaveOccurred())

			// create a new agent with other project and tags
			newAgent := models.Agent{}
			newAgent.Project = "another_project"
			newAgent.Example()
			err = newAgent.Save(db)
			Expect(err).NotTo(HaveOccurred())
			err = newAgent.AddTagAuthorized(db, &authorization, "bird", "piupiu")
			Expect(err).NotTo(HaveOccurred())

			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))

			// check tags
			dbAgent := models.Agent{AgentID: agent.AgentID}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			checkTags := models.JSONBfromString(w.Body.String())
			eq := reflect.DeepEqual(dbAgent.Tags, checkTags)
			Expect(eq).To(Equal(true))
		})

	})

	Describe("saveAgentTags", func() {

		It("returns a 500 error if something is wrong with the query", func() {
			tmp_AddAgentTag := AddAgentTag
			AddAgentTag = `INSERT INTO wrong_table(agent_id,project,value,created_at) VALUES($1,$2,$3) returning agent_id`

			// make request
			req, err := newAuthorizedRequest("POST", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), bytes.NewBufferString("cat=miau&dog=bup"), map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			AddAgentTag = tmp_AddAgentTag
		})

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := newAuthorizedRequest("POST", getUrl(fmt.Sprint("/agents/", "non_existing_agent", "/tags")), bytes.NewBufferString("cat=miau&dog=bup"), map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("POST", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), "tag1, tag2")
			checkNonAuthorizeProjectRequest("POST", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), "tag1, tag2")
		})

		It("saves all tags for the agent", func() {
			// make request
			req, err := newAuthorizedRequest("POST", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags")), bytes.NewBufferString("cat=miau&dog=bup"), map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(200))

			dbAgent := models.Agent{AgentID: agent.AgentID}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			checkTags := models.JSONBfromString(`{"cat":"miau","dog":"bup"}`)
			eq := reflect.DeepEqual(dbAgent.Tags, checkTags)
			Expect(eq).To(Equal(true))
		})

	})

	Describe("deleteAgentTag", func() {

		It("returns a 500 error if something is wrong with the query", func() {
			tmp_DeleteTagQuery := DeleteAgentTagQuery
			DeleteAgentTagQuery = `DELETE FROM wrong_table WHERE agent_id=$1 AND value=$2`

			authorization := auth.Authorization{IdentityStatus: "Confirmed", UserId: "userID", ProjectId: agent.Project}
			err := agent.AddTagAuthorized(db, &authorization, "cat", "miau")
			Expect(err).NotTo(HaveOccurred())

			// make request
			req, err := newAuthorizedRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags", "/cat")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			DeleteAgentTagQuery = tmp_DeleteTagQuery
		})

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := newAuthorizedRequest("DELETE", getUrl(fmt.Sprint("/agents/", "non_existing_agent", "/tags", "/tag_miau")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns no error if Tag not found", func() {
			// make request
			req, err := newAuthorizedRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags", "/non_exiting_tag")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(200))
		})

		It("returns a 401 error if not authorized", func() {
			checkIdentityInvalidRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags", "/tag_miau")), "")
			checkNonAuthorizeProjectRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags", "/tag_miau")), "")
		})

		It("removes the agent tag", func() {
			authorization := auth.Authorization{IdentityStatus: "Confirmed", UserId: "userID", ProjectId: agent.Project}
			err := agent.AddTagAuthorized(db, &authorization, "cat", "miau")
			Expect(err).NotTo(HaveOccurred())
			err = agent.AddTagAuthorized(db, &authorization, "dog", "bup")
			Expect(err).NotTo(HaveOccurred())

			// make request
			req, err := newAuthorizedRequest("DELETE", getUrl(fmt.Sprint("/agents/", agent.AgentID, "/tags", "/dog")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(200))

			dbAgent := models.Agent{AgentID: agent.AgentID}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			checkTags := models.JSONBfromString(`{"cat":"miau"}`)
			eq := reflect.DeepEqual(dbAgent.Tags, checkTags)
			Expect(eq).To(Equal(true))
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

		It("returns a 500 error if something is wrong with the query", func() {
			// wrong query
			tmp_GetLogQuery := GetLogQuery
			GetLogQuery = "SELECT * Wrong_Log_Table logs WHERE job_id=$1"

			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID, "/log")), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))

			// copy back the query defenition
			GetLogQuery = tmp_GetLogQuery
		})

		It("returns a 404 error if Agent not found", func() {
			// make request
			req, err := newAuthorizedRequest("GET", getUrl("/jobs/non_existing_id/log"), bytes.NewBufferString(""), map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 401 error if not authorized", func() {
			// save log for the job
			reply := models.Reply{}
			reply.ExecuteScriptExample(job.RequestID, true, "Log text", 1)
			err := models.ProcessLogReply(db, &reply.Reply, "darwin", true)
			Expect(err).NotTo(HaveOccurred())

			checkIdentityInvalidRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID, "/log")), "")
			checkNonAuthorizeProjectRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID, "/log")), "")
		})

		It("returns the log from the log table", func() {
			// save log for the job
			reply := models.Reply{}
			reply.ExecuteScriptExample(job.RequestID, true, "Log text", 1)
			err := models.ProcessLogReply(db, &reply.Reply, "darwin", true)
			Expect(err).NotTo(HaveOccurred())

			// make request
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID, "/log")), bytes.NewBufferString(""), map[string]string{})
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
			req, err := newAuthorizedRequest("GET", getUrl(fmt.Sprint("/jobs/", job.RequestID, "/log")), bytes.NewBufferString(""), map[string]string{})
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

var _ = Describe("Root Handler", func() {

	It("returns the app name and version as plain text", func() {
		// make request
		req, err := http.NewRequest("GET", "/", bytes.NewBufferString(""))
		Expect(err).NotTo(HaveOccurred())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// check response code and header
		Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
		Expect(w.Code).To(Equal(200))

		// check json body response
		Expect(w.Body.String()).To(Equal(fmt.Sprint("Arc api-server ", version.String())))
	})

})

var _ = Describe("Healthcheck Handler", func() {

	It("returns the app name and version as plain text", func() {
		// make request
		req, err := http.NewRequest("GET", "/healthcheck", bytes.NewBufferString(""))
		Expect(err).NotTo(HaveOccurred())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// check response code and header
		Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
		Expect(w.Code).To(Equal(200))

		// check json body response
		Expect(w.Body.String()).To(Equal(fmt.Sprint("Arc api-server ", version.String())))
	})

})

var _ = Describe("Readiness Handler", func() {

	Describe("DB not reachable", func() {

		JustBeforeEach(func() {
			db.Close()
		})

		AfterEach(func() {
			var err error
			env := os.Getenv("ARC_ENV")
			if env == "" {
				env = "test"
			}
			db, err = NewConnection("db/dbconf.yml", env)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 502 if the db is not reachable", func() {
			// make request
			req, err := http.NewRequest("GET", "/readiness", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(502))

			// check json body response
			var jsonBody Readiness
			err = json.Unmarshal(w.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonBody.Status).To(Equal(502))
			Expect(jsonBody.Message).To(ContainSubstring("DB"))
		})

	})

	Describe("Lost MQTT connection", func() {

		JustBeforeEach(func() {
			config.Identity = "darwin"
			config.Transport = "fake"
			config.Organization = "no-connected"
			var err error
			tp, err = arcNewConnection(config)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			config.Identity = "darwin"
			config.Transport = "fake"
			config.Organization = ""
			var err error
			tp, err = arcNewConnection(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 502 if the connection to MQTT is broken", func() {
			// make request
			req, err := http.NewRequest("GET", "/readiness", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(502))

			// check json body response
			var jsonBody Readiness
			err = json.Unmarshal(w.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonBody.Status).To(Equal(502))
			Expect(jsonBody.Message).To(ContainSubstring("transport"))
		})

	})

	It("returns 200 if DB and MQTT are reachable or connected", func() {
		// make request
		req, err := http.NewRequest("GET", "/readiness", bytes.NewBufferString(""))
		Expect(err).NotTo(HaveOccurred())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// check response code and header
		Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
		Expect(w.Code).To(Equal(200))
	})

})

// private

func checkIdentityInvalidRequest(method, urlStr string, body string) {
	// make a request without any authorization header
	req, err := http.NewRequest(method, urlStr, bytes.NewBufferString(body))
	Expect(err).NotTo(HaveOccurred())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
	Expect(w.Code).To(Equal(401))

	// make a request with X-Identity-Status different from Confirmed
	req, err = http.NewRequest(method, urlStr, bytes.NewBufferString(body))
	Expect(err).NotTo(HaveOccurred())
	req.Header.Add("X-Identity-Status", `something_different_from_confirmed`)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
	Expect(w.Code).To(Equal(401))
}

func checkNonAuthorizeProjectRequest(method, urlStr string, body string) {
	// make a request with X-Identity-Status to Confirmed and X-Project-Id with a different project
	req, err := http.NewRequest(method, urlStr, bytes.NewBufferString(body))
	req.Header.Add("X-Identity-Status", `Confirmed`)
	req.Header.Add("X-Project-Id", `some_different_project`)
	Expect(err).NotTo(HaveOccurred())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
	Expect(w.Code).To(Equal(401))
}

func newAuthorizedRequest(method, urlStr string, body io.Reader, headerOptions map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Identity-Status", `Confirmed`)
	req.Header.Add("X-Project-Id", `test-project`)

	// add extra headers
	for k, v := range headerOptions {
		req.Header.Add(k, v)
	}

	return req, nil
}

func getUrl(url string) string {
	return fmt.Sprint("/api/v1", url)
}
