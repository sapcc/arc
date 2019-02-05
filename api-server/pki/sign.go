package pki

import (
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/info"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/universal"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

var (
	certSigner signer.Signer
)

// SignForbidden should be used to return a 403
type SignForbidden struct {
	Msg string
}

func (e SignForbidden) Error() string {
	return e.Msg
}

// SetupSigner initializes the Signer
func SetupSigner(caCertFile, caKeyFile, configFile string) (err error) {

	if _, err := os.Stat(caCertFile); err != nil {
		return fmt.Errorf("the CA certificate not found at path %#v", caCertFile)
	}
	if _, err := os.Stat(caKeyFile); err != nil {
		return fmt.Errorf("the CA private key not found at path %#v", caKeyFile)
	}
	if _, err := os.Stat(configFile); err != nil {
		return fmt.Errorf("the CA config file not found at path %#v", configFile)
	}
	cfg := &cli.Config{
		CAFile:    caCertFile,
		CAKeyFile: caKeyFile,
	}
	if cfg.CFG, err = config.LoadFile(configFile); err != nil {
		return
	}

	certSigner, err = universal.NewSigner(cli.RootFromConfig(cfg), cfg.CFG.Signing)
	return
}

func Sign(csr []byte, subject signer.Subject, profile string) ([]byte, error) {
	req := signer.SignRequest{
		Request: string(csr),
		Subject: &subject,
		Profile: profile,
	}
	return certSigner.Sign(req)
}

// SignToken sign a given token returning the certificate
func SignToken(db *sql.DB, token string, csr []byte) (*[]byte, string, error) {
	// check db
	if db == nil {
		return nil, "", errors.New("db connection is nil")
	}

	// create db transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback() //#nosec
		} else {
			tx.Commit() //#nosec
		}
	}()

	// retrieve db data
	var profile string
	var subjectData []byte
	err = tx.QueryRow(ownDb.GetTokenQuery, token).Scan(&profile, &subjectData)

	switch {
	case err == sql.ErrNoRows:
		//httpError(w, 403, errors.New("token not found"))
		return nil, "", SignForbidden{Msg: "Token not found"}
	case err != nil:
		// httpError(w, 500, err)
		return nil, "", err
	}

	var subject signer.Subject
	err = json.Unmarshal(subjectData, &subject)
	if err != nil {
		return nil, "", err
	}

	pemCert, err := Sign(csr, subject, profile)
	if err != nil {
		return nil, "", err
	}

	certData, _ := pem.Decode(pemCert)
	if certData == nil {
		//httpError(w, 500, errors.New("Failed to parse PEM encoded certificate."))
		return nil, "", errors.New("failed to parse PEM encoded certificate")
	}

	x509Cert, err := x509.ParseCertificate(certData.Bytes)
	if err != nil {
		//httpError(w, 500, errors.New("Failed to parse signed certificate."))
		return nil, "", errors.New("failed to parse signed certificate")
	}
	certSubject := x509Cert.Subject

	_, err = tx.Exec(ownDb.InsertCertificateQuery,
		certificateFingerprint(*x509Cert),
		certSubject.CommonName,
		firstOrNull(certSubject.Country),
		firstOrNull(certSubject.Locality),
		firstOrNull(certSubject.Organization),
		firstOrNull(certSubject.OrganizationalUnit),
		x509Cert.NotBefore,
		x509Cert.NotAfter,
		pemCert,
	)
	if err != nil {
		//httpError(w, 500, err)
		return nil, "", err
	}

	_, err = tx.Exec("DELETE FROM tokens where id=$1", token)
	if err != nil {
		//httpError(w, 500, errors.New("Failed to delete token."))
		return nil, "", errors.New("failed to delete token")
	}

	//get the signing CA
	caInfo, err := certSigner.Info(info.Req{Profile: profile})
	if err != nil {
		//httpError(w, 500, err)
		return nil, "", err
	}

	return &pemCert, caInfo.Certificate, nil
}

func firstOrNull(s []string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{String: s[0], Valid: true}
}

func PruneCertificates(db *sql.DB) (int64, error) {
	if db == nil {
		return 0, errors.New("clean PKI tokens: db connection is nil")
	}

	res, err := db.Exec(ownDb.CleanPkiCertificatesQuery)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
