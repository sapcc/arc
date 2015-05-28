package main

import (
	"encoding/json"
	"gitHub.***REMOVED***/monsoon/arc/update-server/updates"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"io/ioutil"
	
	"fmt"
)

const StaticRootPath = "/Users/userID/go/src/gitHub.***REMOVED***/monsoon/arc/update-server/static"
const TemplateRootPath = "/Users/userID/go/src/gitHub.***REMOVED***/monsoon/arc/update-server/templates"
const BuildRelativeUrl = "/static/builds/"

type Build struct {
	Files []string
}

func main() {
	// api
	http.HandleFunc("/updates", availableUpdates)
	// serve static files
	fs := http.FileServer(http.Dir(StaticRootPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// serve template
	http.HandleFunc("/", serveTemplate)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func availableUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		update := updates.New(r, StaticRootPath, BuildRelativeUrl)
		if update == nil {
			w.WriteHeader(204)
			return
		}

		if err := json.NewEncoder(w).Encode(update); err != nil {
			log.Println(err.Error())
		}
	} else {
		http.NotFound(w, r)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join(TemplateRootPath, "layout.html")
	fp := path.Join(TemplateRootPath, r.URL.Path)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(err.Error())
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		log.Println("Request is a directory")
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	builds := Build{
		Files: *getAllBuilds(),
	}

	if err := tmpl.ExecuteTemplate(w, "layout", builds); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

// private

func getAllBuilds() *[]string{
	var fileNames []string
	builds, _ := ioutil.ReadDir(fmt.Sprint(StaticRootPath, "/builds"))
	for _, f := range builds {	
		fileNames = append(fileNames, f.Name())
	}
	return &fileNames
}