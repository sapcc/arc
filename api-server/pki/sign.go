package pki

import (
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/info"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/universal"
)

// SignForbidden should be used to return a 403
type SignForbidden struct {
	Msg string
}

func (e SignForbidden) Error() string {
	return e.Msg
}

// TokenLifetime should set the standard token time life time
var TokenLifetime = "1 hour"

// SignToken sign a given token returning the certificate
func SignToken(db *sql.DB, token string, r *http.Request, cfg *cli.Config) (*[]byte, string, error) {
	// check db
	if db == nil {
		return nil, "", errors.New("Db connection is nil")
	}

	// get the csr
	csr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// httpError(w, 500, err)
		return nil, "", err
	}
	r.Body.Close()

	// create db transaction
	tx, err := db.Begin()
	if err != nil {
		//httpError(w, 500, err)
		return nil, "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// retrieve db data
	var profile string
	var subjectData []byte
	err = tx.QueryRow(fmt.Sprintf("SELECT profile, subject FROM tokens WHERE id=$1 AND created_at > NOW() - INTERVAL '%s' FOR UPDATE", TokenLifetime), token).Scan(&profile, &subjectData)

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
		// httpError(w, 500, err)
		return nil, "", err
	}

	root := cli.RootFromConfig(cfg)

	// check if cfg is nil else panic!!
	if cfg.CFG == nil {
		return nil, "", errors.New("Signer configuration is nil.")
	}

	s, err := universal.NewSigner(root, cfg.CFG.Signing)
	if err != nil {
		//httpError(w, 500, err)
		return nil, "", err
	}

	req := signer.SignRequest{
		// Hosts:   signer.SplitHosts(c.Hostname),
		Request: string(csr),
		Subject: &subject,
		Profile: profile,
		// Label:   c.Label,
	}
	pemCert, err := s.Sign(req)
	if err != nil {
		//httpError(w, 500, err)
		return nil, "", err
	}

	certData, _ := pem.Decode(pemCert)
	if certData == nil {
		//httpError(w, 500, errors.New("Failed to parse PEM encoded certificate."))
		return nil, "", errors.New("Failed to parse PEM encoded certificate.")
	}

	x509Cert, err := x509.ParseCertificate(certData.Bytes)
	if err != nil {
		//httpError(w, 500, errors.New("Failed to parse signed certificate."))
		return nil, "", errors.New("Failed to parse signed certificate.")
	}
	certSubject := x509Cert.Subject

	_, err = tx.Exec(`INSERT into certificates
	(fingerprint, common_name, country, locality, organization, organizational_unit, not_before, not_after, pem)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
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
		return nil, "", errors.New("Failed to delete token.")
	}

	//get the signing CA
	caInfo, err := s.Info(info.Req{Profile: profile})
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
