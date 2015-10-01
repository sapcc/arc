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

	Describe("isConcurrencySafe", func() {

		It("return false when the key already exist in the db", func() {
			lock := Lock{}
			lock.Example()
			err := lock.Save(db)
			Expect(err).NotTo(HaveOccurred())

			safe, err := IsConcurrencySafe(db, lock.LockID, lock.AgentID)
			Expect(safe).To(Equal(false))
			Expect(err).NotTo(HaveOccurred())
		})

		It("return true if the key does not exist", func() {
			safe, err := IsConcurrencySafe(db, "miau", "darwin")
			Expect(safe).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())
		})

		It("return true when another error is thrown", func() {
			safe, err := IsConcurrencySafe(nil, "test", "darwin")
			Expect(safe).To(Equal(false))
			Expect(err).To(HaveOccurred())
		})

	})

})
