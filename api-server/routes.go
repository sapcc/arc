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

var routesWithoutPrefix = routes{
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

func newRouter(env string) *mux.Router {
	middlewareChain := alice.New(loggingHandler, combineLogHandler, servedByHandler)
	middlewareChainAuth := alice.New(loggingHandler, combineLogHandler, servedByHandler)

	// remove keystone handler for test and test-local
	if env != "test" && env != "test-local" && ks.Endpoint != "" {
		middlewareChainAuth = alice.New(loggingHandler, combineLogHandler, servedByHandler, ks.Handler, indentityAndPolicyHandler)
	} else {
		middlewareChainAuth = alice.New(loggingHandler, combineLogHandler, servedByHandler, indentityAndPolicyHandler)
	}

	router := mux.NewRouter().StrictSlash(true)

	// add routes without prefix
	for _, r := range routesWithoutPrefix {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(middlewareChain.Then(prometheus.InstrumentHandler(r.Name, r.HandlerFunc)))
	}

	// add metrics outside the loop to not be instrumented by prometheus (loop routesWithoutPrefix) and not have keystone validation (loop v1RoutesDefinition)
	router.
		Methods("GET").
		Path("/metrics").
		Name("Metrics").
		Handler(middlewareChain.Then(prometheus.Handler()))

	// add api/v1 subrouter
	v1SubRouter := router.PathPrefix("/api/v1").Subrouter()

	// add api/v1 routes with authentication
	for _, r := range v1RoutesDefinition {
		v1SubRouter.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(middlewareChainAuth.Then(prometheus.InstrumentHandler(r.Name, r.HandlerFunc)))
	}

	// add pki routes
	if pkiEnabled {
		// route with authentication
		v1SubRouter.
			Methods("POST").
			Path("/agents/init").
			Name("Create node bootstrap credentials").
			Handler(middlewareChainAuth.Then(prometheus.InstrumentHandler("Create node bootstrap credentials", http.HandlerFunc(servePkiToken))))

		// route without authentication
		v1SubRouter.
			Methods("POST").
			Path("/agents/init/{token}").
			Name("Create node").
			Handler(middlewareChain.Then(prometheus.InstrumentHandler("Create node", http.HandlerFunc(signPkiToken))))
		v1SubRouter.
			Methods("POST").
			Path("/agents/renew").
			Name("Renew certificate").
			Handler(middlewareChain.Then(prometheus.InstrumentHandler("Renew certificate", http.HandlerFunc(renewPkiCert))))

		log.Infof("PKI profile config found, adding pki routes...")
	}

	return router
}
