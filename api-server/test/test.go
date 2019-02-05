package test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

func CreateCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (
	cert *x509.Certificate, certPEM []byte, err error) {

	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	}
	// parse the resulting certificate so we can use it again
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}
	// PEM encode the certificate (this is a standard TLS encoding)
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}

// Suported
// SignatureAlgorithm:    x509.SHA256WithRSA,
// SignatureAlgorithm:    x509.ECDSAWithSHA256,
func CertTemplate(signAlgo x509.SignatureAlgorithm) (*x509.Certificate, error) {
	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Test, Inc."}},
		SignatureAlgorithm:    signAlgo,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour), // valid for an hour
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

//
// return a tls httptest.NewUnstartedServer
// defer ts.Close()
//
func NewUnstartedTestTLSServer(certPool *x509.CertPool, servTLSCert *tls.Certificate, handler http.Handler) *httptest.Server {
	ts := httptest.NewUnstartedServer(handler)
	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{*servTLSCert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ClientCAs:    certPool,
	}
	return ts
}

// Optional cert template
// returns
// rootCRT: to generate the client crt
// RootCRTPEM: to generate the trusted cert pool
// Rootkey: to generate the client crt
// error
func RootTLSCRT(rootCertTmpl *x509.Certificate) (*x509.Certificate, []byte, interface{}, error) {
	var err error
	// self signed cert
	if rootCertTmpl == nil {
		rootCertTmpl, err = CertTemplate(x509.ECDSAWithSHA256)
		if err != nil {
			return nil, make([]byte, 0), nil, fmt.Errorf("creating cert template: %v", err)
		}
	}
	// root private and public key
	rootKey, rootPublicKey, err := generateKeys(rootCertTmpl.SignatureAlgorithm)
	if err != nil {
		return nil, make([]byte, 0), nil, fmt.Errorf("generating random key: %v", err)
	}

	// describe what the certificate will be used for
	rootCertTmpl.IsCA = true
	rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	rootCertTmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	// create our self-signed cert
	rootCert, rootCertPEM, err := CreateCert(rootCertTmpl, rootCertTmpl, rootPublicKey, rootKey)
	if err != nil {
		return nil, make([]byte, 0), nil, fmt.Errorf("error creating cert: %v", err)
	}
	return rootCert, rootCertPEM, rootKey, nil
}

// Optional cert template
func ServerTLSCRT(rootCert *x509.Certificate, rootKey interface{}, servCertTmpl *x509.Certificate) (*tls.Certificate, []byte, interface{}, error) {
	var err error
	if servCertTmpl == nil {
		servCertTmpl, err = CertTemplate(x509.ECDSAWithSHA256)
		if err != nil {
			return nil, make([]byte, 0), nil, fmt.Errorf("creating cert template: %v", err)
		}
	}
	// private and public key
	servKey, servPublicKey, err := generateKeys(servCertTmpl.SignatureAlgorithm)
	if err != nil {
		return nil, make([]byte, 0), nil, fmt.Errorf("generating random key: %v", err)
	}

	// describe what the certificate will be used for
	servCertTmpl.KeyUsage = x509.KeyUsageDigitalSignature
	servCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	servCertTmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	// create a certificate which wraps the server's public key, sign it with the root private key
	_, servCertPEM, err := CreateCert(servCertTmpl, rootCert, servPublicKey, rootKey)
	if err != nil {
		return nil, make([]byte, 0), nil, fmt.Errorf("error creating cert: %v", err)
	}
	// create key PEM
	var servKeyPEM []byte
	if servCertTmpl.SignatureAlgorithm == x509.ECDSAWithSHA256 {
		key := servKey.(*ecdsa.PrivateKey)
		serverKeyDer, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			return nil, make([]byte, 0), nil, fmt.Errorf("failed to serialize ECDSA key: %s", err)
		}
		servKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type: "EC PRIVATE KEY", Bytes: serverKeyDer,
		})
	} else {
		key := servKey.(*rsa.PrivateKey)
		servKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
	}

	servTLSCert, err := tls.X509KeyPair(servCertPEM, servKeyPEM)
	if err != nil {
		return nil, make([]byte, 0), nil, fmt.Errorf("invalid key pair: %v", err)
	}
	return &servTLSCert, servCertPEM, servKey, nil
}

