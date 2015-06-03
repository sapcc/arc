package main

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
