package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"	
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	//arc "gitHub.***REMOVED***/monsoon/arc/arc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Log", func() {

	BeforeEach(func() {
		DeleteAllRowsFromTable(db, "logs")
		DeleteAllRowsFromTable(db, "jobs")
	})

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			newLog := Log{}
			err := newLog.Get(nil, "super_id")
			Expect(err).To(HaveOccurred())
		})

	})

})
