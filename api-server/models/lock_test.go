// +build integration

package models_test

import (
	. "github.com/sapcc/arc/api-server/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"
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

	Describe("PruneLocks", func() {

		It("should clean logs older than 6 min", func() {
			// save a locks 6 min old
			lock := Lock{}
			lock.Example()
			lock.CreatedAt = time.Now().Add(-360 * time.Second) // 6 min
			err := lock.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// save a lock just now
			lock2 := Lock{}
			lock2.Example()
			err = lock2.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// clean locks
			affectedLocks, err := PruneLocks(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(affectedLocks).To(Equal(int64(1)))

			// check the right one is still in the db
			dbLock := Lock{LockID: lock2.LockID}
			err = dbLock.Get(db)
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
