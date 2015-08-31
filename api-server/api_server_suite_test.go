// +build integration

package main

import (
	"os"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"

	"testing"
)

var router *mux.Router

func TestApiServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApiServer Suite")
}

var _ = BeforeSuite(func() {
	var err error
	env := os.Getenv("ENV")
	if env == "" {
		env = "test"
	}
	db, err = NewConnection("db/dbconf.yml", env)
	Expect(err).NotTo(HaveOccurred())
	router = newRouter()
})

var _ = AfterSuite(func() {
	db.Close()
})

var _ = BeforeEach(func() {
	DeleteAllRowsFromTable(db, "jobs")
	DeleteAllRowsFromTable(db, "agents")
	DeleteAllRowsFromTable(db, "logs")
	DeleteAllRowsFromTable(db, "log_parts")
})
