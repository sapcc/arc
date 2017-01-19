// +build integration

package pki_test

import (
	"crypto/x509"
	"encoding/pem"

	. "gitHub.***REMOVED***/monsoon/arc/api-server/pki"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = Describe("Sign csr", func() {

	It("Signs a CSR", func() {
		token := CreateTestToken(db, `{}`)
		csr, _, err := CreateCSR("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, csr)
		Expect(err).NotTo(HaveOccurred())
		cert, _ := pem.Decode(*pemCert)
		Expect(cert.Type).To(Equal("CERTIFICATE"))
	})

	It("Requires a valid token", func() {
		token := uuid.New()
		csr, _, err := CreateCSR("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		pemCert, ca, err := SignToken(db, token, csr)
		Expect(err).To(HaveOccurred())
		_, ok := err.(SignForbidden)
		Expect(ok).To(Equal(true))
		var empty *[]byte
		Expect(pemCert).To(Equal(empty))
		Expect(ca).To(Equal(""))
	})

	It("Invalidates a token", func() {
		token := CreateTestToken(db, `{"CN":"blafasel"}`)
		csr, _, err := CreateCSR("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		_, _, err = SignToken(db, token, csr)
		Expect(err).NotTo(HaveOccurred())

		r, err := db.Query("SELECT id from tokens where id=$1", token)
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Next()).To(BeFalse())
	})

	It("Allows CN from CSR if not set in the tokens subject", func() {
		token := CreateTestToken(db, `{"names":[{"OU":"testou"}]}`)
		csr, _, err := CreateCSR("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, csr)
		Expect(err).NotTo(HaveOccurred())

		cert, _ := pem.Decode(*pemCert)
		x509Cert, err := x509.ParseCertificate(cert.Bytes)
		s := x509Cert.Subject
		Expect(err).NotTo(HaveOccurred())
		Expect(s.CommonName).To(Equal("testCsrCN"))
	})

	It("Refuses to Sign CSRs that contain a CN with funny characters", func() {
		token := CreateTestToken(db, `{"names":[{"OU":"testou"}]}`)
		csr, _, err := CreateCSR("testCsrCN; DROP TABLE", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		_, _, err = SignToken(db, token, csr)
		Expect(err).To(HaveOccurred())
	})

	It("Enforces CN, O and OU if set in the tokens subject", func() {
		token := CreateTestToken(db, `{"names":[{"O": "enforced O", "OU":"enforced OU"}]}`)
		csr, _, err := CreateCSR("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, csr)
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
