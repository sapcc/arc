package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/pem"
	"os"
	"path"

	"github.com/pborman/uuid"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

// PathTo generates a path to the pki-package
func PathTo(p string) string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, p)
}

// CreateTestToken save a test token in the db
func CreateTestToken(db *sql.DB, subject string) string {
	token := uuid.New()
	_, err := db.Exec(ownDb.InsertTokenQuery, token, "default", subject)
	if err != nil {
		panic(err)
	}
	return token
}

// GetTestToken get a saved token
func GetTestToken(db *sql.DB, token string) (string, error) {
	var profile string
	var subjectData []byte
	err := db.QueryRow(ownDb.GetTokenQuery, token).Scan(&profile, &subjectData)
	if err != nil {
		return "", err
	}
	return profile, nil
}

// CreateCsr creates a csr
func CreateCsr(CommonName, Organization, OrganizationalUnit string) (*bytes.Buffer, error) {
	// generate key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// create csr template
	csrTemplate := x509.CertificateRequest{
		SignatureAlgorithm: x509.SHA256WithRSA,
		Subject: pkix.Name{
			CommonName:         CommonName,
			Organization:       []string{Organization},
			OrganizationalUnit: []string{OrganizationalUnit},
		},
	}

	csrData, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, key)
	if err != nil {
		return nil, err
	}

	var csr bytes.Buffer
	err = pem.Encode(&csr, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrData})
	if err != nil {
		return nil, err
	}

	return &csr, nil
}
