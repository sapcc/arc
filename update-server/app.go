package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gitHub.***REMOVED***/monsoon/arc/update-server/updates"
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
			Name:  "log-level,l",
			Usage: "log level",
			Value: "info",
		},
		cli.StringFlag{
			Name:  "path,p",
			Usage: "Directory containig update artifacts",
		},
		cli.StringFlag{
			Name:  "bind-address,b",
			Usage: "listen address for the update server",
			Value: "0.0.0.0:3000",
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
	BuildsRootPath = c.GlobalString("path")
	if BuildsRootPath == "" {
		log.Fatal("No path to update artifacts given.")
	}

	// api
	http.HandleFunc("/updates", availableUpdates)

	// serve build files
	fs := http.FileServer(http.Dir(BuildsRootPath))
	http.Handle("/builds/", http.StripPrefix("/builds/", fs))

	// serve static files
	http.Handle("/static/", http.FileServer(FS(false)))

	// serve template
	http.HandleFunc("/", serveTemplate)

	log.Infof("Listening on %q...", c.GlobalString("bind-address"))
	if err := http.ListenAndServe(c.GlobalString("bind-address"), accessLogger(http.DefaultServeMux)); err != nil {
		log.Fatalf("Failed to bind on %s: %s", c.GlobalString("bind-address"), err)
	}
}

func accessLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
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
