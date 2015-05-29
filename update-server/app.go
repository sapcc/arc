package main

import (
	"encoding/json"
	"gitHub.***REMOVED***/monsoon/arc/update-server/updates"
	"html/template"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
)

const BuildsRootPath = "/Users/userID/go/src/gitHub.***REMOVED***/monsoon/arc/update-server/builds"
const BuildRelativeUrl = "/builds/"

type Build struct {
	Files []string
}

func main() {
	// api
	http.HandleFunc("/updates", availableUpdates)
	// serve build files
	fs := http.FileServer(http.Dir(BuildsRootPath))
	http.Handle("/builds/", http.StripPrefix("/builds/", fs))
	
	// serve static files
	http.Handle("/static/", http.FileServer(FS(false)))
	
	// serve template
	http.HandleFunc("/", serveTemplate)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func availableUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		update := updates.New(r, BuildsRootPath, BuildRelativeUrl)
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
	templatesPath := "/static/templates"
	
	// layout
	lp_s, err := FSString(false, fmt.Sprint(templatesPath, "/layout.html"))
	if err != nil {
		log.Printf("Error http.Filesystem for the embedded assets. Got %q", err.Error())
		http.NotFound(w, r)
		return
	}

	// page
	fp_s, err := FSString(false, fmt.Sprint(templatesPath, r.URL.Path))
	if err != nil {
		log.Printf("Error http.Filesystem for the embedded assets. Got %q", err.Error())
		http.NotFound(w, r)
		return
	}
	
	// parse layout
	tmpl, err := template.New("layout").Parse(lp_s)
	if err != nil { 
		log.Printf("Error parsing layout. Got %q", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	
	// parse page
	tmpl, err = tmpl.New("page").Parse(fp_s)
	if err != nil { 
		log.Printf("Error parsing page. Got %q", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	builds := Build{
		Files: *getAllBuilds(),
	}

	if err := tmpl.ExecuteTemplate(w, "layout", builds); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

// private

func getAllBuilds() *[]string{
	var fileNames []string
	builds, _ := ioutil.ReadDir(BuildsRootPath)
	for _, f := range builds {	
		fileNames = append(fileNames, f.Name())
	}
	return &fileNames
}