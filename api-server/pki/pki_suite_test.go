// +build integration

package pki_test

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"

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
})

var _ = AfterSuite(func() {
	db.Close()
})

var _ = BeforeEach(func() {
	DeleteAllRowsFromTable(db, "tokens")
})

func pathTo(p string) string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, "../", p)
}
