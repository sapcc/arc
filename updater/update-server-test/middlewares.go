package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
)

//lint:file-ignore U1000 Ignore all unused code, it is just for testing

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Infof("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func combineLogHandler(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func servedByHandler(next http.Handler) http.Handler {
	var err error
	serverName := os.Getenv("HOSTNAME")
	if serverName == "" {
		serverName, err = os.Hostname()
		if err != nil {
			log.Error(err)
		}
	}
	if serverName == "" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Served-By", serverName)
		next.ServeHTTP(w, r)
	})
}
