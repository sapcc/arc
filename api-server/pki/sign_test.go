// +build integration

package pki_test

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"

	. "gitHub.***REMOVED***/monsoon/arc/api-server/pki"

	"net/http"

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
		cfg.CFG, err = config.LoadFile(PathTo("../etc/pki_default_config.json"))
		Expect(err).NotTo(HaveOccurred())
	})

	It("Signs a CSR", func() {
		token := CreateTestToken(db, `{}`)
		csr, err := CreateCsr("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, csr)
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, req, &cfg)
		Expect(err).NotTo(HaveOccurred())
		cert, _ := pem.Decode(*pemCert)
		Expect(cert.Type).To(Equal("CERTIFICATE"))
	})

	It("Requires a valid token", func() {
		token := uuid.New()
		csr, err := CreateCsr("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

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
		csr, err := CreateCsr("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, csr)
		Expect(err).NotTo(HaveOccurred())

		_, _, err = SignToken(db, token, req, &cfg)
		Expect(err).NotTo(HaveOccurred())

		r, err := db.Query("SELECT id from tokens where id=$1", token)
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Next()).To(BeFalse())
	})

	It("Allows CN from CSR if not set in the tokens subject", func() {
		token := CreateTestToken(db, `{"names":[{"OU":"testou"}]}`)
		csr, err := CreateCsr("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, bytes.NewReader(csr.Bytes()))
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, req, &cfg)
		Expect(err).NotTo(HaveOccurred())

		cert, _ := pem.Decode(*pemCert)
		x509Cert, err := x509.ParseCertificate(cert.Bytes)
		s := x509Cert.Subject
		Expect(err).NotTo(HaveOccurred())
		Expect(s.CommonName).To(Equal("testCsrCN"))
	})

	It("nforces CN, O and OU if set in the tokens subject", func() {
		token := CreateTestToken(db, `{"names":[{"O": "enforced O", "OU":"enforced OU"}]}`)
		csr, err := CreateCsr("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		req, err := http.NewRequest("POST", "/api/v1/pki/sign/"+token, bytes.NewReader(csr.Bytes()))
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, req, &cfg)
		Expect(err).NotTo(HaveOccurred())

		cert, _ := pem.Decode(*pemCert)
		x509Cert, err := x509.ParseCertificate(cert.Bytes)
		s := x509Cert.Subject
		Expect(err).NotTo(HaveOccurred())
		Expect(s.CommonName).To(Equal("testCsrCN"))
		Expect(s.Organization[0]).To(Equal("enforced O"))
		Expect(s.OrganizationalUnit[0]).To(Equal("enforced OU"))
	})

})
