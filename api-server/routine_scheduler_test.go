// +build integration

package main

import (
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"
)

var _ = Describe("Job Handlers", func() {

	It("should clean jobs which no heartbeat was send back after created_at + 60 sec", func() {
		// save a job
		job := models.Job{}
		job.CustomExecuteScriptExample(arc.Queued, time.Now().Add(-61*time.Second), 120)
		err := job.Save(db)
		Expect(err).NotTo(HaveOccurred())

		// clean jobs
		cleanJobs(db)

		// check job
		dbJob := models.Job{Request: arc.Request{RequestID: job.RequestID}}
		err = dbJob.Get(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbJob.Status).To(Equal(arc.Failed))
	})

	It("should clean jobs which the timeout + 60 sec has exceeded and still in queued or executing status", func() {
		// save a job
		job := models.Job{}
		job.CustomExecuteScriptExample(arc.Executing, time.Now().Add((-20-60)*time.Second), 15) // 60 sec extra to be sure
		err := job.Save(db)
		Expect(err).NotTo(HaveOccurred())

		job2 := models.Job{}
		job2.CustomExecuteScriptExample(arc.Queued, time.Now().Add((-20-60)*time.Second), 15) // 60 sec extra to be sure
		err = job2.Save(db)
		Expect(err).NotTo(HaveOccurred())

		// clean jobs
		cleanJobs(db)

		// check job
		dbJob := models.Job{Request: arc.Request{RequestID: job.RequestID}}
		err = dbJob.Get(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbJob.Status).To(Equal(arc.Failed))

		dbJob2 := models.Job{Request: arc.Request{RequestID: job2.RequestID}}
		err = dbJob2.Get(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbJob2.Status).To(Equal(arc.Failed))
	})

})
