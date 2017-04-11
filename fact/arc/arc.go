package arc

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

type Source struct {
	config arc_config.Config
}

func New(config arc_config.Config) Source {
	return Source{config: config}
}

func (h Source) Name() string {
	return "arc"
}

func (h Source) Facts() (map[string]interface{}, error) {
	facts := make(map[string]interface{})

	// build path to the cert.pem
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return facts, fmt.Errorf("Failed to build path to the cert.pem file: %v", err)
	}
	certPath := path.Join(dir, "cert.pem")

	// read cert.pem
	pemCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return facts, fmt.Errorf("Failed to read %s: %v", certPath, err)
	}

	// decode
	certData, _ := pem.Decode(pemCert)
	if certData == nil {
		return facts, errors.New("Failed to parse PEM encoded certificate.")
	}

	// parse
	x509Cert, err := x509.ParseCertificate(certData.Bytes)
	if err != nil {
		return facts, err
	}
	certSubject := x509Cert.Subject

	facts["arc_version"] = version.String()
	if len(certSubject.OrganizationalUnit) > 0 {
		facts["project"] = certSubject.OrganizationalUnit[0]
	}
	if certSubject.CommonName != "" {
		facts["identity"] = certSubject.CommonName
	}
	if len(certSubject.Organization) > 0 {
		facts["organization"] = certSubject.Organization[0]
	}
	return facts, nil
}
