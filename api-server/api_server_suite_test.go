// +build integration

package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sapcc/arc/api-server/db"
	"github.com/sapcc/arc/api-server/pki"

	"testing"
)

var router *mux.Router

func TestApiServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApiServer Suite")
}

var _ = BeforeSuite(func() {
	var err error
	env = os.Getenv("ARC_ENV")
	if env == "" {
		env = "test"
	}
	// set the test database
	db, err = NewConnection("db/dbconf.yml", env)
	Expect(err).NotTo(HaveOccurred())

	// set the pki configuration
	err = pki.SetupSigner("test/ca.pem", "test/ca-key.pem", "etc/pki.json")
	pkiEnabled = true
	Expect(err).NotTo(HaveOccurred())

	// adding policies before createing router
	for _, pol := range policies {
		if polErr := warden.Manager.Create(pol); err != nil {
			log.Fatal(polErr)
		}
	}

	router = newRouter(env)
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
	DeleteAllRowsFromTable(db, "tokens")
	DeleteAllRowsFromTable(db, "certificates")
})
