// +build integration

package pki_test

import (
	"time"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/pki"

	"crypto/x509"
	"encoding/pem"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = Describe("Sign csr", func() {

	It("Signs a CSR", func() {
		token := CreateTestToken(db, `{}`)
		csr, _, err := CreateSignReqCertAndPrivKey("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, csr)
		Expect(err).NotTo(HaveOccurred())
		cert, _ := pem.Decode(*pemCert)
		Expect(cert.Type).To(Equal("CERTIFICATE"))
	})

	It("Requires a valid token", func() {
		token := uuid.New()
		csr, _, err := CreateSignReqCertAndPrivKey("testCsrCN", "test O", "test OU")
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
		csr, _, err := CreateSignReqCertAndPrivKey("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		_, _, err = SignToken(db, token, csr)
		Expect(err).NotTo(HaveOccurred())

		r, err := db.Query("SELECT id from tokens where id=$1", token)
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Next()).To(BeFalse())
	})

	It("Allows CN from CSR if not set in the tokens subject", func() {
		token := CreateTestToken(db, `{"names":[{"OU":"testou"}]}`)
		csr, _, err := CreateSignReqCertAndPrivKey("testCsrCN", "test O", "test OU")
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
		csr, _, err := CreateSignReqCertAndPrivKey("testCsrCN; DROP TABLE", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		_, _, err = SignToken(db, token, csr)
		Expect(err).To(HaveOccurred())
	})

	It("Enforces CN, O and OU if set in the tokens subject", func() {
		token := CreateTestToken(db, `{"CN":"enforced_CN", "names":[{"O": "enforced_O", "OU":"enforced_OU"}]}`)
		csr, _, err := CreateSignReqCertAndPrivKey("testCsrCN", "test O", "test OU")
		Expect(err).NotTo(HaveOccurred())

		pemCert, _, err := SignToken(db, token, csr)
		Expect(err).NotTo(HaveOccurred())

		cert, _ := pem.Decode(*pemCert)
		x509Cert, err := x509.ParseCertificate(cert.Bytes)
		s := x509Cert.Subject
		Expect(err).NotTo(HaveOccurred())
		Expect(s.CommonName).To(Equal("enforced_CN"))
		Expect(s.Organization[0]).To(Equal("enforced_O"))
		Expect(s.OrganizationalUnit[0]).To(Equal("enforced_OU"))
	})

})

var _ = Describe("PruneCertificates", func() {

	It("returns an error if no db connection is given", func() {
		occurrencies, err := PruneCertificates(nil)
		Expect(err).To(HaveOccurred())
		Expect(occurrencies).To(Equal(int64(0)))
	})

	It("should clean expired certificates", func() {
		// add expired certificate (15 min)
		_, err := db.Exec(ownDb.InsertCertificateQuery,
			uuid.New(),
			"common name",
			"country",
			"Locality",
			"Organization",
			"OrganizationalUnit",
			time.Now().Add((-60)*time.Minute),
			time.Now().Add((-15)*time.Minute),
			"pemCert",
		)
		Expect(err).NotTo(HaveOccurred())

		occurrencies, err := PruneCertificates(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(occurrencies).To(Equal(int64(1)))
	})

	It("should NOT clean not expired certificates ", func() {
		// add not expired certificate (still valid 15 min)
		_, err := db.Exec(ownDb.InsertCertificateQuery,
			uuid.New(),
			"common name",
			"country",
			"Locality",
			"Organization",
			"OrganizationalUnit",
			time.Now().Add((-60)*time.Minute),
			time.Now().Add((+15)*time.Minute),
			"pemCert",
		)
		Expect(err).NotTo(HaveOccurred())

		occurrencies, err := PruneCertificates(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(occurrencies).To(Equal(int64(0)))
	})

})
