package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type Server struct {
	tlsServerCert string
	tlsServerKey  string
	tlsConfig     *tls.Config
	bindAdress    string
	server        *http.Server
}

func NewSever(tlsServerCert string, tlsServerKey string, bindAdress string, router *mux.Router) *Server {
	var tlsConfig *tls.Config
	cer, err := tls.LoadX509KeyPair(tlsServerCert, tlsServerKey)
	if err == nil {
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cer},
			ClientAuth:   tls.VerifyClientCertIfGiven,
		}
	}

	// create server
	svr := &Server{
		tlsServerCert: tlsServerCert,
		tlsServerKey:  tlsServerKey,
		bindAdress:    bindAdress,
		tlsConfig:     tlsConfig,
		server: &http.Server{
			Addr:    bindAdress,
			Handler: router,
		},
	}

	// check tls configuration and add to the server
	if tlsConfig != nil {
		svr.server.TLSConfig = tlsConfig
	}

	return svr
}

func (s *Server) run() {
	ln, err := net.Listen("tcp", s.bindAdress)
	if err != nil {
		log.Fatalf("Failed to create listener: %s", err)
	}
	defer ln.Close()

	if s.tlsConfig == nil {
		log.Infof("Listening on %q... for incoming connections", s.bindAdress)
		err = s.server.Serve(ln)
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		log.Infof("Listening on %q for incoming TLS connections...", s.bindAdress)
		err = s.server.ServeTLS(ln, "", "")
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
