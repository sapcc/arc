package pki

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/pem"
	"errors"
	"net/http"

	"github.com/cloudflare/cfssl/csr"
	"github.com/pborman/uuid"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

var (
	TLS_CERTIFICATE_MISSING = "TLS client authentication required for certificate renewal."
	TLS_CRT_SUBJECT_MISSING = "TLS client authentication requires certificat with CommonName, OrganizationalUnit and Organization."
	TLS_CRT_REQUEST_INVALID = "Certificate request invalid. "
)

type Subject struct {
	CommonName         string
	OrganizationalUnit string
	Organization       string
}

func TLSRequestSubject(r *http.Request) (Subject, error) {
	// check if certificate is given
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		return Subject{}, errors.New(TLS_CERTIFICATE_MISSING)
	}

	// Extract CN, OU and O from the certificate used on the tls communication
	cn := r.TLS.PeerCertificates[0].Subject.CommonName
	ou := ""
	if len(r.TLS.PeerCertificates[0].Subject.OrganizationalUnit) > 0 {
		ou = r.TLS.PeerCertificates[0].Subject.OrganizationalUnit[0]
	}
	o := ""
	if len(r.TLS.PeerCertificates[0].Subject.Organization) > 0 {
		o = r.TLS.PeerCertificates[0].Subject.Organization[0]
	}

	if len(cn) == 0 || len(ou) == 0 || len(o) == 0 {
		return Subject{}, errors.New(TLS_CRT_SUBJECT_MISSING)
	}

	return Subject{
		CommonName:         cn,
		OrganizationalUnit: ou,
		Organization:       o,
	}, nil
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

// CreateSignReqCert
// Creates a signing request cert PEM Block from a private key
func CreateSignReqCert(commonName, organization, organizationalUnit string, privKey interface{}) (csreq []byte, err error) {
	// create csr request template
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{Organization: []string{organization}, OrganizationalUnit: []string{organizationalUnit}, CommonName: commonName},
	}
	csrData, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privKey)
	if err != nil {
		return nil, err
	}
	block := pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrData,
	}
	csrData = pem.EncodeToMemory(&block)
	return csrData, nil
}

// CreateSignReqCertAndPrivKey
// Creates a signing request cert and private key
// SignatureAlgorithm is x509.ECDSAWithSHA256
func CreateSignReqCertAndPrivKey(commonName, organization, organizationalUnit string) (csreq, key []byte, err error) {
	req := csr.New()
	req.CN = commonName
	req.Names = []csr.Name{csr.Name{O: organization, OU: organizationalUnit}}
	csreq, key, err = csr.ParseRequest(req)
	return
}
