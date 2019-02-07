package main

import (
	"encoding/json"
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
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			if errJSONEncode := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJSONEncode != nil {
				log.Errorf("error encoding JSON: %v", errJSONEncode)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}

		// check policy
		if err := authorization.CheckPolicy(*warden); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			if errJSONEncode := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJSONEncode != nil {
				log.Errorf("error encoding JSON: %v", errJSONEncode)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
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
