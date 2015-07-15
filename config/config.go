package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"runtime"
)

type Config struct {
	Endpoints    []string
	CACerts      *x509.CertPool
	ClientCert   *tls.Certificate
	Transport    string
	Identity     string
	Project      string
	Organization string
	LogLevel     string
}

func (config *Config) Load(c *cli.Context) error {
	config.Endpoints = c.GlobalStringSlice("endpoint")
	config.Transport = c.GlobalString("transport")

	if c.GlobalString("tls-client-cert") != "" || c.GlobalString("tls-client-key") != "" || c.GlobalString("tls-ca-cert") != "" {
		if err := config.loadTLSConfig(c.GlobalString("tls-client-cert"), c.GlobalString("tls-client-key"), c.GlobalString("tls-ca-cert")); err != nil {
			return err
		}
	} else {
		//This is only for testing when running without a tls certificate
		config.Identity = runtime.GOOS
		config.Project = "test-project"
		config.Organization = "test-org"
	}

	config.LogLevel = c.GlobalString("log-level")

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("Endpoints: %s, CACerts: %s, ClientCert: %s, Transport: %s, Identity: %s, Project: %s, Organization: %s, LogLevel: %s", c.Endpoints, c.CACerts != nil, c.ClientCert != nil, c.Transport, c.Identity, c.Project, c.Organization, c.LogLevel)
}

func (c *Config) loadTLSConfig(client_cert, client_key, ca_certs string) error {
	cert, err := tls.LoadX509KeyPair(client_cert, client_key)
	if err != nil {
		return fmt.Errorf("Failed to load client certificate/key: %s", err)
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("Failed to parse client certificate: %s", err)
	}
	c.ClientCert = &cert
	//Extract org, project and identity from the client cert
	if len(cert.Leaf.Subject.Organization) > 0 {
		c.Organization = cert.Leaf.Subject.Organization[0]
	}
	if len(cert.Leaf.Subject.OrganizationalUnit) > 0 {
		c.Project = cert.Leaf.Subject.OrganizationalUnit[0]
	}
	c.Identity = cert.Leaf.Subject.CommonName
	pemCerts, err := ioutil.ReadFile(ca_certs)
	if err != nil {
		return fmt.Errorf("Failed to load CA certificate: %s", err)
	}
	certpool := x509.NewCertPool()
	if !certpool.AppendCertsFromPEM(pemCerts) {
		return fmt.Errorf("Given CA file does not contain a PEM encoded x509 certificate")
	}
	c.CACerts = certpool

	return nil
}
