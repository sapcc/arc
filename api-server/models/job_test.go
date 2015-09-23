// +build integration

package models_test

import (
	"fmt"
	"time"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"
)

var _ = Describe("Jobs", func() {

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			job := Job{}
			job.ExecuteScriptExample()
			err := job.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should return all jobs", func() {
			jobs := Jobs{}
			jobs.CreateAndSaveRpcVersionExamples(db, 3)

			dbJobs := Jobs{}
			err := dbJobs.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJobs[0].RequestID).To(Equal(jobs[0].RequestID))
			Expect(dbJobs[1].RequestID).To(Equal(jobs[1].RequestID))
			Expect(dbJobs[2].RequestID).To(Equal(jobs[2].RequestID))
		})

	})

})

var _ = Describe("Job", func() {

	Describe("CreateJob", func() {

		It("returns an error data is no json conform", func() {
			noValidJson := `"to":"darwin"`
			strSlice := []byte(noValidJson)
			job, err := CreateJob(&strSlice, uuid.New())
			Expect(err).To(HaveOccurred())
			var newJob *Job
			Expect(job).To(Equal(newJob))
		})

		It("returns an error data is not valid", func() {
			noValidData := `{"to":"darwin","timeout":60,"agent":"execute","payload":"echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""}` // action is missing
			strSlice := []byte(noValidData)
			job, err := CreateJob(&strSlice, uuid.New())
			Expect(err).To(HaveOccurred())
			var newJob *Job
			Expect(job).To(Equal(newJob))
		})

		It("should return a job", func() {
			to := "darwin"
			timeout := 60
			agent := "execute"
			action := "script"
			payload := `"payload":"echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""`
			noValidData := fmt.Sprintf(`{"to":%q,"timeout":%v,"agent":%q,"action":%q,"payload":%q}`, to, timeout, agent, action, payload)
			strSlice := []byte(noValidData)
			job, err := CreateJob(&strSlice, uuid.New())
			Expect(err).NotTo(HaveOccurred())
			Expect(job.To).To(Equal(to))
			Expect(job.Timeout).To(Equal(timeout))
			Expect(job.Agent).To(Equal(agent))
			Expect(job.Action).To(Equal(action))
			Expect(job.Payload).To(Equal(payload))
		})

	})

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			job := Job{}
			job.ExecuteScriptExample()
			err := job.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if job not found", func() {
			job := Job{}
			job.ExecuteScriptExample()
			err := job.Get(db)
			Expect(err).To(HaveOccurred())
		})

		It("should return the job", func() {
			// create and save a job
			job := Job{}
			job.ExecuteScriptExample()
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// get the job
			dbJob := Job{Request: arc.Request{RequestID: job.RequestID}}
			err = dbJob.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJob.RequestID).To(Equal(job.RequestID))
			Expect(dbJob.To).To(Equal(job.To))
			Expect(dbJob.Timeout).To(Equal(job.Timeout))
			Expect(dbJob.Agent).To(Equal(job.Agent))
			Expect(dbJob.Action).To(Equal(job.Action))
			Expect(dbJob.Payload).To(Equal(job.Payload))
		})

	})

	Describe("Save", func() {

		It("returns an error if no db connection is given", func() {
			job := Job{}
			job.ExecuteScriptExample()
			err := job.Save(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should save a job", func() {
			job := Job{}
			job.ExecuteScriptExample()
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// get the job and check
			dbJob := Job{Request: arc.Request{RequestID: job.RequestID}}
			err = dbJob.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJob.RequestID).To(Equal(job.RequestID))
			Expect(dbJob.To).To(Equal(job.To))
			Expect(dbJob.Timeout).To(Equal(job.Timeout))
			Expect(dbJob.Agent).To(Equal(job.Agent))
			Expect(dbJob.Action).To(Equal(job.Action))
			Expect(dbJob.Payload).To(Equal(job.Payload))
		})

	})

	Describe("Update", func() {

		It("returns an error if no db connection is given", func() {
			job := Job{}
			job.ExecuteScriptExample()
			err := job.Update(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should update a job", func() {
			// save a job
			job := Job{}
			job.ExecuteScriptExample()
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// update a job and check
			status := arc.Complete
			updated_at := time.Now()
			newJob := Job{Request: arc.Request{RequestID: job.RequestID}, Status: status, UpdatedAt: updated_at}
			err = newJob.Update(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(newJob.RequestID).To(Equal(job.RequestID))
			Expect(newJob.Status).To(Equal(status))
			Expect(newJob.UpdatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(updated_at.Format("2006-01-02 15:04:05.99")))
		})

	})

	Describe("CleanJobs", func() {

		It("returns an error if no db connection is given", func() {
			_, _, err := CleanJobs(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should clean jobs which no heartbeat was send back after created_at + 60 sec", func() {
			// save a job
			job := Job{}
			job.CustomExecuteScriptExample(arc.Queued, time.Now().Add(-61*time.Second), 120)
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// clean jobs
			occurHeartbeat, occurTimeOut, err := CleanJobs(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(occurHeartbeat).To(Equal(int64(1)))
			Expect(occurTimeOut).To(Equal(int64(0)))

			// check job
			dbJob := Job{Request: arc.Request{RequestID: job.RequestID}}
			err = dbJob.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJob.Status).To(Equal(arc.Failed))
		})

		It("should clean jobs which the timeout + 60 sec has exceeded and still in queued or executing status", func() {
			// save a job
			job := Job{}
			job.CustomExecuteScriptExample(arc.Executing, time.Now().Add((-20-60)*time.Second), 15) // 60 sec extra to be sure
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			job2 := Job{}
			job2.CustomExecuteScriptExample(arc.Queued, time.Now().Add((-20-60)*time.Second), 15) // 60 sec extra to be sure
			err = job2.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// clean jobs
			occurHeartbeat, occurTimeOut, err := CleanJobs(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(occurHeartbeat).To(Equal(int64(1))) // this overlap between hearbeat and timeout criteria
			Expect(occurTimeOut).To(Equal(int64(1)))

			// check job
			dbJob := Job{Request: arc.Request{RequestID: job.RequestID}}
			err = dbJob.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJob.Status).To(Equal(arc.Failed))

			dbJob2 := Job{Request: arc.Request{RequestID: job2.RequestID}}
			err = dbJob2.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJob2.Status).To(Equal(arc.Failed))
		})

	})

})
