// +build integration

package models_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"

	auth "github.com/sapcc/arc/api-server/auth"
	. "github.com/sapcc/arc/api-server/models"
	"github.com/sapcc/arc/api-server/pagination"
	arc "github.com/sapcc/arc/arc"
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
			Expect(len(dbJobs)).To(Equal(3))
			// check that the jobs are sorted by created_at descending
			Expect(dbJobs[0].CreatedAt.Sub(dbJobs[1].CreatedAt)).To(BeNumerically(">", 0))
			Expect(dbJobs[1].CreatedAt.Sub(dbJobs[2].CreatedAt)).To(BeNumerically(">", 0))
		})

	})

	Describe("GetAuthorized", func() {

		var (
			authorization = auth.Authorization{}
			pagination    = pagination.Pagination{}
		)

		JustBeforeEach(func() {
			jobs := Jobs{}
			jobs.CreateAndSaveRpcVersionExamples(db, 3)
			authorization.IdentityStatus = "Confirmed"
			authorization.User = auth.User{Id: "123456789", Name: "Arturo", DomainId: "monsoon2_id", DomainName: "monsoon_name"}
			authorization.ProjectId = "test-project"
			pagination.Offset = 0
			pagination.Limit = 25
		})

		It("returns an error if no db connection is given", func() {
			jobs := Jobs{}
			err := jobs.GetAuthorized(nil, &authorization, "", &pagination)
			Expect(err).To(HaveOccurred())
		})

		It("should return all jobs with same project", func() {
			// add a new job with different project
			job := Job{}
			job.ExecuteScriptExample()
			job.Project = "miau"
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// change authorization
			authorization.ProjectId = "miau"

			dbJobs := Jobs{}
			err = dbJobs.GetAuthorized(db, &authorization, "", &pagination)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbJobs)).To(Equal(1))
		})

		It("should return a project authorization error", func() {
			authorization.ProjectId = "Some other project"

			dbJobs := Jobs{}
			err := dbJobs.GetAuthorized(db, &authorization, "", &pagination)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbJobs)).To(Equal(0))
		})

		It("should filter results per agent", func() {
			// add a new job with project and to attr
			job := Job{}
			job.ExecuteScriptExample()
			job.Project = "miau"
			job.To = "my_test_laptop"
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// add a new job just with to attr
			job2 := Job{}
			job2.ExecuteScriptExample()
			job2.Project = "miau"
			job2.To = "other_laptop"
			err = job2.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// change authorization
			authorization.ProjectId = "miau"

			dbJobs := Jobs{}
			err = dbJobs.GetAuthorized(db, &authorization, "", &pagination)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbJobs)).To(Equal(2))
			test1 := dbJobs[0].RequestID == job.RequestID || dbJobs[0].RequestID == job2.RequestID
			test2 := dbJobs[1].RequestID == job.RequestID || dbJobs[1].RequestID == job2.RequestID
			Expect(test1).To(Equal(true))
			Expect(test2).To(Equal(true))

			dbJobs = Jobs{}
			err = dbJobs.GetAuthorized(db, &authorization, "my_test_laptop", &pagination)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbJobs)).To(Equal(1))
			Expect(dbJobs[0].RequestID).To(Equal(job.RequestID))
		})

		It("should return the results ordered by update descendent", func() {
			// add a new job with a small update time
			job := Job{}
			job.ExecuteScriptExample()
			job.Project = "miau"
			job.UpdatedAt = time.Now().Add(-30 * time.Minute)
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// add a new job with a bigger update time
			job2 := Job{}
			job2.ExecuteScriptExample()
			job2.Project = "miau"
			job.UpdatedAt = time.Now().Add(-5 * time.Minute)
			err = job2.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// change authorization
			authorization.ProjectId = "miau"

			dbJobs := Jobs{}
			err = dbJobs.GetAuthorized(db, &authorization, "", &pagination)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbJobs)).To(Equal(2))
			Expect(dbJobs[0].RequestID).To(Equal(job2.RequestID))
			Expect(dbJobs[1].RequestID).To(Equal(job.RequestID))
		})

	})

})

