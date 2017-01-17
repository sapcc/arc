// +build integration

package main

import (
	"database/sql"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pki"
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

		runRoutineTasks(db)

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

	It("should clean tokens older than 1 hour", func() {
		token := uuid.New()
		_, err := db.Exec(ownDb.InsertTokenWithCreatedAtQuery, token, "default", `{"CN":"blafasel"}`, time.Now().Add(-3700*time.Second))
		Expect(err).NotTo(HaveOccurred())

		runRoutineTasks(db)

		_, err = pki.GetTestToken(db, token)
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(sql.ErrNoRows))
	})

	It("should not clean tokens younger than 1 hour", func() {
		token := uuid.New()
		_, err := db.Exec(ownDb.InsertTokenWithCreatedAtQuery, token, "default", `{"CN":"blafasel"}`, time.Now().Add(-60*time.Second))
		Expect(err).NotTo(HaveOccurred())

		runRoutineTasks(db)

		res, err := pki.GetTestToken(db, token)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal("default"))
	})

	It("should clean certificates older than 2 years", func() {
		_, err := db.Exec(ownDb.InsertCertificateQuery,
			`certificateFingerprint(*x509Cert)`,
			`certSubject.CommonName`,
			`certSubject.Country`,
			`certSubject.Locality`,
			`certSubject.Organization`,
			`certSubject.OrganizationalUnit`,
			time.Now().Add(-17521*time.Hour),
			time.Now().Add(-1*time.Hour),
			`pemCert`,
		)
		Expect(err).NotTo(HaveOccurred())

		runRoutineTasks(db)

		var rows int
		err = db.QueryRow("SELECT COUNT(*) FROM certificates").Scan(&rows)
		Expect(err).NotTo(HaveOccurred())
		Expect(rows).To(Equal(0))
	})

	It("should not clean certificates younger than 2 years", func() {
		_, err := db.Exec(ownDb.InsertCertificateQuery,
			`certificateFingerprint(*x509Cert)`,
			`certSubject.CommonName`,
			`certSubject.Country`,
			`certSubject.Locality`,
			`certSubject.Organization`,
			`certSubject.OrganizationalUnit`,
			time.Now(),
			time.Now().Add(17520*time.Hour),
			`pemCert`,
		)
		Expect(err).NotTo(HaveOccurred())

		runRoutineTasks(db)

		var rows int
		err = db.QueryRow("SELECT COUNT(*) FROM certificates").Scan(&rows)
		Expect(err).NotTo(HaveOccurred())
		Expect(rows).To(Equal(1))
	})

})
