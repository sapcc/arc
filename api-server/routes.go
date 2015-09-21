package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var standardRoutesDefinition = routes{
	route{
		"Root",
		"GET",
		"/",
		serveVersion,
	},
	route{
		"Healthcheck",
		"GET",
		"/healthcheck",
		serveVersion,
	},
}

var v1RoutesDefinition = routes{
	route{
		"Jobs",
		"GET",
		"/jobs",
		serveJobs,
	},
	route{
		"Execute Job",
		"POST",
		"/jobs",
		executeJob,
	},
	route{
		"Job",
		"GET",
		"/jobs/{jobId}",
		serveJob,
	},
	route{
		"Job logs",
		"GET",
		"/jobs/{jobId}/log",
		serveJobLog,
	},
	route{
		"Agents",
		"GET",
		"/agents",
		serveAgents,
	},
	route{
		"Agent",
		"GET",
		"/agents/{agentId}",
		serveAgent,
	},
	route{
		"Facts",
		"GET",
		"/agents/{agentId}/facts",
		serveFacts,
	},
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range standardRoutesDefinition {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(r.HandlerFunc)
	}
	
	v1SubRouter := router.PathPrefix("/api/v1").Subrouter()
	for _, r := range v1RoutesDefinition {
		v1SubRouter.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(r.HandlerFunc)
	}

	return router
}
