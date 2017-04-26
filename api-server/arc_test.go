// +build integration

package main

import (
	"time"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/transport/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"reflect"
)

var _ = Describe("Arc", func() {

	JustBeforeEach(func() {
		config.Identity = "darwin"
		config.Transport = "fake"
		var err error
		tp, err = arcNewConnection(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should receive all registries and update facts", func(done Done) {
		// subscribe to all replies
		go arcSubscribeReplies(tp)
		//This is a drity hack, we need to do better
		time.Sleep(100 * time.Millisecond)

		// write to chan
		reg := models.Registration{}
		reg.Example()
		tp.Registration(&reg.Registration)

		// wait till done
		ftp, ok := tp.(*fake.FakeClient)
		Expect(ok).Should(BeTrue())
		<-ftp.Done

		// check registry is been saved
		dbAgent := models.Agent{AgentID: reg.Sender}
		err := dbAgent.Get(db)
		Expect(err).NotTo(HaveOccurred())

		// check facts
		checkFacts, err := models.JSONBfromString(reg.Payload)
		Expect(err).NotTo(HaveOccurred())
		eq := reflect.DeepEqual(dbAgent.Facts, *checkFacts)
		Expect(eq).To(Equal(true))
		close(done)
	}, 2.0)

	It("should receive all replies, update job and save log part", func(done Done) {
		// save a job
		job := models.Job{}
		job.RpcVersionExample()
		err := job.Save(db)
		Expect(err).NotTo(HaveOccurred())

		// subscribe to all replies
		go arcSubscribeReplies(tp)
		//This is a drity hack, we need to do better
		time.Sleep(100 * time.Millisecond)

		// write to chan
		reply := models.Reply{}
		reply.ExecuteScriptExample(job.RequestID, true, "Chunky chunk", 2)
		tp.Reply(&reply.Reply)

		// wait till done
		ftp, ok := tp.(*fake.FakeClient)
		Expect(ok).Should(BeTrue())
		<-ftp.Done

		// check job has been updated
		dbJob := models.Job{Request: arc.Request{RequestID: job.RequestID}}
		err = dbJob.Get(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbJob.Status).To(Equal(arc.Complete))
		Expect(dbJob.UpdatedAt.Format("2006-01-02 15:04:05.999")).NotTo(Equal(job.UpdatedAt.Format("2006-01-02 15:04:05.999")))

		// check the log part is saved
		dbLogPart := models.LogPart{JobID: job.RequestID, Number: 2}
		err = dbLogPart.Get(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbLogPart.Content).To(Equal(reply.Payload))
		close(done)
	}, 2.0)

})
