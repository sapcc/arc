package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Agents",
		"GET",
		"/agent",
		serveAgents,
	},
	Route{
		"Agent",
		"GET",
		"/agent/{agentId}",
		serveAgent,
	},
	Route{
		"Facts",
		"GET",
		"/agent/{agentId}/facts",
		serveFacts,
	},
	Route{
		"Fact",
		"GET",
		"/agent/{agentId}/facts/{factId}",
		serveFact,
	},
}

func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
