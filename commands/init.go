package commands

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"gitHub.***REMOVED***/monsoon/arc/service"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

var configTemplate = template.Must(template.New("config").Parse(`{{if .Transport }}transport: {{ .Transport }}
{{end}}{{if .Endpoint }}endpoint: {{ .Endpoint }}
{{end}}tls-client-cert: {{ .Cert }}
tls-client-key: {{ .Key }}
{{if .Ca }}tls-ca-cert: {{ .Ca }}
{{end}}{{if .UpdateUri}}update-uri: {{ .UpdateUri }}
{{end}}{{if .UpdateInterval}}update-interval: {{ .UpdateInterval }}
{{end}}`))

func Init(c *cli.Context, appName string) (int, error) {
	keySize := 2048
	dir := c.String("install-dir")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return 1, fmt.Errorf("Failed to create %s", dir)
		}
	}

	if c.IsSet("registration-url") {
		log.Infof("Generating %d bit private key", keySize)
		key, err := rsa.GenerateKey(rand.Reader, keySize)
		if err != nil {
			return 1, fmt.Errorf("Can't generate a private key")
		}

		keyfile, err := os.OpenFile(path.Join(dir, "cert.key"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return 1, fmt.Errorf("Failed to create %s: %v", path.Join(dir, "cert.key"), err)
		}
		if err := pem.Encode(keyfile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
			return 1, fmt.Errorf("Failed to write private key to file: %v", err)
		}
		keyfile.Close()

		csrTemplate := x509.CertificateRequest{SignatureAlgorithm: x509.SHA256WithRSA}
		log.Info("Creating signing request")
		csrData, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, key)
		if err != nil {
			return 1, fmt.Errorf("Failed to generate csr: %v", err)
		}
		var csr bytes.Buffer
		if err := pem.Encode(&csr, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrData}); err != nil {
			return 1, fmt.Errorf("Failed to pem encode certificate request")
		}

		log.Info("Requesting certificate from registration endpoint")
		client := http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return errors.New("stopped after 10 redirects")
				}
				req.Header["Accept"] = via[0].Header["Accept"]
				req.Header["User-Agent"] = via[0].Header["User-Agent"]
				return nil
			},
		}
		req, _ := http.NewRequest("POST", c.String("registration-url"), &csr)
		req.Header["Accept"] = []string{"application/json"}
		req.Header["User-Agent"] = []string{appName + " " + version.String()}
		resp, err := client.Do(req)
		if err != nil {
			return 1, fmt.Errorf("http request failed: %v", err)
		}
		defer resp.Body.Close()
		switch {
		case resp.StatusCode == 403:
			return 1, fmt.Errorf("Invalid registration token given.")
		case resp.StatusCode != 200:
			return 1, fmt.Errorf("Unknown error while fetching certificate: %s", resp.Status)
		}
		var certs struct {
			Ca          string
			Certificate string
		}
		err = json.NewDecoder(resp.Body).Decode(&certs)
		if err != nil {
			return 1, fmt.Errorf("Failed to parse json reponse: %v", err)
		}

		err = ioutil.WriteFile(path.Join(dir, "cert.pem"), []byte(certs.Certificate), 0644)
		if err != nil {
			return 1, fmt.Errorf("Failed to write certificate to disk: %v", err)
		}
		err = ioutil.WriteFile(path.Join(dir, "ca.pem"), []byte(certs.Ca), 0644)
		if err != nil {
			return 1, fmt.Errorf("Failed to write CA certificate to disk: %v", err)
		}

		log.Info("Retrieved and stored certificate")

		cfgFile, err := os.OpenFile(path.Join(dir, "arc.cfg"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return 1, fmt.Errorf("Failed to create %s: %v", path.Join(dir, "arc.cfg"), err)
		}

		templateVars := map[string]string{
			"Cert":      path.Join(dir, "cert.pem"),
			"Key":       path.Join(dir, "cert.key"),
			"Endpoint":  strings.Join(c.StringSlice("endpoint"), ","),
			"Transport": c.String("transport"),
		}
		err = configTemplate.Execute(cfgFile, templateVars)
		if err != nil {
			return 1, fmt.Errorf("Failed to write config file: %v", err)
		}

	}

	if err := service.New(dir).Install(); err != nil {
		return 1, fmt.Errorf("Failed to install service: %s", err)
	}
	return 0, nil
}
