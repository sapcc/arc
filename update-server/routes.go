package main

//to re generate the asset you need to install esc once: go get -u github.com/mjibson/esc
//go:generate esc -o assets.go static/

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"gitHub.***REMOVED***/monsoon/arc/update-server/storage"
)

func newRouter(storageType storage.StorageType) *mux.Router {
	middlewareChain := alice.New(loggingHandler, combineLogHandler, servedByHandler)
	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("POST").
		Path("/updates").
		Name("Get available updates").
		Handler(middlewareChain.ThenFunc(http.HandlerFunc(serveAvailableUpdates)))

	router.
		Methods("POST").
		Path("/updates/updates").
		Name("Get available updates. Should fix the problem with arc wrong updates site").
		Handler(middlewareChain.ThenFunc(http.HandlerFunc(serveAvailableUpdates)))

	if storageType == storage.Local {
		router.
			Methods("POST").
			Path("/upload").
			Name("Upload files").
			Handler(middlewareChain.ThenFunc(http.HandlerFunc(uploadHandler)))
	}
	router.
		Methods("GET").
		PathPrefix("/static/").
		Name("Serve static files").
		Handler(middlewareChain.Then(http.FileServer(FS(false))))

	router.
		Methods("GET").
		Path("/healthcheck").
		Name("Healthcheck").
		Handler(middlewareChain.ThenFunc(http.HandlerFunc(serveVersion)))

	router.
		Methods("GET").
		Path("/readiness").
		Name("Readiness").
		Handler(middlewareChain.ThenFunc(http.HandlerFunc(serveReadiness)))

	router.
		Methods("GET").
		Path("/").
		Name("Serve templates").
		Handler(middlewareChain.ThenFunc(http.HandlerFunc(serveTemplate)))

	//
	// builds subrouter
	//
	buildsSubRouter := router.PathPrefix("/builds").Subrouter()
	buildsSubRouter.
		Methods("GET").
		PathPrefix("/latest/").
		Name("Serve latest version").
		Handler(middlewareChain.Then(http.StripPrefix("/builds/latest/", http.HandlerFunc(serveLatestBuild))))

	if storageType == storage.Swift {
		buildsSubRouter.
			Methods("GET").
			PathPrefix("/").
			Name("Serve build files from swift").
			Handler(middlewareChain.Then(http.StripPrefix("/builds/", http.HandlerFunc(serveSwiftBuilds))))
	} else if storageType == storage.Local {
		buildsSubRouter.
			Methods("GET").
			PathPrefix("/").
			Name("Serve local build files").
			Handler(middlewareChain.Then(http.StripPrefix("/builds/", http.FileServer(http.Dir(st.GetStoragePath())))))
	}

	return router
}
