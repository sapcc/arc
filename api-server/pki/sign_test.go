// +build integration

package pki_test

import (
	"encoding/pem"

	. "gitHub.***REMOVED***/monsoon/arc/api-server/pki"

	"net/http"
	"os"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/go-uuid/uuid"
)

var _ = Describe("Sign csr", func() {

	var (
		cfg cli.Config
	)

	It("Signs a CSR", func() {
		token := createToken(`{}`)
		csr, err := os.Open(pathTo("test/test.csr"))
		Expect(err).NotTo(HaveOccurred())
		defer csr.Close()

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, csr)
		Expect(err).NotTo(HaveOccurred())

		cfg.CAFile = pathTo("test/ca.pem")
		cfg.CAKeyFile = pathTo("test/ca-key.pem")
		cfg.CFG, err = config.LoadFile(pathTo("test/local.json"))
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, req, &cfg)
		Expect(err).NotTo(HaveOccurred())
		cert, _ := pem.Decode(*pemCert)
		Expect(cert.Type).To(Equal("CERTIFICATE"))
	})

})

func createToken(subject string) string {
	token := uuid.New()
	_, err := db.Exec("INSERT INTO tokens (id, profile, subject) VALUES($1, $2, $3)", token, "default", subject)
	if err != nil {
		panic(err)
	}
	return token
}
