package main

import (
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Handlers", func() {

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
			serveJobs(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))
		})

		It("returns empty json arry if no jobs found", func() {
			// make a request
			req, err := http.NewRequest("POST", "/jobs", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			serveJobs(w, req)

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
			serveJobs(w, req)

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

		It("returns a 404 error if job not found", func() {
			// make request
			req, err := http.NewRequest("GET", "/jobs/non_existing_id", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			serveJob(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(404))
		})

		It("returns a 500 error if something goes wrong", func() {
			// save bad data
			job := models.Job{}
			job.RpcVersionExample()
			job.Status = 6 // not existing status
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// make a request
			req, err := http.NewRequest("GET", fmt.Sprint("/jobs/", job.RequestID), bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			router.HandleFunc("/employees/{id}", serveJob)
			router.ServeHTTP(w, req)

			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			Expect(w.Code).To(Equal(500))
		})

		It("should return the job", func() {})
	})

	Describe("executeJob", func() {})
	Describe("serveJobLog", func() {})
	Describe("serveAgents", func() {})
	Describe("serveAgent", func() {})
	Describe("serveFacts", func() {})
	Describe("root", func() {})

})
