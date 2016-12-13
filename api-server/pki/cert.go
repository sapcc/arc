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

func certificateSubject(cert x509.Certificate) string {
	rdns := cert.Subject.ToRDNSequence()
	fields := make([]string, 0, len(rdns))
	for _, rdn := range rdns {
		atv := rdn[0]
		value, ok := atv.Value.(string)
		if !ok {
			continue
		}
		t := atv.Type

		//TODO: According to https://www.ietf.org/rfc/rfc4514.txt
		//the null charcter needs to be handeled and we don't need to esacpe ' ' and '#' always
		r := strings.NewReplacer(" ", "\\ ", "#", "\\#", "\"", "\\\"", "+", "\\+", ",", "\\,", ";", "\\;", "<", "\\<", ">", "\\>", "\\", "\\\\")
		value = r.Replace(value)

		if len(t) == 4 && t[0] == 2 && t[1] == 5 && t[2] == 4 {
			switch t[3] {
			case 3:
				fields = append(fields, fmt.Sprintf("CN=%s", value))
			//case 5:
			//	n.SerialNumber = value
			case 6:
				fields = append(fields, fmt.Sprintf("C=%s", value))
			case 7:
				fields = append(fields, fmt.Sprintf("L=%s", value))
			case 8:
				fields = append(fields, fmt.Sprintf("ST=%s", value))
			//case 9:
			//	n.StreetAddress = append(n.StreetAddress, value)
			case 10:
				fields = append(fields, fmt.Sprintf("O=%s", value))
			case 11:
				fields = append(fields, fmt.Sprintf("OU=%s", value))
				//case 17:
				//	n.PostalCode = append(n.PostalCode, value)
			}
		}
	}

	return strings.Join(fields, ",")
}
