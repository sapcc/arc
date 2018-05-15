package pki

import (
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"strings"
)

func certificateFingerprint(cert x509.Certificate) string {
	return strings.Replace(fmt.Sprintf("%x", sha1.Sum(cert.Raw)), " ", ":", -1)
}
