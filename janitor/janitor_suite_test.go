// +build integration

package janitor

import (
	"database/sql"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

func TestJanitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Janitor Suite")
}

var db *sql.DB

var _ = BeforeSuite(func() {
	var err error
	env := os.Getenv("ARC_ENV")
	if env == "" {
		env = "test"
	}
	// set the test database
	db, err = dbConnection("../api-server/db/dbconf.yml", env)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	db.Close()
})

var _ = BeforeEach(func() {
	ownDb.DeleteAllRowsFromTable(db, "jobs")
	ownDb.DeleteAllRowsFromTable(db, "logs")
	ownDb.DeleteAllRowsFromTable(db, "log_parts")
	ownDb.DeleteAllRowsFromTable(db, "locks")
	ownDb.DeleteAllRowsFromTable(db, "tokens")
	ownDb.DeleteAllRowsFromTable(db, "certificates")
})
