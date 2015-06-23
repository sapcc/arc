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


func (agents *Agents) Get(db *sql.DB) error{
	if db == nil {
		return errors.New("Db is nil")
	}
	
	*agents = make(Agents,0)
	rows, err := db.Query(ownDb.GetAgentsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var agent Agent
	for rows.Next() {
		err = rows.Scan(&agent.AgentID,&agent.CreatedAt,&agent.UpdatedAt)
		if err != nil {
			log.Errorf("Error scaning agent results. Got ", err.Error())
			continue
		}
		*agents = append(*agents, agent)
	}

	return nil
}

func (agent *Agent) Get(db *sql.DB, agent_id string) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	err := db.QueryRow(ownDb.GetAgentQuery, agent_id).Scan(&agent.AgentID, &agent.CreatedAt, &agent.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
