package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/inconshreveable/go-update/check"
	"gitHub.***REMOVED***/monsoon/arc/update-server/storage/helpers"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

func serveAvailableUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	update, err := st.GetAvailableUpdate(r)
	if err == helpers.UpdateArgumentError {
		log.Errorf(err.Error())
		http.Error(w, http.StatusText(400), 400)
		return
	} else if err != nil {
		log.Errorf(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if update == nil {
		w.WriteHeader(204)
		return
	}

	if err := json.NewEncoder(w).Encode(update); err != nil {
		log.Errorf(err.Error())
	}
}

func serveSwiftBuilds(w http.ResponseWriter, r *http.Request) {
	err := st.GetUpdate(r.URL.Path, w)
	if err == helpers.ObjectNotFoundError {
		checkErrAndReturnStatus(w, err, "Error getting swift update. ", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func serveLatestBuild(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	uriSegments := strings.Split(path, "/")
	if len(uriSegments) != 3 {
		checkErrAndReturnStatus(w, fmt.Errorf("Error getting lastest update. Not enough arguments."), "", http.StatusNotFound)
		return
	}

	// save params
	app := uriSegments[0]
	os := uriSegments[1]
	arch := uriSegments[2]

	params := check.Params{AppId: app, Tags: map[string]string{"os": os, "arch": arch}}
	latestUpdate, err := st.GetLastestUpdate(&params)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Error getting latest update. ", http.StatusInternalServerError)
		return
	}

	if latestUpdate == "" {
		checkErrAndReturnStatus(w, fmt.Errorf("Error getting lastest update. No latest update available for this configuration."), "", http.StatusNotFound)
		return
	}

	// set header for the filename
	w.Header().Set("Content-disposition", fmt.Sprint("attachment; filename=", latestUpdate))
	err = st.GetUpdate(latestUpdate, w)
	if err == helpers.ObjectNotFoundError {
		checkErrAndReturnStatus(w, err, "Error serving latest update. ", http.StatusNotFound)
		return
	} else if err != nil {
		checkErrAndReturnStatus(w, err, "Error serving latest update. ", http.StatusInternalServerError)
		return
	}
}

type tmplData struct {
	AppName     string
	AppVersion  string
	LastUpdates []string
	AllUpdates  []string
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

	lastUpdates, allUpdates, err := st.GetWebUpdates()
	if err != nil {
		log.Errorf("Error getting the build file names. Got %q", err)
		http.Error(w, http.StatusText(500), 500)
	}

	// get build infos
	data := tmplData{
		AppName:     appName,
		AppVersion:  version.String(),
		LastUpdates: *lastUpdates,
		AllUpdates:  *allUpdates,
	}

	// render template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Errorf("Error executing template. Got %q", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("filename")
	if len(fileName) == 0 {
		checkErrAndReturnStatus(w, errors.New("No filename parameter found."), "", http.StatusBadRequest)
		return
	}

	// create the file if not exists
	path := path.Join(st.GetStoragePath(), fileName)
	out, err := os.Create(path)
	if err != nil {
		checkErrAndReturnStatus(w, err, "Unable to create the file for writing.", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, r.Body)
	if err != nil {
		checkErrAndReturnStatus(w, err, "", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully.\n")
	return
}

/*
 * Readiness
 */

type Readiness struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

func serveReadiness(w http.ResponseWriter, r *http.Request) {
	if !st.IsConnected() {
		ready := Readiness{
			Status:  http.StatusBadGateway,
			Message: "Storage not reachable",
		}

		// convert struct to json
		body, err := json.Marshal(ready)
		checkErrAndReturnStatus(w, err, "Error encoding Agent to JSON", http.StatusInternalServerError)

		// return the error with json body
		http.Error(w, string(body), ready.Status)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		log.Errorf("Error, returning status %v. %s", ready.Status, ready.Message)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Ready!!!"))
}

/*
 * Healthcheck
 */

func serveVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Arc update-server " + version.String()))
}

// private

func checkErrAndReturnStatus(w http.ResponseWriter, err error, msg string, status int) {
	if err != nil {
		log.Errorf("Error, returning status %v. %s %s", status, msg, err.Error())
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, http.StatusText(status), status)
	}
}
