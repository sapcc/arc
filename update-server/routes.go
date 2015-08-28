package main

//to re generate the asset you need to install esc once: go get -u github.com/mjibson/esc
//go:generate esc -o assets.go static/

import (
	"net/http"
	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("POST").
		Path("/updates").
		Name("Get available updates").
		Handler( http.HandlerFunc(serveAvailableUpdates) )

	router.
		Methods("GET").
		PathPrefix("/builds/").
		Name("Serve build files").
		Handler( http.StripPrefix("/builds/", http.FileServer(http.Dir(buildsRootPath))) )
		
	router.
		Methods("GET").
		PathPrefix("/static/").
		Name("Serve static files").
		Handler(http.FileServer(FS(false)))

	router.
		Methods("GET").
		Path("/").
		Name("Serve templates").
		Handler( http.HandlerFunc(serveTemplate) )

	return router
}
