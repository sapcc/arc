// +build integration

package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registry", func() {

	Describe("Get and Save", func() {

		It("returns an error if no db connection is given", func() {
			registry := Registry{}
			err := registry.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if no registry found", func() {
			registry := Registry{}
			registry.Example()
			err := registry.Get(db)
			Expect(err).To(HaveOccurred())
		})

		It("return a registry", func() {
			registry := Registry{}
			registry.Example()
			err := registry.Save(db)
			Expect(err).NotTo(HaveOccurred())

			dbRegistry := Registry{RegistryID: registry.RegistryID}
			err = dbRegistry.Get(db)
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
