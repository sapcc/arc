package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"html/template"
	"io/ioutil"
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

func getAllBuilds() *[]string {
	var fileNames []string
	builds, _ := ioutil.ReadDir(buildsRootPath)
	for _, f := range builds {
		fileNames = append(fileNames, f.Name())
	}

	if len(fileNames) == 0 {
		fileNames = append(fileNames, "No files found")
	}

	return &fileNames
}
