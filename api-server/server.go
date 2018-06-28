package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type Server struct {
	tlsServerCert   string
	tlsServerKey    string
	tlsConfig       *tls.Config
	httpBindAdress  string
	httpsBindAdress string
	httpServer      *http.Server
	httpsServer     *http.Server
}

// NewSever creates a server struct
func NewSever(tlsServerCert, tlsServerKey, httpBindAdress, httpsBindAdress string, router *mux.Router) *Server {
	var tlsConfig *tls.Config

	if tlsServerCert != "" && tlsServerKey != "" {
		cer, err := tls.LoadX509KeyPair(tlsServerCert, tlsServerKey)
		if err != nil {
			log.Fatalf("Failed to load tls cert/key: %s", err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cer},
			ClientAuth:   tls.VerifyClientCertIfGiven,
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
func (s *Server) run() {
	// http
	httpLn, err := net.Listen("tcp", s.httpBindAdress)
	if err != nil {
		log.Fatalf("Failed to create http listener: %s", err)
	}
	defer httpLn.Close()
	log.Infof("Listening on %q... for incoming connections", s.httpBindAdress)
	err = s.httpServer.Serve(httpLn)
	if err != nil {
		log.Fatalf("Error starting http server: %q", err)
	}

	// https
	if s.tlsConfig != nil {
		s.httpsServer.TLSConfig = s.tlsConfig
		httspLn, err := net.Listen("tcp", s.httpsBindAdress)
		if err != nil {
			log.Fatalf("Failed to create https listener: %s", err)
		}
		log.Infof("Listening on %q for incoming TLS connections...", s.httpsBindAdress)
		err = s.httpsServer.ServeTLS(httspLn, "", "")
		if err != nil {
			log.Fatalf("Error starting https server: %q", err)
		}
	} else {
		log.Infof("TLS configuration nil. No TLS server will be started.")
	}
}

func (s *Server) close() error {
	err := s.httpServer.Close()
	if err != nil {
		return fmt.Errorf("Error closing http server: %s", err)
	}
	err = s.httpsServer.Close()
	return fmt.Errorf("Error closing https server: %s", err)
}

func (s *Server) shutdown() error {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("Error shuting down http server: %s", err)
	}
	err = s.httpsServer.Shutdown(context.Background())
	return fmt.Errorf("Error shuting down https server: %s", err)
}
