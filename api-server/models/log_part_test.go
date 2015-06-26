package models_test

import (
	"fmt"
	"strings"
	"time"
	"code.google.com/p/go-uuid/uuid"

	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			job_id := uuid.New()
			// add a job related to the log chuncks
			newJob := Job{Request: arc.Request{RequestID: job_id}}
			newJob.Save(db)
			// save different chuncks
			var contentSlice [3]string
			for i := 0; i < 3; i++ {
				chunck := fmt.Sprintf("This is the %d chunck", i)
				contentSlice[i] = chunck
				logPart := LogPart{job_id, uint(i), chunck, false, time.Now()}
				err := logPart.Save(db)
				if err != nil {
					fmt.Println(err)
				}
			}
			content := strings.Join(contentSlice[:], "")

			// collect
			logPart := LogPart{JobID: job_id}
			dbContent, err := logPart.Collect(db)
			Expect(dbContent).To(Equal(&content))
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Describe("Save", func() {

		It("returns an error if no db connection is given", func() {
			logPart := LogPart{"the_job_id", 1, "the log chunck", false, time.Now()}
			err := logPart.Save(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should not save a log part if the job with the same id does not exist", func() {
			job_id := uuid.New()
			
			// save chunck
			logPart := LogPart{job_id, 1, "the log chunck", false, time.Now()}
			err := logPart.Save(db)
			Expect(err).To(HaveOccurred())
		})

		It("should save a log part", func() {
			job_id := uuid.New()
			// add a job related to the log chuncks
			newJob := Job{Request: arc.Request{RequestID: job_id}}
			newJob.Save(db)

			// save chunck
			logPart := LogPart{job_id, 1, "the log chunck", false, time.Now()}
			err := logPart.Save(db)

			dbLogPart := LogPart{}
			db.QueryRow(GetLogPartQuery, job_id).Scan(&dbLogPart.JobID, &dbLogPart.Number, &dbLogPart.Content, &dbLogPart.Final, &dbLogPart.CreatedAt)
			Expect(dbLogPart.JobID).To(Equal(logPart.JobID))
			Expect(dbLogPart.Number).To(Equal(logPart.Number))
			Expect(dbLogPart.Content).To(Equal(logPart.Content))
			Expect(dbLogPart.Final).To(Equal(logPart.Final))
			Expect(dbLogPart.CreatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(logPart.CreatedAt.Format("2006-01-02 15:04:05.99")))
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
