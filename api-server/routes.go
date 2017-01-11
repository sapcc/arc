package main

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
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
	route{
		"Readiness",
		"GET",
		"/readiness",
		serveReadiness,
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
		"Delete Agent",
		"DELETE",
		"/agents/{agentId}",
		deleteAgent,
	},
	route{
		"Facts",
		"GET",
		"/agents/{agentId}/facts",
		serveFacts,
	},
	route{
		"Agent Tags",
		"GET",
		"/agents/{agentId}/tags",
		serveAgentTags,
	},
	route{
		"Save Agent Tag",
		"POST",
		"/agents/{agentId}/tags",
		saveAgentTags,
	},
	route{
		"Delete an Agent Tag",
		"DELETE",
		"/agents/{agentId}/tags/{value}",
		deleteAgentTag,
	},
}

var v1PkiRoutesDefinition = routes{
	route{
		"Validate token",
		"POST",
		"/pki/sign/{token}",
		signPkiToken,
	},
	route{
		"Create one time token",
		"POST",
		"/pki/token",
		servePkiToken,
	},
}

func newRouter(env string) *mux.Router {
	middlewareChain := alice.New(loggingHandler, combineLogHandler, servedByHandler)
	middlewareChainApiV1 := alice.New(loggingHandler, combineLogHandler, servedByHandler)

	// remove keystone handler for test and test-local
	if env != "test" && env != "test-local" && ks.Endpoint != "" {
		middlewareChainApiV1 = alice.New(loggingHandler, combineLogHandler, servedByHandler, ks.Handler)
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, r := range standardRoutesDefinition {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(middlewareChain.Then(prometheus.InstrumentHandler(r.Name, r.HandlerFunc)))
	}

	// add metrics
	router.
		Methods("GET").
		Path("/metrics").
		Name("Metrics").
		Handler(middlewareChain.Then(prometheus.Handler()))

		// add api/v1 std routes
	v1SubRouter := router.PathPrefix("/api/v1").Subrouter()
	for _, r := range v1RoutesDefinition {
		v1SubRouter.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(middlewareChainApiV1.Then(prometheus.InstrumentHandler(r.Name, r.HandlerFunc)))
	}

	// add pki routes
	if pkiConfig.CFG != nil {
		for _, r := range v1PkiRoutesDefinition {
			v1SubRouter.
				Methods(r.Method).
				Path(r.Pattern).
				Name(r.Name).
				Handler(middlewareChainApiV1.Then(prometheus.InstrumentHandler(r.Name, r.HandlerFunc)))
		}
		log.Infof("PKI profile config found, adding pki routes...")
	}

	return router
}
