package pki

/*
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
	"github.com/cloudflare/cfssl/log"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/universal"
	"github.com/gorilla/mux"
)

func signHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	token := vars["token"]

	csr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpError(w, 500, err)
		return
	}
	r.Body.Close()

	tx, err := db.Begin()
	if err != nil {
		httpError(w, 500, err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var profile string
	var subjectData []byte

	err = tx.QueryRow(fmt.Sprintf("SELECT profile, subject FROM tokens WHERE id=$1 AND created_at > NOW() - INTERVAL '%s' FOR UPDATE", TokenLifetime), token).Scan(&profile, &subjectData)

	switch {
	case err == sql.ErrNoRows:
		httpError(w, 403, errors.New("token not found"))
		return
	case err != nil:
		httpError(w, 500, err)
		return
	}

	var subject signer.Subject
	err = json.Unmarshal(subjectData, &subject)
	if err != nil {
		httpError(w, 500, err)
		return
	}

	root := cli.RootFromConfig(&cfg)

	s, err := universal.NewSigner(root, cfg.CFG.Signing)
	if err != nil {
		httpError(w, 500, err)
		return
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
		httpError(w, 500, err)
		return
	}

	certData, _ := pem.Decode(pemCert)
	if certData == nil {
		httpError(w, 500, errors.New("Failed to parse PEM encoded certificate."))
		return
	}

	x509Cert, err := x509.ParseCertificate(certData.Bytes)
	if err != nil {
		httpError(w, 500, errors.New("Failed to parse signed certificate."))
		return
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
		httpError(w, 500, err)
		return
	}

	_, err = tx.Exec("DELETE FROM tokens where id=$1", token)

	if err != nil {
		httpError(w, 500, errors.New("Failed to delete token."))
		return
	}

	acceptHeader := ""
	for _, v := range r.Header["Accept"] {
		acceptHeader = v
		break
	}

	if acceptHeader == "application/json" {
		//get the signing CA
		caInfo, err := s.Info(info.Req{Profile: profile})
		if err != nil {
			httpError(w, 500, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"certificate": string(pemCert),
			"ca":          caInfo.Certificate,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		log.Error("plain")
		w.Header().Set("Content-Type", "application/pkix-cert")
		w.Write(pemCert)
	}

}*/
