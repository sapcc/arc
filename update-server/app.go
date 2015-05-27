package main

import (
  "html/template"
  "log"
  "net/http"
  "os"
  "path"
	"encoding/json"
	
	"gitHub.***REMOVED***/monsoon/arc/update-server/updates"
)

func main() {
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))	
  http.HandleFunc("/updates", availableUpdates)	
  http.HandleFunc("/", serveTemplate)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
}

func availableUpdates(w http.ResponseWriter, r *http.Request) {	
	w.Header().Set("Content-Type", "application/json")
	update := updates.New(r)
  if err := json.NewEncoder(w).Encode(update); err != nil {
      log.Println(err.Error())
  }
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {	
  lp := path.Join("templates", "layout.html")
  fp := path.Join("templates", r.URL.Path)

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

  if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
    log.Println(err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
}