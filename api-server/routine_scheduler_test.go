// +build integration

package main

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"
)

var _ = Describe("Routine scheduler", func() {

	It("should return an error if the db connection is nil", func() {

	})

	It("should clean jobs", func() {
		// save a job
		job := Job{}
		job.CustomExecuteScriptExample(arc.Queued, time.Now().Add(-61*time.Second), 120)
		err := job.Save(db)
		Expect(err).NotTo(HaveOccurred())

		//runRoutineTasks(db)

		// check job
		dbJob := Job{Request: arc.Request{RequestID: job.RequestID}}
		err = dbJob.Get(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbJob.Status).To(Equal(arc.Failed))
	})

	It("should clean log parts", func() {
		// add a job related to the log chuncks
		job := Job{}
		job.ExecuteScriptExample()
		job.Save(db)

		// log part
		logPart := LogPart{job.RequestID, 1, "Some chunk of code", true, time.Now().Add(-601 * time.Second)} // bit more than 10 min
		err := logPart.Save(db)
		Expect(err).NotTo(HaveOccurred())

		runRoutineTasks(db)

		// check log parts
		dbLogPart := LogPart{JobID: job.RequestID}
		_, err = dbLogPart.Collect(db)
		Expect(err).To(HaveOccurred())

		// check log
		dbLog := Log{JobID: job.RequestID}
		err = dbLog.Get(db)
		Expect(err).NotTo(HaveOccurred())
	})

})
