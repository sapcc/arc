// +build integration

package main

import (
	"os"

	cfssl_config "github.com/cloudflare/cfssl/config"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pki"

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
	router = newRouter(env)

	// set the pki configuration
	pkiConfig.CAFile = pki.PathTo("test/ca.pem")
	pkiConfig.CAKeyFile = pki.PathTo("test/ca-key.pem")
	pkiConfig.CFG, err = cfssl_config.LoadFile(pki.PathTo("test/local.json"))
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
	DeleteAllRowsFromTable(db, "tokens")
	DeleteAllRowsFromTable(db, "certificates")
})
