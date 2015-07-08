// +build integration

package main

import (
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/transport/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Arc", func() {

	JustBeforeEach(func() {
		config.Identity = "darwin"
		config.Transport = "fake"
		var err error
		tp, err = arcNewConnection(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should receive all registries and update facts", func() {
		// subscribe to all replies
		go arcSubscribeReplies(tp)

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
		Expect(dbAgent.Facts).To(Equal(reg.Payload))
	})

	It("should receive all replies, update job and save log", func() {
		// save a job
		job := models.Job{}
		job.RpcVersionExample()
		err := job.Save(db)
		Expect(err).NotTo(HaveOccurred())

		// subscribe to all replies
		go arcSubscribeReplies(tp)

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

		// check the log is been saved
		dbLog := models.Log{JobID: job.RequestID}
		err = dbLog.Get(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbLog.Content).To(Equal(reply.Payload))
	})

})
