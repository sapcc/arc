package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const templatesPath = "/static/templates/"

var pages = []string{"home", "healthcheck"}

func getTemplates() map[string]*template.Template {
	tmplCache := make(map[string]*template.Template)

	// get layout as string
	stringLayout, err := FSString(false, fmt.Sprint(templatesPath, "layout.html"))
	if err != nil {
		log.Errorf("Error http.Filesystem for the embedded assets. Got %q", err)
		return nil
	}

	// loop over the pages, get strings and parse to the templates
	for i := 0; i < len(pages); i++ {
		// get page as string
		stringPage, err := FSString(false, fmt.Sprint(templatesPath, pages[i], ".html"))
		if err != nil {
			log.Errorf("Error http.Filesystem for the embedded assets. Got %q", err)
			continue
		}

		// create a new template
		tmpl, err := template.New("layout").Parse(stringLayout)
		if err != nil {
			log.Errorf("Error parsing layout. Got %q", err)
			continue
		}

		// parse page to the template
		tmpl, err = tmpl.New(pages[i]).Parse(stringPage)
		if err != nil {
			log.Errorf("Error parsing page. Got %q", err)
			continue
		}

		// add template to the template array
		tmplCache[pages[i]] = tmpl
	}

	return tmplCache
}

type Release struct {
	Uid      string `json:"uid"`
	Filename string `json:"filename"`
	App      string `json:"app"`
	Os       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Date     string `json:"date"`
}

func getReleasesConfigFile() ([]byte, error) {
	releasesPath := path.Join(buildsRootPath, "releases.yml")

	// check if path exists
	if _, err := os.Stat(releasesPath); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(releasesPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getBuildExtraInfo() (map[string]Release, error) {
	data, err := getReleasesConfigFile()
	if err != nil {
		return nil, err
	}

	releases := make(map[string]Release)
	err = yaml.Unmarshal([]byte(data), &releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func getAllBuilds() *[]string {
	var fileNames []string
	builds, _ := ioutil.ReadDir(buildsRootPath)
	for _, f := range builds {
		// filter config file
		if strings.ToLower(f.Name()) != "releases.yml" {
			fileNames = append(fileNames, f.Name())	
		}
	}

	if len(fileNames) == 0 {
		fileNames = append(fileNames, "No files found")
	}

	return &fileNames
}
