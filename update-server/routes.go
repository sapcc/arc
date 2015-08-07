package main

//to re generate the asset you need to install esc once: go get -u github.com/mjibson/esc
//go:generate esc -o assets.go static/

import (
	"net/http"
)

func newRouter() http.Handler {
	// api
	http.HandleFunc("/updates", serveAvailableUpdates)

	// serve build files
	fs := http.FileServer(http.Dir(buildsRootPath))
	http.Handle("/builds/", http.StripPrefix("/builds/", fs))

	// serve static files
	http.Handle("/static/", http.FileServer(FS(false)))

	// serve templates
	http.HandleFunc("/", serveTemplate)

	return http.DefaultServeMux
}
