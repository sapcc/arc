package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/arc/update-server/updates"
	"gitHub.***REMOVED***/monsoon/arc/version"

	"net/http"
	"strings"
	"fmt"
	"path"
	"os"
	"io"
	"errors"
)

func serveAvailableUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
}

type tmplData struct {
	AppName    string
	AppVersion string
	Files      []string
	BuildInfos map[string]Release
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

	// get builds information
	buildsInfo := Releases{}
	buildsInfo.Read()

	// get build infos
	data := tmplData{
		AppName:    appName,
		AppVersion: version.String(),
		Files:      *getAllBuilds(),
		BuildInfos: buildsInfo,
	}

	// render template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Errorf("Error executing template. Got %q", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {	
	fileName := r.URL.Query().Get("fileName")
	if len(fileName) == 0 {
		checkErrAndReturnStatus(w, errors.New("No filename parameter found."), "", http.StatusBadRequest)
		return
	}

	// create the file
	path := path.Join(buildsRootPath, fileName)		
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
	
	fmt.Fprintf(w, "File uploaded successfully\n")	
	return
 }
 
 // private

 func checkErrAndReturnStatus(w http.ResponseWriter, err error, msg string, status int) {
 	if err != nil {
 		log.Errorf("Error, returning status %v. %s %s", status, msg, err.Error())
 		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
 		http.Error(w, http.StatusText(status), status)
 	}
 }