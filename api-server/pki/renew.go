package pki

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
)

var (
	RENEW_CFG_PRIVKEY_MISSING     = "configuration is nil or TLS client private key not found."
	RENEW_TLS_CERTIFICATE_MISSING = "no TLS Certificate found to check expiration date."
	RENEW_CFG_CERT_PATH_MISSING   = "configuration is nil or client cert path is missing."
)

// CheckAndRenewCert check with the threshold and renew the cert
// int64 --> hours left to the expiration date. If int64 > 0 means that hoursLeft > threshold and there is no need to renew the cert
// error --> something wrong happend
func CheckAndRenewCert(cfg *arc_config.Config, renewURI string, renewThreshold int64, httpClientInsecureSkipVerify bool) (int64, error) {
	// first check expiration date
	notAfter, err := CertExpirationDate(cfg)
	if err != nil {
		return 0, err
	}

	// hours negative (ex: -5) the certificate is already 5 hours expired
	hoursLeft := CertExpiresIn(notAfter)

	// renew cert if experition time is less than the given hours
	if hoursLeft > renewThreshold {
		return hoursLeft, nil
	}
	// renew Cert
	err = RenewCert(cfg, renewURI, httpClientInsecureSkipVerify)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// RenewCert renew the cert
func RenewCert(cfg *arc_config.Config, renewURI string, httpClientInsecureSkipVerify bool) error {
	// #nosec: TLS InsecureSkipVerify set true
	// client creation
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{*cfg.ClientCert},
				InsecureSkipVerify: httpClientInsecureSkipVerify,
			},
		},
	}

	// lets get the new cert
	certPEMBlock, err := SendCertificateRequest(client, renewURI, cfg)
	if err != nil {
		return err
	}

	return SaveCertificate(certPEMBlock, cfg)
}

// CertExpiresIn returns expiration time in hours (int64)
func CertExpiresIn(notAfter *time.Time) int64 {
	return int64(time.Until(*notAfter).Hours())
}

// CertExpirationDate return the notAfter attribute of the cert
func CertExpirationDate(cfg *arc_config.Config) (*time.Time, error) {
	if cfg == nil || cfg.ClientCert == nil || len(cfg.ClientCert.Certificate) == 0 {
		return nil, errors.New(RENEW_TLS_CERTIFICATE_MISSING)
	}

	cert, err := x509.ParseCertificate(cfg.ClientCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse client certificate: %s", err)
	}

	return &cert.NotAfter, nil
}

func SaveCertificate(certPEMBlock []byte, cfg *arc_config.Config) error {
	if cfg == nil || cfg.ClientCertPath == "" {
		return errors.New(RENEW_CFG_CERT_PATH_MISSING)
	}

	return ioutil.WriteFile(cfg.ClientCertPath, certPEMBlock, 0644)
}

func SendCertificateRequest(client *http.Client, endpoint string, cfg *arc_config.Config) ([]byte, error) {
	if cfg == nil || cfg.ClientKey == nil {
		return nil, errors.New(RENEW_CFG_PRIVKEY_MISSING)
	}

	// create cert request
	csrData, err := CreateSignReqCert("", "", "", cfg.ClientKey)
	if err != nil {
		return nil, err
	}

	// send post with a sign request
	res, err := client.Post(endpoint, "application/pkix-cert", bytes.NewReader(csrData))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// validate cert with private key
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	clientKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(cfg.ClientKey),
	})
	_, err = tls.X509KeyPair(responseBody, clientKeyPEM)
	if err != nil {
		return nil, err
	}

	// decode cert
	pBlockCert, _ := pem.Decode(responseBody)
	if pBlockCert == nil {
		return nil, fmt.Errorf("no pem block found")
	}

	return pem.EncodeToMemory(pBlockCert), nil
}
