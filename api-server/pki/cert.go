package pki

import (
	"crypto/sha1" // #nosec Use of weak cryptographic primitive
	"crypto/x509"
	"fmt"
	"strings"
)

// #nosec Use of weak cryptographic primitive
func certificateFingerprint(cert x509.Certificate) string {
	return strings.Replace(fmt.Sprintf("%x", sha1.Sum(cert.Raw)), " ", ":", -1)
}
