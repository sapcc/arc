package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"encoding/json"
	"fmt"
	"gitHub.***REMOVED***/monsoon/arc/update-server/updates"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var BuildsRootPath string
const BuildRelativeUrl = "/builds/"

type Build struct {
	Files []string
}

func main() {
	app := cli.NewApp()

	app.Name = "update-server"
	app.Authors = []cli.Author{
		{
			Name:  "Fabian Ruff",
			Email: "fabian.ruff@sap.com",
		},
		{
			Name:  "Arturo Reuschenbach Puncernau",
			Email: "a.reuschenbach.puncernau@sap.com",
		},
	}
	app.Usage = "web server to to check and deliver from updates"
	app.Version = "0.1.0-dev"
  app.Action = runServer
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log-level,l",
			Usage:  "log level",
			Value:  "info",
		},		
		cli.StringFlag{
			Name:   "builds-path",
			Usage:  "Path to builds in the file system",
		},
		cli.StringFlag{
			Name:   "server-port",
			Usage:  "Update server port",
			Value: "3000",
		},
	}
	
	app.Before = func(c *cli.Context) error {
		lvl, err := log.ParseLevel(c.GlobalString("log-level"))
		if err != nil {
			log.Fatalf("Invalid log level: %s\n", c.GlobalString("log-level"))
			return err
		}
		log.SetLevel(lvl)
		return nil
	}
	
	app.Run(os.Args)
}

// private

func runServer(c *cli.Context) {
	// check mandatory params
	if len(c.GlobalString("builds-path")) == 0 {
		log.Fatalf("Build path is missing. Got %q", c.GlobalString("builds-path"))
		return
	} 
	
	BuildsRootPath = c.GlobalString("builds-path")
	
	// api
	http.HandleFunc("/updates", availableUpdates)

	// serve build files
	fs := http.FileServer(http.Dir(BuildsRootPath))
	http.Handle("/builds/", http.StripPrefix("/builds/", fs))

	// serve static files
	http.Handle("/static/", http.FileServer(FS(false)))

	// serve template
	http.HandleFunc("/", serveTemplate)

	log.Infof("Listening on port %q...", c.GlobalString("server-port"))
	http.ListenAndServe( fmt.Sprint(":", c.GlobalString("server-port")), nil)
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
			log.Errorf(err.Error())
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
		log.Errorf("Error http.Filesystem for the embedded assets. Got %q", err.Error())
		http.NotFound(w, r)
		return
	}

	// page
	fp_s, err := FSString(false, fmt.Sprint(templatesPath, r.URL.Path))
	if err != nil {
		log.Errorf("Error http.Filesystem for the embedded assets. Got %q", err.Error())
		http.NotFound(w, r)
		return
	}

	// parse layout
	tmpl, err := template.New("layout").Parse(lp_s)
	if err != nil {
		log.Errorf("Error parsing layout. Got %q", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// parse page
	tmpl, err = tmpl.New("page").Parse(fp_s)
	if err != nil {
		log.Errorf("Error parsing page. Got %q", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	builds := Build{
		Files: *getAllBuilds(),
	}

	if err := tmpl.ExecuteTemplate(w, "layout", builds); err != nil {
		log.Errorf(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func getAllBuilds() *[]string {
	var fileNames []string
	builds, _ := ioutil.ReadDir(BuildsRootPath)
	for _, f := range builds {
		fileNames = append(fileNames, f.Name())
	}
	return &fileNames
}
