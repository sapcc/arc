package models

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"		
)

type Agent struct {
	AgentID			string	 		`json:"agent_id"`
	CreatedAt   time.Time		`json:"created_at"`
	UpdatedAt   time.Time 	`json:"updated_at"`
}

type Agents []Agent

func GetAgents(db *sql.DB) (*Agents, error) {
	if db == nil {
		return nil, errors.New("Db is nil")
	}
	
	agents := make(Agents,0)
	rows, err := db.Query(ownDb.GetAgentsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agent Agent
	for rows.Next() {
		err = rows.Scan(&agent.AgentID,&agent.CreatedAt,&agent.UpdatedAt)
		if err != nil {
			log.Errorf("Error scaning agent results. Got ", err.Error())
			continue
		}
		agents = append(agents, agent)
	}

	return &agents, nil
}

func GetAgent(db *sql.DB, agent_id string) (*Agent, error) {
	if db == nil {
		return nil, errors.New("Db is nil")
	}
	
	var agent Agent
	err := db.QueryRow(ownDb.GetAgentQuery, agent_id).Scan(&agent.AgentID, &agent.CreatedAt, &agent.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}