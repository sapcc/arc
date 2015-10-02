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
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
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

	buildFilesNames, err := st.GetAllUpdates()
	if err != nil {
		log.Errorf("Error getting the build file names. Got %q", err)
		http.Error(w, http.StatusText(500), 500)
	}

	// get build infos
	data := tmplData{
		AppName:    appName,
		AppVersion: version.String(),
		Files:      *buildFilesNames,
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