// Optional cert template
func ClientTLSCRT(rootCert *x509.Certificate, rootKey interface{}, clientCertTmpl *x509.Certificate) (*tls.Certificate, []byte, interface{}, error) {
	var err error
	// create a template for the client
	if clientCertTmpl == nil {
		clientCertTmpl, err = CertTemplate(x509.ECDSAWithSHA256)
		if err != nil {
			return nil, make([]byte, 0), nil, fmt.Errorf("creating cert template: %v", err)
		}
	}
	// private and public key
	clientKey, clientPublicKey, err := generateKeys(clientCertTmpl.SignatureAlgorithm)
	if err != nil {
		return nil, make([]byte, 0), nil, fmt.Errorf("generating random key: %v", err)
	}

	clientCertTmpl.KeyUsage = x509.KeyUsageDigitalSignature
	clientCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	// the root cert signs the cert by again providing its private key
	_, clientCertPEM, err := CreateCert(clientCertTmpl, rootCert, clientPublicKey, rootKey)
	if err != nil {
		return nil, make([]byte, 0), nil, fmt.Errorf("error creating cert: %v", err)
	}

	// create key PEM
	var clientKeyPEM []byte
	if clientCertTmpl.SignatureAlgorithm == x509.ECDSAWithSHA256 {
		key := clientKey.(*ecdsa.PrivateKey)
		clientKeyDer, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			return nil, make([]byte, 0), nil, fmt.Errorf("failed to serialize ECDSA key: %s", err)
		}
		clientKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type: "EC PRIVATE KEY", Bytes: clientKeyDer,
		})
	} else {
		key := clientKey.(*rsa.PrivateKey)
		clientKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
	}

	// client authentication cert
	clientTLSCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return nil, make([]byte, 0), nil, fmt.Errorf("invalid key pair: %v", err)
	}
	return &clientTLSCert, clientCertPEM, clientKey, nil
}

func NewTestTLSClient(certPool *x509.CertPool, clientTLSCert *tls.Certificate) *http.Client {
	// create client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{*clientTLSCert},
				RootCAs:      certPool,
			},
		},
	}

	return client
}

type TestTLSServer struct {
	Server  *httptest.Server
	Crt     *tls.Certificate
	PrivKey interface{}
	Subject pkix.Name
}

type TestTLSClient struct {
	Client  *http.Client
	Crt     *tls.Certificate
	PrivKey interface{}
	Subject pkix.Name
}

//
// defer ts.Close() on succeess
//
func NewTLSCommunication(handler http.Handler) (TestTLSServer, TestTLSClient, error) {
	// create root crt
	rootCert, rootCertPEM, rootKey, err := RootTLSCRT(nil)
	if err != nil {
		return TestTLSServer{}, TestTLSClient{}, err
	}
	// pool of trusted certs
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(rootCertPEM) {
		return TestTLSServer{}, TestTLSClient{}, fmt.Errorf("given CA file does not contain a PEM encoded x509 certificate")
	}
	// create server crt
	servCertTmpl, err := CertTemplate(x509.ECDSAWithSHA256)
	if err != nil {
		return TestTLSServer{}, TestTLSClient{}, fmt.Errorf("creating cert template: %v", err)
	}
	servCertTmpl.Subject = pkix.Name{Organization: []string{"test_server_o"}, OrganizationalUnit: []string{"test_server_ou"}, CommonName: "Test_server_cn"}
	servTLSCert, _, servPrivKey, err := ServerTLSCRT(rootCert, rootKey, servCertTmpl)
	if err != nil {
		return TestTLSServer{}, TestTLSClient{}, err
	}
	// create client crt
	clientCertTmpl, err := CertTemplate(x509.ECDSAWithSHA256)
	if err != nil {
		return TestTLSServer{}, TestTLSClient{}, fmt.Errorf("creating cert template: %v", err)
	}
	clientCertTmpl.Subject = pkix.Name{Organization: []string{"test_client_o"}, OrganizationalUnit: []string{"test_client_ou"}, CommonName: "Test_client_cn"}
	clientTLSCert, _, clientPrivKey, err := ClientTLSCRT(rootCert, rootKey, clientCertTmpl)
	if err != nil {
		return TestTLSServer{}, TestTLSClient{}, err
	}
	// create server
	ts := NewUnstartedTestTLSServer(certPool, servTLSCert, handler)
	// create client
	client := NewTestTLSClient(certPool, clientTLSCert)
	return TestTLSServer{
			ts,
			servTLSCert,
			servPrivKey,
			servCertTmpl.Subject,
		}, TestTLSClient{
			client,
			clientTLSCert,
			clientPrivKey,
			clientCertTmpl.Subject,
		}, nil
}

func generateKeys(signAlgo x509.SignatureAlgorithm) (interface{}, interface{}, error) {
	var err error
	var privKey interface{}
	var pubKey interface{}
	if signAlgo == x509.ECDSAWithSHA256 {
		privKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		pubKey = &privKey.(*ecdsa.PrivateKey).PublicKey
	} else {
		privKey, err = rsa.GenerateKey(rand.Reader, 2048)
		pubKey = &privKey.(*rsa.PrivateKey).PublicKey
	}
	if err != nil {
		return nil, nil, fmt.Errorf("generating random key: %v", err)
	}
	return privKey, pubKey, nil
}
