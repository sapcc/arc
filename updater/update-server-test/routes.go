package main

//to re generate the asset you need to install esc once: go get -u github.com/mjibson/esc
//go:generate esc -o assets.go static/

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func NewRouter() *mux.Router {
	middlewareChain := alice.New(loggingHandler, combineLogHandler, servedByHandler)
	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("GET").
		PathPrefix("/updates/").
		Name("Serve static files").
		Handler(middlewareChain.Then(http.FileServer(FS(false))))

	return router
}
