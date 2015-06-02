package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"net/http"
)

func serveAgents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agents)
}

func serveAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	agent := getAgent(agentId)
	if agent == nil {
		log.Errorf("Agent with id %q not found.", agentId)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agent)
}

func serveFacts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["agentId"]

	agent := getAgent(agentId)
	if agent == nil {
		log.Errorf("Agent with id %q not found.", agentId)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(agent.Facts)
}

func serveFact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentId := vars["agentId"]
	factId := vars["factId"]

	agent := getAgent(agentId)
	if agent == nil {
		log.Errorf("Agent with id %q not found.", agentId)
		http.NotFound(w, r)
		return
	}

	var fact models.Fact
	for _, f := range agent.Facts {
		if f.Id == factId {
			fact = models.Fact(f)
		}
	}
	if len(fact.Id) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(fact)
}

// private

func getAgent(agentId string) *models.Agent {
	var agent models.Agent
	for _, a := range agents {
		if a.Id == agentId {
			agent = models.Agent(a)
		}
	}
	if len(agent.Id) == 0 {
		return nil
	}
	return &agent
}
