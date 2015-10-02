package main

//to re generate the asset you need to install esc once: go get -u github.com/mjibson/esc
//go:generate esc -o assets.go static/

import (
	"net/http"

	"github.com/gorilla/mux"

	"gitHub.***REMOVED***/monsoon/arc/update-server/storage"
)

func newRouter(storageType storage.StorageType) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("POST").
		Path("/updates").
		Name("Get available updates").
		Handler(http.HandlerFunc(serveAvailableUpdates))

	if storageType == storage.Local {
		router.
			Methods("POST").
			Path("/upload").
			Name("Upload files").
			Handler(http.HandlerFunc(uploadHandler))

		router.
			Methods("GET").
			PathPrefix("/builds/").
			Name("Serve local build files").
			Handler(http.StripPrefix("/builds/", http.FileServer(http.Dir(st.GetStoragePath()))))
	}

	if storageType == storage.Swift {
		router.
			Methods("GET").
			PathPrefix("/builds/").
			Name("Serve build files from swift").
			Handler(http.StripPrefix("/builds/", http.HandlerFunc(serveSwiftBuilds)))
	}

	router.
		Methods("GET").
		PathPrefix("/static/").
		Name("Serve static files").
		Handler(http.FileServer(FS(false)))

	router.
		Methods("GET").
		Path("/healthcheck").
		Name("Healthcheck").
		Handler(http.HandlerFunc(serveVersion))

	router.
		Methods("GET").
		Path("/readiness").
		Name("Readiness").
		Handler(http.HandlerFunc(serveReadiness))

	router.
		Methods("GET").
		Path("/").
		Name("Serve templates").
		Handler(http.HandlerFunc(serveTemplate))

	return router
}
