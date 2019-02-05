package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"sync"

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
func NewSever(tlsServerCert, tlsServerKey, tlsServerCA, httpBindAdress, httpsBindAdress string, router *mux.Router) *Server {
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
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cer},
			ClientAuth:   tls.VerifyClientCertIfGiven,
			ClientCAs:    certpool,
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
