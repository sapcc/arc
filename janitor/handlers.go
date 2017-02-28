package janitor

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"runtime"

	"github.com/bamzi/jobrunner"
	"github.com/cloudflare/cfssl/log"
)

func serveVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(landingPage)
}

func serveJobrunner(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(jobrunner.StatusJson())
		if err != nil {
			logAndReturnHttpError(w, http.StatusInternalServerError, err)
		}
	default:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t := template.New("status_page")
		t, err := t.Parse(statusPage)
		if err != nil {
			logAndReturnHttpError(w, http.StatusInternalServerError, err)
		}
		t.ExecuteTemplate(w, "status_page", jobrunner.StatusPage())
	}
}

func logAndReturnHttpError(w http.ResponseWriter, status int, err error) {
	_, file, line, _ := runtime.Caller(1)
	log.Errorf("Janitor request error. status=%d location=%s:%d error=%v", status, path.Base(file), line, err)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	http.Error(w, err.Error(), status)
}
