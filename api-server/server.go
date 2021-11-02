package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
)

type Server struct {
	tlsServerCert   string
	tlsServerKey    string
	tlsConfig       *tls.Config
	httpBindAdress  string
	httpsBindAdress string
	httpServer      *http.Server
	httpsServer     *http.Server
	close           func()
}

// NewSever creates a server struct
func NewSever(tlsServerCert, tlsServerKey, tlsServerCA, tlsServerCRL, httpBindAdress, httpsBindAdress string, router *mux.Router) *Server {
	var tlsConfig *tls.Config

	if tlsServerCert != "" && tlsServerKey != "" {
		cer, err := tls.LoadX509KeyPair(tlsServerCert, tlsServerKey)
		if err != nil {
			log.Fatalf("Failed to load tls cert/key: %s", err)
		}
		pemCerts, err := ioutil.ReadFile(filepath.Clean(tlsServerCA))
		if err != nil {
			log.Fatalf("Failed to load CA certificate: %s", err)
		}
		certpool := x509.NewCertPool()
		if !certpool.AppendCertsFromPEM(pemCerts) {
			log.Fatalf("Given CA file does not contain a PEM encoded x509 certificate")
		}
		var verifyFunc func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error
		if tlsServerCRL != "" {
			crlData, err := ioutil.ReadFile(filepath.Clean(tlsServerCRL))
			if err != nil {
				log.Fatalf("Failed to read given crl file %s: %s", tlsServerCRL, err)
			}
			crlList, err := x509.ParseCRL(crlData)
			if err != nil {
				log.Fatalf("Failed to parse crl list: %s", err)
			}
			caCerts, err := loadCertsFromPEM(pemCerts)
			if err != nil {
				log.Fatalf("Failed to load CA certs", err)
			}
			crlVerified := false
			for _, c := range caCerts {
				if c.CheckCRLSignature(crlList) == nil {
					crlVerified = true
				}
			}
			if !crlVerified || crlList.HasExpired(time.Now()) {
				log.Fatal("Crl list couldn't be verified by given CA")
			}
			revokedCerts := make(map[string]struct{}, len(crlList.TBSCertList.RevokedCertificates))
			for _, r := range crlList.TBSCertList.RevokedCertificates {
				revokedCerts[r.SerialNumber.String()] = struct{}{}
			}
			verifyFunc = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				for _, c := range verifiedChains {
					for _, b := range c {
						if _, revoked := revokedCerts[b.SerialNumber.String()]; revoked {
							return fmt.Errorf("Certificate %s is revoked", b.Subject.String())
						}
					}
				}
				return nil
			}

		}

		tlsConfig = &tls.Config{
			Certificates:          []tls.Certificate{cer},
			ClientAuth:            tls.VerifyClientCertIfGiven,
			ClientCAs:             certpool,
			VerifyPeerCertificate: verifyFunc,
		}
	} else {
		log.Infof("No TLS Cert or key given, no TLS configuration will be created.")
	}

	// create server
	svr := &Server{
		tlsServerCert:   tlsServerCert,
		tlsServerKey:    tlsServerKey,
		httpBindAdress:  httpBindAdress,
		httpsBindAdress: httpsBindAdress,
		tlsConfig:       tlsConfig,
		httpServer: &http.Server{
			Handler: router,
		},
		httpsServer: &http.Server{
			Handler: router,
		},
	}

	return svr
}

func loadCertsFromPEM(pemCerts []byte) ([]*x509.Certificate, error) {
	certs := []*x509.Certificate{}
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		certBytes := block.Bytes
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			continue
		}
		certs = append(certs, cert)
	}
	if len(certs) < 1 {
		return nil, errors.New("No certs found")
	}
	return certs, nil

}

// don't return error. Runing in a go rutine
func (s *Server) run() error {

	var runner run.Group
	// http
	httpLn, err := net.Listen("tcp", s.httpBindAdress)
	if err != nil {
		log.Fatalf("Failed to create http listener: %s", err)
	}

	closeCh := make(chan struct{})
	var once sync.Once
	s.close = func() {
		once.Do(func() { close(closeCh) })
	}

	//add stop handler
	runner.Add(func() error {
		<-closeCh
		return nil
	}, func(_ error) {
		s.close()
	})

	//add http listener
	runner.Add(func() error {
		log.Infof("Listening on %q... for incoming connections", s.httpBindAdress)
		return s.httpServer.Serve(httpLn)
	}, func(_ error) {
		if err = httpLn.Close(); err != nil {
			log.Errorf("error closing http listener: %s", err)
		}
	})

	// https
	if s.tlsConfig != nil {
		s.httpsServer.TLSConfig = s.tlsConfig
		httpsLn, err := net.Listen("tcp", s.httpsBindAdress)
		if err != nil {
			log.Fatalf("Failed to create https listener: %s", err)
		}
		// add https listener
		runner.Add(func() error {
			log.Infof("Listening on %q for incoming TLS connections...", s.httpsBindAdress)
			return s.httpsServer.ServeTLS(httpsLn, "", "")
		}, func(_ error) {
			if err = httpLn.Close(); err != nil {
				log.Errorf("error closing http listener: %s", err)
			}
		})
	}

	return runner.Run()
}

func (s *Server) shutdown() error {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("error shuting down http server: %s", err)
	}
	err = s.httpsServer.Shutdown(context.Background())
	return fmt.Errorf("error shuting down https server: %s", err)
}
