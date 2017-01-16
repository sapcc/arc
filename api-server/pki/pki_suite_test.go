// +build integration

package pki_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pki"

	"database/sql"
	"testing"
)

var (
	db *sql.DB
)

func TestPki(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pki Suite")
}

var _ = BeforeSuite(func() {
	var err error
	env := os.Getenv("ARC_ENV")
	if env == "" {
		env = "test"
	}
	db, err = NewConnection("../db/dbconf.yml", env)
	Expect(err).NotTo(HaveOccurred())
	// set the pki configuration
	err = pki.SetupSigner("../test/ca.pem", "../test/ca-key.pem", "../etc/pki_default_config.json")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	db.Close()
})

var _ = BeforeEach(func() {
	DeleteAllRowsFromTable(db, "tokens")
	DeleteAllRowsFromTable(db, "certificates")
})
