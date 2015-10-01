// +build integration

package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lock", func() {

	Describe("Get and Save", func() {

		It("returns an error if no db connection is given", func() {
			lock := Lock{}
			err := lock.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if no lock found", func() {
			lock := Lock{}
			lock.Example()
			err := lock.Get(db)
			Expect(err).To(HaveOccurred())
		})

		It("return a lock", func() {
			lock := Lock{}
			lock.Example()
			err := lock.Save(db)
			Expect(err).NotTo(HaveOccurred())

			dbLock := Lock{LockID: lock.LockID}
			err = dbLock.Get(db)
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
