package scan

import "github.com/cloudflare/cf-tls/tls"

// TLSSession contains tests of host TLS Session Resumption via
// Session Tickets and Session IDs
var TLSSession = &Family{
	Description: "Scans host's implementation of TLS session resumption using session tickets/session IDs",
	Scanners: map[string]*Scanner{
		"SessionResume": {
			"Host is able to resume sessions across all addresses",
			sessionResumeScan,
		},
	},
}

// SessionResumeScan tests that host is able to resume sessions across all addresses.
func sessionResumeScan(host string) (grade Grade, output Output, err error) {
	config := defaultTLSConfig(host)
	config.ClientSessionCache = tls.NewLRUClientSessionCache(1)

	conn, err := tls.DialWithDialer(Dialer, Network, host, config)
	if err != nil {
		return
	}
	if err = conn.Close(); err != nil {
		return
	}

	return multiscan(host, func(addrport string) (g Grade, o Output, e error) {
		g = Good
		conn, e1 := tls.DialWithDialer(Dialer, Network, addrport, config)
		if e1 != nil {
			return
		}
		conn.Close()

		o = conn.ConnectionState().DidResume
		if !conn.ConnectionState().DidResume {
			grade = Bad
		}
		return
	})
}
