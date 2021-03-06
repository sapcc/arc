package config

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/codegangsta/cli"
)

type Config struct {
	Endpoints      []string
	CACerts        *x509.CertPool
	ClientCert     *tls.Certificate
	ClientKey      *rsa.PrivateKey
	ClientCertPath string
	Transport      string
	Identity       string
	Project        string
	Organization   string
	LogLevel       string
}

func New() Config {
	//FIXME: This is only for testing when running without a tls certificate
	//Should be moved to test files at some point
	return Config{
		Identity:     runtime.GOOS,
		Project:      "",
		Organization: "",
	}
}

func (config *Config) Load(c *cli.Context) error {
	if len(c.StringSlice("endpoint")) > 0 {
		config.Endpoints = c.StringSlice("endpoint")
	}
	if c.String("transport") != "" {
		config.Transport = c.String("transport")
	}

	if c.String("tls-client-cert") != "" || c.String("tls-client-key") != "" || c.String("tls-ca-cert") != "" {
		if err := config.loadTLSConfig(c.String("tls-client-cert"), c.String("tls-client-key"), c.String("tls-ca-cert")); err != nil {
			return err
		}
	}

	config.LogLevel = c.GlobalString("log-level")

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("Endpoints: %s, CACerts: %t, ClientCert: %t, Transport: %s, Identity: %s, Project: %s, Organization: %s, LogLevel: %s", c.Endpoints, c.CACerts != nil, c.ClientCert != nil, c.Transport, c.Identity, c.Project, c.Organization, c.LogLevel)
}

func (c *Config) loadTLSConfig(client_cert, client_key, ca_certs string) error {
	cert, err := tls.LoadX509KeyPair(client_cert, client_key)
	if err != nil {
		return fmt.Errorf("failed to load client certificate/key: %s", err)
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("failed to parse client certificate: %s", err)
	}
	c.ClientCert = &cert

	// save client cert path
	c.ClientCertPath = client_cert

	// Extract key
	keyPEMBlock, err := ioutil.ReadFile(filepath.Clean(client_key))
	if err != nil {
		return fmt.Errorf("failed to load client key: %s", err)
	}
	keyPemDecoded, _ := pem.Decode(keyPEMBlock)
	if keyPemDecoded == nil {
		return fmt.Errorf("no pem block found")
	}
	c.ClientKey, err = x509.ParsePKCS1PrivateKey(keyPemDecoded.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse rsa key: %s", err)
	}

	//Extract org, project and identity from the client cert
	if len(cert.Leaf.Subject.Organization) > 0 {
		c.Organization = cert.Leaf.Subject.Organization[0]
	}
	if len(cert.Leaf.Subject.OrganizationalUnit) > 0 {
		c.Project = cert.Leaf.Subject.OrganizationalUnit[0]
	}
	c.Identity = cert.Leaf.Subject.CommonName
	pemCerts, err := ioutil.ReadFile(filepath.Clean(ca_certs))
	if err != nil {
		return fmt.Errorf("failed to load CA certificate: %s", err)
	}
	certpool := x509.NewCertPool()
	if !certpool.AppendCertsFromPEM(pemCerts) {
		return fmt.Errorf("given CA file does not contain a PEM encoded x509 certificate")
	}
	c.CACerts = certpool

	return nil
}
