// +build integration

package models_test

import (
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"

	. "github.com/sapcc/arc/api-server/models"
)

var _ = Describe("LogParts", func() {

	Describe("Collect", func() {

		It("returns an error if no db connection is given", func() {
			logPart := LogPart{JobID: "the_job_id"}
			_, err := logPart.Collect(nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if no id found", func() {
			logPart := LogPart{JobID: "the_job_id"}
			dbContent, err := logPart.Collect(db)
			var res *string
			Expect(dbContent).To(Equal(res))
			Expect(err).To(HaveOccurred())
		})

		It("should collect all log chuncks", func() {
			// add a job related to the log chuncks
			newJob := Job{}
			newJob.ExecuteScriptExample()
			newJob.Save(db)
			// save different chuncks
			var contentSlice [3]string
			for i := 0; i < 3; i++ {
				chunck := fmt.Sprintf("This is the %d chunck", i)
				contentSlice[i] = chunck
				logPart := LogPart{newJob.RequestID, uint(i), chunck, false, time.Now()}
				err := logPart.Save(db)
				Expect(err).NotTo(HaveOccurred())
			}
			content := strings.Join(contentSlice[:], "")

			// collect
			logPart := LogPart{JobID: newJob.RequestID}
			dbContent, err := logPart.Collect(db)
			Expect(dbContent).To(Equal(&content))
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Describe("Get and Save", func() {

		It("returns an error if no db connection is given for get", func() {
			logPart := LogPart{"the_job_id", 1, "the log chunck", false, time.Now()}
			err := logPart.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if no db connection is given for save", func() {
			logPart := LogPart{"the_job_id", 1, "the log chunck", false, time.Now()}
			err := logPart.Save(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should not save a log part if the job with the same id does not exist", func() {
			// save chunck
			logPart := LogPart{uuid.New(), 1, "the log chunck", false, time.Now()}
			err := logPart.Save(db)
			Expect(err).To(HaveOccurred())
		})

		It("should save a log part and get the results", func() {
			// add a job related to the log chuncks
			newJob := Job{}
			newJob.ExecuteScriptExample()
			err := newJob.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// save chunck
			logPart := LogPart{newJob.RequestID, 1, "the log chunck", false, time.Now()}
			err = logPart.Save(db)
			Expect(err).NotTo(HaveOccurred())

			dbLogPart := LogPart{JobID: logPart.JobID, Number: logPart.Number}
			err = dbLogPart.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbLogPart.JobID).To(Equal(logPart.JobID))
			Expect(dbLogPart.Number).To(Equal(logPart.Number))
			Expect(dbLogPart.Content).To(Equal(logPart.Content))
			Expect(dbLogPart.Final).To(Equal(logPart.Final))
			Expect(dbLogPart.CreatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(logPart.CreatedAt.Format("2006-01-02 15:04:05.99")))
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
