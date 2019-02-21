package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"gitHub.***REMOVED***/monsoon/arc/api-server/auth"
)

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

func indentityAndPolicyHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorization := auth.GetIdentity(r)

		// check authentication
		if err := authorization.CheckIdentity(); err != nil {
			status := http.StatusUnauthorized
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(status)
			apiError := NewApiError(http.StatusText(status), status, "", err, r)
			http.Error(w, apiError.toString(), status)
			return
		}

		// check policy
		if err := authorization.CheckPolicy(*warden); err != nil {
			status := http.StatusUnauthorized
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(status)
			apiError := NewApiError(http.StatusText(status), status, "", err, r)
			http.Error(w, apiError.toString(), status)
			return
		}

		// call next handler
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func servedByHandler(next http.Handler) http.Handler {
	serverName := os.Getenv("HOSTNAME")
	if serverName == "" {
		serverName, _ = os.Hostname() //#nosec
	}
	if serverName == "" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Served-By", serverName)
		next.ServeHTTP(w, r)
	})
}
