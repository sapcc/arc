// +build integration

package pki_test

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gitHub.***REMOVED***/monsoon/arc/api-server/auth"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/pki"
	"gitHub.***REMOVED***/monsoon/arc/api-server/test"

	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
)

var _ = Describe("CheckAndRenewCert", func() {

	It("should update the cert", func() {
		notAfter := time.Now().Add(600 * time.Hour)
		rootCertFile, servCertFile, servKeyFile, clientCertFile, clientKeyFile, err := createCfgTmpCerts(&notAfter)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(rootCertFile.Name())
		defer os.Remove(servCertFile.Name())
		defer os.Remove(servKeyFile.Name())
		defer os.Remove(clientCertFile.Name())
		defer os.Remove(clientKeyFile.Name())

		// save cert notAfter from creation
		clientCertPair, err := tls.LoadX509KeyPair(clientCertFile.Name(), clientKeyFile.Name())
		Expect(err).NotTo(HaveOccurred())
		cert, err := x509.ParseCertificate(clientCertPair.Certificate[0])
		Expect(err).NotTo(HaveOccurred())
		clientCertNotAfter := cert.NotAfter
		Expect(int64(clientCertNotAfter.Sub(time.Now()).Hours())).To(BeNumerically("==", 599)) // expire the certificate

		conf, err := testConfig(rootCertFile, clientCertFile, clientKeyFile)
		Expect(err).NotTo(HaveOccurred())

		// read manuelly server cert
		servCert, err := tls.LoadX509KeyPair(servCertFile.Name(), servKeyFile.Name())
		Expect(err).NotTo(HaveOccurred())

		ts := test.NewUnstartedTestTLSServer(conf.CACerts, &servCert, http.HandlerFunc(renewPkiCertTest))
		ts.StartTLS()
		defer ts.Close()
		Expect(err).NotTo(HaveOccurred())

		hoursLeft, err := CheckAndRenewCert(conf, ts.URL, 744, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(hoursLeft).To(BeNumerically("==", int64(0)))

		// test cert NotAfter
		clientCertPair, err = tls.LoadX509KeyPair(clientCertFile.Name(), clientKeyFile.Name())
		Expect(err).NotTo(HaveOccurred())
		cert, err = x509.ParseCertificate(clientCertPair.Certificate[0])
		Expect(err).NotTo(HaveOccurred())
		renewedClientCertNotAfter := cert.NotAfter

		// we should check the expire date from clientCertFile
		diff := renewedClientCertNotAfter.Sub(time.Now())
		Expect(int64(diff.Hours())).To(BeNumerically("==", 17519))
	})

	It("should NOT update the cert", func() {
		notAfter := time.Now().Add(600 * time.Hour)
		rootCertFile, servCertFile, servKeyFile, clientCertFile, clientKeyFile, err := createCfgTmpCerts(&notAfter)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(rootCertFile.Name())
		defer os.Remove(servCertFile.Name())
		defer os.Remove(servKeyFile.Name())
		defer os.Remove(clientCertFile.Name())
		defer os.Remove(clientKeyFile.Name())

		// save cert notAfter from creation
		clientCertPair, err := tls.LoadX509KeyPair(clientCertFile.Name(), clientKeyFile.Name())
		Expect(err).NotTo(HaveOccurred())
		cert, err := x509.ParseCertificate(clientCertPair.Certificate[0])
		Expect(err).NotTo(HaveOccurred())
		clientCertNotAfter := cert.NotAfter
		// test cert not after
		Expect(int64(clientCertNotAfter.Sub(time.Now()).Hours())).To(BeNumerically("==", 599)) // NOT expire the certificate

		conf, err := testConfig(rootCertFile, clientCertFile, clientKeyFile)
		Expect(err).NotTo(HaveOccurred())

		// read manuelly server cert
		servCert, err := tls.LoadX509KeyPair(servCertFile.Name(), servKeyFile.Name())
		Expect(err).NotTo(HaveOccurred())

		ts := test.NewUnstartedTestTLSServer(conf.CACerts, &servCert, http.HandlerFunc(renewPkiCertTest))
		ts.StartTLS()
		defer ts.Close()
		Expect(err).NotTo(HaveOccurred())

		hoursLeft, err := CheckAndRenewCert(conf, ts.URL, 500, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(hoursLeft).To(BeNumerically("==", int64(599)))
	})

})

var _ = Describe("CertExpirationDate", func() {

	It("should return left hours to the expiration date", func() {
		notAfter := time.Now().Add(5 * time.Hour) // expiration in 5 hours
		rootCertFile, servCertFile, servKeyFile, clientCertFile, clientKeyFile, err := createCfgTmpCerts(&notAfter)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(rootCertFile.Name())
		defer os.Remove(servCertFile.Name())
		defer os.Remove(servKeyFile.Name())
		defer os.Remove(clientCertFile.Name())
		defer os.Remove(clientKeyFile.Name())

		conf, err := testConfig(rootCertFile, clientCertFile, clientKeyFile)
		Expect(err).NotTo(HaveOccurred())

		expiration, err := CertExpirationDate(conf)
		Expect(err).NotTo(HaveOccurred())
		Expect(int64(time.Until(*expiration).Hours())).To(Equal(int64(4)))
	})

	It("should return expired hours", func() {
		notAfter := time.Now().Add(-5 * time.Hour)
		rootCertFile, servCertFile, servKeyFile, clientCertFile, clientKeyFile, err := createCfgTmpCerts(&notAfter)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(rootCertFile.Name())
		defer os.Remove(servCertFile.Name())
		defer os.Remove(servKeyFile.Name())
		defer os.Remove(clientCertFile.Name())
		defer os.Remove(clientKeyFile.Name())

		conf, err := testConfig(rootCertFile, clientCertFile, clientKeyFile)
		Expect(err).NotTo(HaveOccurred())

		expiration, err := CertExpirationDate(conf)
		Expect(err).NotTo(HaveOccurred())
		Expect(int64(time.Until(*expiration).Hours())).To(Equal(int64(-5)))
	})

})

var _ = Describe("SendCertificateRequest", func() {

	It("should return error if cfg is nil", func() {
		rootCertFile, servCertFile, servKeyFile, clientCertFile, clientKeyFile, err := createCfgTmpCerts(nil)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(rootCertFile.Name())
		defer os.Remove(servCertFile.Name())
		defer os.Remove(servKeyFile.Name())
		defer os.Remove(clientCertFile.Name())
		defer os.Remove(clientKeyFile.Name())

		conf, err := testConfig(rootCertFile, clientCertFile, clientKeyFile)
		Expect(err).NotTo(HaveOccurred())

		ts := test.NewUnstartedTestTLSServer(conf.CACerts, conf.ClientCert, http.HandlerFunc(renewPkiCertTest))
		client := test.NewTestTLSClient(conf.CACerts, conf.ClientCert)

		_, err = SendCertificateRequest(client, ts.URL, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(RENEW_CFG_PRIVKEY_MISSING))
	})

	It("should return error if cfg priv key missing", func() {
		rootCertFile, servCertFile, servKeyFile, clientCertFile, clientKeyFile, err := createCfgTmpCerts(nil)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(rootCertFile.Name())
		defer os.Remove(servCertFile.Name())
		defer os.Remove(servKeyFile.Name())
		defer os.Remove(clientCertFile.Name())
		defer os.Remove(clientKeyFile.Name())

		conf, err := testConfig(rootCertFile, clientCertFile, clientKeyFile)
		Expect(err).NotTo(HaveOccurred())

		conf.ClientKey = nil

		ts := test.NewUnstartedTestTLSServer(conf.CACerts, conf.ClientCert, http.HandlerFunc(renewPkiCertTest))
		client := test.NewTestTLSClient(conf.CACerts, conf.ClientCert)

		_, err = SendCertificateRequest(client, ts.URL, conf)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(RENEW_CFG_PRIVKEY_MISSING))
	})

	It("should return a certificate", func() {
		// create certs for the config
		rootCertFile, servCertFile, servKeyFile, clientCertFile, clientKeyFile, err := createCfgTmpCerts(nil)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(rootCertFile.Name())
		defer os.Remove(servCertFile.Name())
		defer os.Remove(servKeyFile.Name())
		defer os.Remove(clientCertFile.Name())
		defer os.Remove(clientKeyFile.Name())

		// create config
		conf, err := testConfig(rootCertFile, clientCertFile, clientKeyFile)
		Expect(err).NotTo(HaveOccurred())

		// get a tls communication between server and client
		ts, tc, err := test.NewTLSCommunication(http.HandlerFunc(renewPkiCertTest))
		ts.Server.StartTLS()
		defer ts.Server.Close()
		Expect(err).NotTo(HaveOccurred())

		_, err = SendCertificateRequest(tc.Client, ts.Server.URL, conf)
		Expect(err).NotTo(HaveOccurred())
	})

})

// just sign a certificate request with mock data
func renewPkiCertTest(w http.ResponseWriter, r *http.Request) {
	// create a TokenRequest with the CN
	var tr TokenRequest
	tr.CN = "test_cn"

	// add OU and O over the authorization because it will override the tokenRequest
	auth := auth.Authorization{}
	auth.ProjectId = "test_ou"
	auth.ProjectDomainId = "test_o"

	// read request body
	csr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a token
	token, err := CreateToken(db, &auth, tr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// sign tocken with csr
	pemCert, _, err := SignToken(db, token, csr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pkix-cert")
	w.Write(*pemCert)
}

// default notAfter time.Now().Add(time.Hour)
// defer os.Remove(rootCertFile.Name())
// defer os.Remove(servCertFile.Name())
// defer os.Remove(servKeyFile.Name())
func createCfgTmpCerts(notAfter *time.Time) (*os.File, *os.File, *os.File, *os.File, *os.File, error) {
	// create root crt
	rootCertTmpl, err := test.CertTemplate(x509.SHA256WithRSA)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	rootCert, rootCertPEM, rootKey, err := test.RootTLSCRT(rootCertTmpl)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	// create server crt
	servCertTmpl, err := test.CertTemplate(x509.SHA256WithRSA)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	servCertTmpl.Subject = pkix.Name{Organization: []string{"test_server_o"}, OrganizationalUnit: []string{"test_server_ou"}, CommonName: "Test_server_cn"}
	_, servCertPEM, servPrivKey, err := test.ServerTLSCRT(rootCert, rootKey, servCertTmpl)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	servKey := servPrivKey.(*rsa.PrivateKey)
	servKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(servKey),
	})

	// create client crt
	clientCertTmpl, err := test.CertTemplate(x509.SHA256WithRSA)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	clientCertTmpl.Subject = pkix.Name{Organization: []string{"test_client_o"}, OrganizationalUnit: []string{"test_client_ou"}, CommonName: "Test_client_cn"}
	clientCertTmpl.NotAfter = time.Now().Add(time.Hour)
	if notAfter != nil {
		clientCertTmpl.NotAfter = *notAfter
	}
	_, clientCertPEM, clientPrivKey, err := test.ClientTLSCRT(rootCert, rootKey, clientCertTmpl)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	clientKey := clientPrivKey.(*rsa.PrivateKey)
	clientKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey),
	})

	// save root cert to tmp file
	rootCertFile, err := ioutil.TempFile(os.TempDir(), "rootCert")
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	err = ioutil.WriteFile(rootCertFile.Name(), rootCertPEM, 0644)
	if err != nil {
		os.Remove(rootCertFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	// save serv cert to tmp file
	servCertFile, err := ioutil.TempFile(os.TempDir(), "servCert")
	if err != nil {
		os.Remove(rootCertFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	err = ioutil.WriteFile(servCertFile.Name(), servCertPEM, 0644)
	if err != nil {
		os.Remove(rootCertFile.Name())
		os.Remove(servCertFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	// save server client cert
	servKeyFile, err := ioutil.TempFile(os.TempDir(), "servKey")
	if err != nil {
		os.Remove(rootCertFile.Name())
		os.Remove(servCertFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	err = ioutil.WriteFile(servKeyFile.Name(), servKeyPEM, 0644)
	if err != nil {
		os.Remove(rootCertFile.Name())
		os.Remove(servCertFile.Name())
		os.Remove(servKeyFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	// save client cert to tmp file
	clientCertFile, err := ioutil.TempFile(os.TempDir(), "clientCert")
	if err != nil {
		os.Remove(rootCertFile.Name())
		os.Remove(servCertFile.Name())
		os.Remove(servKeyFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	err = ioutil.WriteFile(clientCertFile.Name(), clientCertPEM, 0644)
	if err != nil {
		os.Remove(rootCertFile.Name())
		os.Remove(servCertFile.Name())
		os.Remove(servKeyFile.Name())
		os.Remove(clientCertFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	// save serv key to tmp file
	clientKeyFile, err := ioutil.TempFile(os.TempDir(), "clientKey")
	if err != nil {
		os.Remove(rootCertFile.Name())
		os.Remove(servCertFile.Name())
		os.Remove(servKeyFile.Name())
		os.Remove(clientCertFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	err = ioutil.WriteFile(clientKeyFile.Name(), clientKeyPEM, 0644)
	if err != nil {
		os.Remove(rootCertFile.Name())
		os.Remove(servCertFile.Name())
		os.Remove(servKeyFile.Name())
		os.Remove(clientCertFile.Name())
		os.Remove(clientKeyFile.Name())
		return nil, nil, nil, nil, nil, err
	}
	return rootCertFile, servCertFile, servKeyFile, clientCertFile, clientKeyFile, nil
}

func testConfig(rootCertFile, servCertFile, servKeyFile *os.File) (*arc_config.Config, error) {
	// prepare flags
	globalSet := flag.NewFlagSet("test", 0)
	globalSet.String("transport", "mqtt", "test")
	globalSet.String("tls-client-cert", servCertFile.Name(), "test")
	globalSet.String("tls-client-key", servKeyFile.Name(), "test")
	globalSet.String("tls-ca-cert", rootCertFile.Name(), "test")
	globalSet.String("log-level", "info", "test")
	globalContext := cli.NewContext(nil, globalSet, nil)

	stringSlice := cli.StringSlice{}
	stringSlice.Set("tcp://localhost:1883")
	flag := cli.StringSliceFlag{Name: "endpoint", Value: &stringSlice}
	flag.Apply(globalSet)
	ctx := cli.NewContext(nil, globalSet, globalContext)

	// load context to the config
	conf := arc_config.Config{}
	err := conf.Load(ctx)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
