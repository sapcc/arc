package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"	

	"database/sql"
	"testing"
)

var (
	db *sql.DB
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	var err error
	db, err = NewConnection("../db/dbconf.yml", "test")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	db.Close()
})

var _ = BeforeEach(func() {
	DeleteAllRowsFromTable(db, "jobs")
	DeleteAllRowsFromTable(db, "facts")
	DeleteAllRowsFromTable(db, "logs")
	DeleteAllRowsFromTable(db, "log_parts")
})