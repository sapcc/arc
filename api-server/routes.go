package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var routesDefinition = routes{
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
	route{
		"Fact",
		"GET",
		"/agents/{agentId}/facts/{factId}",
		serveFact,
	},
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routesDefinition {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(r.HandlerFunc)
	}

	return router
}
