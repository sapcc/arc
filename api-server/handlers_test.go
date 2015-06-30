package main

import (	
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	
	"net/http/httptest"
	"encoding/json"	
	"bytes"
	"net/http"
)

var _ = Describe("Handlers", func() {

	Describe("serveJobs", func() {
	
		It("returns a 500 error if something goes wrong", func() {})
		
		It("returns empty json arry if no jobs found", func() {
			req, err := http.NewRequest("POST", "/jobs", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())		
			w := httptest.NewRecorder()
			serveJobs(w, req)
			
			// check response code
			Expect(w.Code).To(Equal(200))
			
			// check json body response						
			jobs := make(models.Jobs,0)
			err = json.Unmarshal(w.Body.Bytes(), &jobs)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobs).To(Equal(make(models.Jobs,0)))			
		})
		
		It("returns all jobs", func() {
			// fill db
			jobs := []models.Job{models.ExecuteSctiptJob(), models.ExecuteSctiptJob(), models.ExecuteSctiptJob()}
			for i := 0; i < len(jobs); i++ {
				job := jobs[i]
				err := job.Save(db);
				Expect(err).NotTo(HaveOccurred())			
			}
			
			// make request
			req, err := http.NewRequest("POST", "/jobs", bytes.NewBufferString(""))
			Expect(err).NotTo(HaveOccurred())			
			w := httptest.NewRecorder()
			serveJobs(w, req)
			
			// check response code and header
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(w.Code).To(Equal(200))						
			
			// check json body response						
			dbJobs := make(models.Jobs,0)
			err = json.Unmarshal(w.Body.Bytes(), &dbJobs)
			Expect(err).NotTo(HaveOccurred())						
			Expect(dbJobs[0].RequestID).To(Equal(jobs[0].RequestID))
			Expect(dbJobs[1].RequestID).To(Equal(jobs[1].RequestID))
			Expect(dbJobs[2].RequestID).To(Equal(jobs[2].RequestID))									
			
		})		
		
	})
	
	Describe("serveJob", func() {})
	Describe("executeJob", func() {})
	Describe("serveJobLog", func() {})
	Describe("serveAgents", func() {})
	Describe("serveAgent", func() {})
	Describe("serveFacts", func() {})
	Describe("root", func() {})

})
