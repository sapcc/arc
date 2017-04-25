// +build integration

package models_test

import (
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"

	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"
)

var _ = Describe("Log", func() {

	Describe("Get and save", func() {

		It("returns an error if no db connection is given for get", func() {
			newLog := Log{JobID: uuid.New()}
			err := newLog.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if no db connection is given for save", func() {
			newLog := Log{JobID: uuid.New()}
			err := newLog.Save(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should save and get a log", func() {
			// add a job related to the log
			newJob := Job{}
			newJob.ExecuteScriptExample()
			err := newJob.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// insert a log
			content := "Log content"
			created_at := time.Now().Add((-5) * time.Minute)
			updated_at := time.Now().Add((-5) * time.Minute)
			log := Log{JobID: newJob.RequestID, Content: content, CreatedAt: created_at, UpdatedAt: updated_at}
			err = log.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// get the log
			newLog := Log{JobID: newJob.RequestID}
			err = newLog.GetOrCollect(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(newLog.Content).To(Equal(content))
			Expect(newLog.CreatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(created_at.Format("2006-01-02 15:04:05.99")))
			Expect(newLog.UpdatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(updated_at.Format("2006-01-02 15:04:05.99")))
		})

	})

	Describe("GetOrCollect", func() {

		It("returns an error if no db connection is given", func() {
			newLog := Log{JobID: uuid.New()}
			err := newLog.GetOrCollect(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should return a log string from the log table if exists", func() {
			// add a job related to the log
			newJob := Job{}
			newJob.ExecuteScriptExample()
			newJob.Save(db)

			// insert a log
			content := "Log content"
			log := Log{JobID: newJob.RequestID, Content: content, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			err := log.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// get the log
			newLog := Log{JobID: newJob.RequestID}
			err = newLog.GetOrCollect(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(Equal(newLog.Content))
		})

		It("should collect the log parts if a log from the log table doesn't exist", func() {
			// add a job related to the log
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
				if err != nil {
					fmt.Println(err)
				}
			}
			content := strings.Join(contentSlice[:], "")

			// get the log
			newLog := Log{JobID: newJob.RequestID}
			err := newLog.GetOrCollect(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(Equal(newLog.Content))
		})

	})

	Describe("ProcessLogReply", func() {

		It("returns an error if no db connection is given", func() {
			reply := arc.Reply{RequestID: uuid.New()}
			err := ProcessLogReply(nil, &reply, "darwin", true)
			Expect(err).To(HaveOccurred())
		})

		It("should not save a log part entry if the payload is empty", func() {
			// add a job related to the log
			newJob := Job{}
			newJob.ExecuteScriptExample()
			err := newJob.Save(db)
			Expect(err).NotTo(HaveOccurred())

			reply := arc.Reply{RequestID: newJob.RequestID, Number: 0, Payload: "", Final: false}
			err = ProcessLogReply(db, &reply, "darwin", true)
			Expect(err).NotTo(HaveOccurred())

			// check log
			newLog := Log{JobID: newJob.RequestID}
			err = newLog.Get(db)
			Expect(err).To(HaveOccurred())

			// check log parts
			logPart := LogPart{JobID: newJob.RequestID}
			_, err = logPart.Collect(db)
			Expect(err).To(HaveOccurred())
		})

		It("should save a log part entry if the payload is not empty", func() {
			chunck := "This is a chunck log"

			// add a job related to the log
			newJob := Job{}
			newJob.ExecuteScriptExample()
			err := newJob.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// process reply
			reply := arc.Reply{RequestID: newJob.RequestID, Payload: chunck, Final: false}
			err = ProcessLogReply(db, &reply, "darwin", true)
			Expect(err).NotTo(HaveOccurred())

			// check log
			newLog := Log{JobID: newJob.RequestID}
			err = newLog.Get(db)
			Expect(err).To(HaveOccurred())

			// check log parts
			logPart := LogPart{JobID: newJob.RequestID}
			dbContent, err := logPart.Collect(db)
			Expect(dbContent).To(Equal(&chunck))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should get an error if no job with the same id exists", func() {
			job_id := uuid.New()
			chunck := "This is a chunck log"

			// process reply
			reply := arc.Reply{RequestID: job_id, Payload: chunck}
			err := ProcessLogReply(db, &reply, "darwin", true)
			Expect(err).To(HaveOccurred())

			// check log parts
			logPart := LogPart{JobID: job_id}
			_, err = logPart.Collect(db)
			Expect(err).To(HaveOccurred())

			// check log
			newLog := Log{JobID: job_id}
			err = newLog.Get(db)
			Expect(err).To(HaveOccurred())
		})

		It("should check concurrency safe", func() {
			chunck := "This is a chunck log 1"
			chunck2 := "This is a chunck log 2"

			// add a job related to the log
			newJob := Job{}
			newJob.ExecuteScriptExample()
			err := newJob.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// process reply
			reply := arc.Reply{RequestID: newJob.RequestID, Number: 0, Payload: chunck, Final: false}
			err = ProcessLogReply(db, &reply, "darwin", true)
			Expect(err).NotTo(HaveOccurred())

			// process new reply
			newReply := arc.Reply{RequestID: newJob.RequestID, Number: 0, Payload: chunck2, Final: false}
			err = ProcessLogReply(db, &newReply, "darwin", true)
			_, ok := err.(ReplyExistsError)
			Expect(ok).To(Equal(true))
		})

		It("should not check concurrency safe", func() {
			chunck := "This is a chunck log 1"
			chunck2 := "This is a chunck log 2"

			// add a job related to the log
			newJob := Job{}
			newJob.ExecuteScriptExample()
			err := newJob.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// process reply
			reply := arc.Reply{RequestID: newJob.RequestID, Number: 0, Payload: chunck, Final: false}
			err = ProcessLogReply(db, &reply, "darwin", false)
			Expect(err).NotTo(HaveOccurred())

			// process new reply
			newReply := arc.Reply{RequestID: newJob.RequestID, Number: 0, Payload: chunck2, Final: false}
			err = ProcessLogReply(db, &newReply, "darwin", false)
			Expect(err).To(HaveOccurred())
			_, ok := err.(ReplyExistsError)
			Expect(ok).NotTo(Equal(true))
		})

	})

	Describe("AggregateLogs", func() {

		It("returns an error if no db connection is given", func() {
			occurrencies, err := AggregateLogs(nil)
			Expect(err).To(HaveOccurred())
			Expect(occurrencies).To(Equal(0))
		})

		It("should clean log parts with final state which are longer then 10 min", func() {
			// add a job related to the log chuncks
			job := Job{}
			job.ExecuteScriptExample()
			job.Save(db)

			// save different chuncks
			var contentSlice [3]string
			for i := 0; i < 3; i++ {
				chunck := fmt.Sprintf("This is the %d chunck", i)
				contentSlice[i] = chunck
				logPart := LogPart{job.RequestID, uint(i), chunck, false, time.Now().Add(-601 * time.Second)} // bit more than 10 min
				if i == 2 {
					logPart.Final = true
				}
				err := logPart.Save(db)
				Expect(err).NotTo(HaveOccurred())
			}
			content := strings.Join(contentSlice[:], "")

			// clean log parts
			occurrencies, err := AggregateLogs(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(occurrencies).To(Equal(1))

			// check log parts
			logPart := LogPart{JobID: job.RequestID}
			_, err = logPart.Collect(db)
			Expect(err).To(HaveOccurred())

			// check log
			newLog := Log{JobID: job.RequestID}
			err = newLog.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(newLog.Content).To(Equal(content))
		})

	})

	Describe("truncate payload log", func() {

		It("should truncate if text length is greater than 100 and not trhough an error", func() {
			newJob := Job{}
			newJob.ExecuteScriptExample()
			newJob.Payload = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec qu"
			err := newJob.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
