package models_test

import (
	"time"
	"fmt"
	"strings"
	
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"	
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"code.google.com/p/go-uuid/uuid"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Log", func() {

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			newLog := Log{JobID:uuid.New()}
			err := newLog.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should return a log string from the log table if exists", func() {
			job_id := uuid.New()
			// add a job related to the log
			newJob := Job{Request: arc.Request{RequestID: job_id}}
			newJob.Save(db)
			
			content := "Log content"
			
			// insert a log
			_, err := db.Exec(InsertLogQuery, job_id, content, time.Now(), time.Now())
			Expect(err).NotTo(HaveOccurred())
			
			// get the log
			newLog := Log{JobID:job_id}
			err = newLog.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(Equal(newLog.Content))
		})

		It("should collect the log parts if a log from the log table doesn't exist", func() {
			job_id := uuid.New()
			// add a job related to the log
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
			
			// get the log
			newLog := Log{JobID:job_id}
			err := newLog.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(Equal(newLog.Content))
		})

	})
	
	Describe("ProcessLogReply", func() {
		
		It("returns an error if no db connection is given", func() {
			reply := arc.Reply{RequestID: uuid.New()}
			err := ProcessLogReply(nil, &reply)
			Expect(err).To(HaveOccurred())
		})

		It("should not save a log part entry if the payload is empty", func() {})
		It("should save a log part entry if the payload is not empty", func() {})
		It("should collect all log parts, save a new log entry and remove all log parts if final attribute is true ", func() {})		
	})	

})
