package commands

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/pborman/uuid"

	"gitHub.***REMOVED***/monsoon/arc/service"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

var configTemplate = template.Must(template.New("config").Parse(`{{if .Transport }}transport: {{ .Transport }}{{.Eol}}{{end}}{{if .Endpoint }}endpoint: {{ .Endpoint }}{{.Eol}}{{end}}tls-client-cert: {{ .Cert }}{{.Eol}}tls-client-key: {{ .Key }}{{.Eol}}{{if .Ca }}tls-ca-cert: {{ .Ca }}{{.Eol}}{{end}}{{if .UpdateUri}}update-uri: {{ .UpdateUri }}{{.Eol}}{{end}}{{if .UpdateInterval}}update-interval: {{ .UpdateInterval }}{{.Eol}}{{end}}`))

// Init install an arc agent/node
func Init(c *cli.Context, appName string) (int, error) {
	keySize := 2048
	dir := c.String("install-dir")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); /* #nosec */ err != nil {
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
		if err = pem.Encode(keyfile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
			return 1, fmt.Errorf("Failed to write private key to file: %v", err)
		}
		if err := keyfile.Close(); err != nil {
			return 1, fmt.Errorf("Failed to close handle to private key file: %s", err)
		}

		cn := discoverIdentity(c)
		// create csr template
		csrTemplate := x509.CertificateRequest{
			SignatureAlgorithm: x509.SHA256WithRSA,
			Subject: pkix.Name{
				CommonName: cn,
			},
		}
		log.Infof("Creating signing request for identity %#v", cn)
		csrData, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, key)
		if err != nil {
			return 1, fmt.Errorf("Failed to generate csr: %v", err)
		}
		var csr bytes.Buffer
		if err = pem.Encode(&csr, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrData}); err != nil {
			return 1, fmt.Errorf("Failed to PEM encode certificate request")
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
		req, err := http.NewRequest("POST", c.String("registration-url"), &csr)
		if err != nil {
			return 1, fmt.Errorf("Failed to create registration http request: %s", err)
		}
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

		cfgFile, err := os.OpenFile(path.Join(dir, "arc.cfg"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
		if err != nil {
			return 1, fmt.Errorf("Failed to create %s: %v", path.Join(dir, "arc.cfg"), err)
		}

		updateInterval := ""
		if c.IsSet("update-interval") {
			updateInterval = fmt.Sprintf("%d", c.Int("update-interval"))
		}

		eol := "\n"
		if runtime.GOOS == "windows" {
			eol = "\r\n"
		}

		templateVars := map[string]string{
			"Cert":           path.Join(dir, "cert.pem"),
			"Key":            path.Join(dir, "cert.key"),
			"Ca":             path.Join(dir, "ca.pem"),
			"Endpoint":       strings.Join(c.StringSlice("endpoint"), ","),
			"UpdateUri":      c.String("update-uri"),
			"UpdateInterval": updateInterval,
			"Transport":      c.String("transport"),
			"Eol":            eol,
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

func discoverIdentity(c *cli.Context) string {
	if c.String("common-name") != "" {
		return c.String("common-name")
	} else if instID := instanceID(); instID != "" {
		return instID
	} else if name, err := os.Hostname(); err == nil {
		return name
	}
	return uuid.New()
}

type metaDataID struct {
	UUID string `json:"uuid"`
}

var metadataURL = "http://169.254.169.254/openstack/latest/meta_data.json"

// InstanceID returns the instance id from the metadata
func instanceID() string {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	r, err := client.Get(metadataURL)
	if err != nil {
		log.Warnf(fmt.Sprint("Error requesting metadata. ", err.Error()))
		return ""
	}
	defer r.Body.Close()

	var metadata = new(metaDataID)
	err = json.NewDecoder(r.Body).Decode(metadata)
	if err != nil {
		log.Warnf(fmt.Sprint("Error parsing metadata. ", err.Error()))
		return ""
	}

	return metadata.UUID
}