var _ = Describe("Job", func() {

	Describe("CreateJob", func() {

		var (
			user  = auth.User{Id: "123456789", Name: "Arturo", DomainId: "monsoon2_id", DomainName: "monsoon_name"}
			agent = Agent{}
		)

		JustBeforeEach(func() {
			agent.Example()
			agent.AgentID = "darwin"
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error data is no json conform", func() {
			noValidJson := `"to":"darwin"`
			strSlice := []byte(noValidJson)
			job, err := CreateJob(db, &strSlice, uuid.New(), &user)
			Expect(err).To(HaveOccurred())
			var newJob *Job
			Expect(job).To(Equal(newJob))
		})

		It("returns error data is not valid", func() {
			noValidData := `{"to":"darwin","timeout":60,"agent":"execute","payload":"echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""}` // action is missing
			strSlice := []byte(noValidData)
			job, err := CreateJob(db, &strSlice, uuid.New(), &user)
			Expect(err).To(HaveOccurred())
			var newJob *Job
			Expect(job).To(Equal(newJob))
		})

		It("returns error user id is blank", func() {
			data := `{"to":"darwin","timeout":60,"agent":"execute","action":"script","payload":"echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""}`
			strSlice := []byte(data)
			job, err := CreateJob(db, &strSlice, uuid.New(), &auth.User{})
			Expect(err).To(HaveOccurred())
			var newJob *Job
			Expect(job).To(Equal(newJob))
		})

		It("should create a job", func() {
			to := "darwin"
			timeout := 60
			arcAgent := "execute"
			action := "script"
			payload := `"payload":"echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""`
			noValidData := fmt.Sprintf(`{"to":%q,"timeout":%v,"agent":%q,"action":%q,"payload":%q}`, to, timeout, arcAgent, action, payload)
			strSlice := []byte(noValidData)
			job, err := CreateJob(db, &strSlice, uuid.New(), &user)
			Expect(err).NotTo(HaveOccurred())
			Expect(job.To).To(Equal(to))
			Expect(job.Timeout).To(Equal(timeout))
			Expect(job.Agent).To(Equal(arcAgent))
			Expect(job.Action).To(Equal(action))
			Expect(job.Payload).To(Equal(payload))
			// should create a job with the project id from the target agent
			Expect(job.Project).To(Equal(agent.Project))

			// should save also the given user
			eq, err := CompareUserWithJobUser(user, job.User)
			Expect(err).NotTo(HaveOccurred())
			Expect(eq).To(Equal(true))
		})

	})

	Describe("CreateJobAuthorized", func() {

		var (
			user          = auth.User{Id: "123456789", Name: "Arturo", DomainId: "monsoon2_id", DomainName: "monsoon_name"}
			authorization = auth.Authorization{}
			agent         = Agent{}
		)

		JustBeforeEach(func() {
			// authorization
			authorization.IdentityStatus = "Confirmed"
			authorization.User = user
			authorization.ProjectId = "test-project"
			// agent
			agent.Example()
			agent.AgentID = "darwin"
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should save the user id from the token", func() {
			data := fmt.Sprintf(`{"to":%q,"timeout":60,"agent":"execute","action":"script","payload":"echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""}`, agent.AgentID)
			strSlice := []byte(data)
			job, err := CreateJobAuthorized(db, &strSlice, uuid.New(), &authorization)
			Expect(err).NotTo(HaveOccurred())

			// should save also the given user
			eq, err := CompareUserWithJobUser(user, job.User)
			Expect(err).NotTo(HaveOccurred())
			Expect(eq).To(Equal(true))
		})

		It("should return a project authorization error", func() {
			authorization.ProjectId = "Some other project"

			// create job
			data := fmt.Sprintf(`{"to":%q,"timeout":60,"agent":"execute","action":"script","payload":"echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""}`, agent.AgentID)
			strSlice := []byte(data)
			job, err := CreateJobAuthorized(db, &strSlice, uuid.New(), &authorization)
			Expect(err).To(HaveOccurred())
			_, ok := err.(auth.NotAuthorized)
			Expect(ok).To(Equal(true))
			var newJob *Job
			Expect(job).To(Equal(newJob))
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

		It("should update the status and update at", func() {
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

	Describe("FailQueuedJobs", func() {

		It("returns an error if no db connection is given", func() {
			_, err := FailQueuedJobs(nil)
			Expect(err).To(HaveOccurred())
		})

		It("fail jobs which no heartbeat was send back after created_at + 60 sec", func() {
			// save a job
			job := Job{}
			job.CustomExecuteScriptExample(arc.Queued, time.Now().Add(-61*time.Second), 120)
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// clean jobs
			occurHeartbeat, err := FailQueuedJobs(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(occurHeartbeat).To(Equal(int64(1)))

			// check job
			dbJob := Job{Request: arc.Request{RequestID: job.RequestID}}
			err = dbJob.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJob.Status).To(Equal(arc.Failed))
		})

	})

	Describe("FailJobsTimeoutQuery", func() {

		It("returns an error if no db connection is given", func() {
			_, err := FailExpiredJobs(nil)
			Expect(err).To(HaveOccurred())
		})

		It("fail jobs which the timeout + 60 sec has exceeded and still in queued or executing status", func() {
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
			occurTimeOut, err := FailExpiredJobs(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(occurTimeOut).To(Equal(int64(2)))

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

	Describe("PruneJobs", func() {

		It("returns an error if no db connection is given", func() {
			_, err := PruneJobs(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should clean jobs are older than 30 days", func() {
			// save a job older than 30 days
			job := Job{}
			job.CustomExecuteScriptExample(arc.Complete, time.Now().Add((-24*31)*time.Hour), 15) // 60 sec extra to be sure
			err := job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// save a job not older than 30 days
			job.CustomExecuteScriptExample(arc.Complete, time.Now().Add((-24*15)*time.Hour), 15) // 60 sec extra to be sure
			err = job.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// clean jobs
			occurOld, err := PruneJobs(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(occurOld).To(Equal(int64(1)))
		})

	})

})
