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
	RENEW_CFG_PRIVKEY_MISSING     = "Configuration is nil or TLS client private key not found."
	RENEW_TLS_CERTIFICATE_MISSING = "No TLS Certificate found to check expiration date."
	RENEW_CFG_CERT_PATH_MISSING   = "Configuration is nil or client cert path is missing."
)

// renewThreshold in hours
func RenewCert(cfg *arc_config.Config, renewURI string, renewThreshold int64, httpClientInsecureSkipVerify bool) (bool, int64, error) {
	// first check expiration date
	hoursLeft, err := CertExpirationDate(cfg)
	if err != nil {
		return false, 0, err
	}
	// renew cert if experition time is less than the given hours
	if hoursLeft > renewThreshold {
		return false, hoursLeft, nil
	}

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
		return false, hoursLeft, err
	}

	err = SaveCertificate(certPEMBlock, cfg)
	if err != nil {
		return false, hoursLeft, err
	}

	return true, hoursLeft, nil
}

// returns expiration time in hours (int64)
func CertExpirationDate(cfg *arc_config.Config) (int64, error) {
	if cfg.ClientCert == nil || len(cfg.ClientCert.Certificate) == 0 {
		return 0, errors.New(RENEW_TLS_CERTIFICATE_MISSING)
	}

	cert, err := x509.ParseCertificate(cfg.ClientCert.Certificate[0])
	if err != nil {
		return 0, fmt.Errorf("Failed to parse client certificate: %s", err)
	}

	expiresIn := int64(time.Until(cert.NotAfter).Hours())

	return expiresIn, nil
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
		return nil, fmt.Errorf("No pem block found.")
	}

	return pem.EncodeToMemory(pBlockCert), nil
}
