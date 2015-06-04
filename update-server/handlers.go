package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/arc/update-server/updates"
	"gitHub.***REMOVED***/monsoon/arc/version"
	"net/http"
	"strings"
)

func serveAvailableUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		update, err := updates.New(r, buildsRootPath)
		if err == updates.ArgumentError {
			log.Errorf(err.Error())
			http.Error(w, http.StatusText(500), 500)
			return
		} else if err != nil {
			log.Errorf(err.Error())
			http.Error(w, http.StatusText(400), 400)
			return
		}
		if update == nil {
			w.WriteHeader(204)
			return
		}

		if err := json.NewEncoder(w).Encode(update); err != nil {
			log.Errorf(err.Error())
		}
	} else {
		http.NotFound(w, r)
	}
}

type tmplData struct {
	AppName    string
	AppVersion string
	Files      []string
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// get the page name and the associated template
	name := strings.Replace(r.URL.Path[1:], ".html", "", 1)

	// root path redirect to home page
	if len(name) == 0 {
		name = "home"
	}

	// get the template defined by the url
	tmpl, ok := templates[name]
	if !ok {
		log.Errorf("The template %s does not exist.", name)
		http.NotFound(w, r)
		return
	}

	// get build files
	data := tmplData{
		AppName:    appName,
		AppVersion: version.String(),
		Files:      *getAllBuilds(),
	}

	// render template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Errorf("Error executing template. Got %q", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
