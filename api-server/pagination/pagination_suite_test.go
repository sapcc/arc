// +build integration

package pagination_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sapcc/arc/api-server/db"

	"database/sql"
	"testing"
)

var (
	db *sql.DB
)

func TestPagination(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pagination Suite")
}

var _ = BeforeSuite(func() {
	var err error
	env := os.Getenv("ARC_ENV")
	if env == "" {
		env = "test"
	}
	db, err = NewConnection("../db/dbconf.yml", env)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	db.Close()
})

var _ = BeforeEach(func() {
	DeleteAllRowsFromTable(db, "jobs")
	DeleteAllRowsFromTable(db, "agents")
	DeleteAllRowsFromTable(db, "logs")
	DeleteAllRowsFromTable(db, "log_parts")
	DeleteAllRowsFromTable(db, "locks")
	DeleteAllRowsFromTable(db, "tags")
})
