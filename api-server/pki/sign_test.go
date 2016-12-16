// +build integration

package pki_test

import (
	"crypto/x509"
	"encoding/pem"

	. "gitHub.***REMOVED***/monsoon/arc/api-server/pki"

	"net/http"
	"os"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = Describe("Sign csr", func() {

	var (
		cfg cli.Config
	)

	JustBeforeEach(func() {
		var err error
		cfg.CAFile = PathTo("../test/ca.pem")
		cfg.CAKeyFile = PathTo("../test/ca-key.pem")
		cfg.CFG, err = config.LoadFile(PathTo("../test/local.json"))
		Expect(err).NotTo(HaveOccurred())
	})

	It("Signs a CSR", func() {
		token := CreateTestToken(db, `{}`)
		csr, err := os.Open(PathTo("../test/test.csr"))
		Expect(err).NotTo(HaveOccurred())
		defer csr.Close()

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, csr)
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, req, &cfg)
		Expect(err).NotTo(HaveOccurred())
		cert, _ := pem.Decode(*pemCert)
		Expect(cert.Type).To(Equal("CERTIFICATE"))
	})

	It("Requires a valid token", func() {
		token := uuid.New()
		csr, err := os.Open(PathTo("../test/test.csr"))
		Expect(err).NotTo(HaveOccurred())
		defer csr.Close()

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, csr)
		Expect(err).NotTo(HaveOccurred())

		pemCert, ca, err := SignToken(db, token, req, &cfg)
		Expect(err).To(HaveOccurred())
		_, ok := err.(SignForbidden)
		Expect(ok).To(Equal(true))
		var empty *[]byte
		Expect(pemCert).To(Equal(empty))
		Expect(ca).To(Equal(""))
	})

	It("Invalidates a token", func() {
		token := CreateTestToken(db, `{"CN":"blafasel"}`)
		csr, err := os.Open(PathTo("../test/test.csr"))
		Expect(err).NotTo(HaveOccurred())
		defer csr.Close()

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, csr)
		Expect(err).NotTo(HaveOccurred())

		_, _, err = SignToken(db, token, req, &cfg)
		Expect(err).NotTo(HaveOccurred())

		r, err := db.Query("SELECT id from tokens where id=$1", token)
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Next()).To(BeFalse())
	})

	It("Enforces the certificate subject", func() {
		token := CreateTestToken(db, `{"CN":"testcn", "names":[{"OU":"testou"}]}`)
		csr, err := os.Open(PathTo("../test/test.csr"))
		Expect(err).NotTo(HaveOccurred())
		defer csr.Close()

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, csr)
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, req, &cfg)
		Expect(err).NotTo(HaveOccurred())

		cert, _ := pem.Decode(*pemCert)
		x509Cert, err := x509.ParseCertificate(cert.Bytes)
		s := x509Cert.Subject
		Expect(err).NotTo(HaveOccurred())
		Expect(s.CommonName).To(Equal("testcn"))
		Expect(s.OrganizationalUnit[0]).To(Equal("testou"))
		Expect(len(s.Country)).To(BeZero())
		Expect(len(s.Organization)).To(BeZero())
	})

})
