package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gitHub.***REMOVED***/monsoon/arc/update-server/updates"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

// TODO: global var templates path
// TODO: testing

var appName = "arc-update-server"
var BuildsRootPath string
var templates map[string]*template.Template
const BuildRelativeUrl = "/builds/"

type TmplData struct {
	AppName string
	AppVersion string
	Files []string
}

func main() {
	app := cli.NewApp()

	app.Name = appName
	app.Version = version.String()
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

	// read and cache templates
	cacheTemplates()

	// api
	http.HandleFunc("/updates", availableUpdates)

	// serve build files
	fs := http.FileServer(http.Dir(BuildsRootPath))
	http.Handle("/builds/", http.StripPrefix("/builds/", fs))

	// serve static files
	http.Handle("/static/", http.FileServer(FS(false)))

	// serve templates
	http.HandleFunc("/", serveTemplate)

	// run server
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

func cacheTemplates() {
	templatesPath := "/static/templates/"
	pages := []string{"home", "healthcheck"}
	
	// init templates
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	
	// get layout as string
	stringLayout, err := FSString(false, fmt.Sprint(templatesPath, "layout.html"))
	if err != nil {
		log.Errorf("Error http.Filesystem for the embedded assets. Got %q", err)
		return
	}	
	
	// loop over the pages, get strings and parse to the templates
	for i := 0; i < len(pages); i++ {
		// get page as string
		stringPage, err := FSString(false, fmt.Sprint(templatesPath, pages[i], ".html"))
		if err != nil {
			log.Errorf("Error http.Filesystem for the embedded assets. Got %q", err)
			return
		}
		
		// create a new template
		tmpl, err := template.New("layout").Parse(stringLayout)
		if err != nil {
			log.Errorf("Error parsing layout. Got %q", err)
			return
		}

		// parse page to the template
		tmpl, err = tmpl.New(pages[i]).Parse(stringPage)
		if err != nil {
			log.Errorf("Error parsing page. Got %q", err)
			return
		}
		
		// add template to the template array
		templates[pages[i]] = tmpl
	}
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
	data := TmplData{
		AppName: appName,
		AppVersion: version.String(),
		Files: *getAllBuilds(),
	}

	// render template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Errorf("Error executing template. Got %q", err)
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
