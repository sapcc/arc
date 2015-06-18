package models

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"		
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type Agent struct {
	AgentID		string 		`json:"agent_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Agents []Agent

type Fact struct {
	Agent			
	Facts  	string 	`json:"facts"`
}

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
	var agent Agent
	err := db.QueryRow(ownDb.GetAgentQuery, agent_id).Scan(&agent.AgentID, &agent.CreatedAt, &agent.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func GetFact(db *sql.DB, agent_id string) (*Fact, error) {
	var fact Fact
	err := db.QueryRow(ownDb.GetFactQuery, agent_id).Scan(&fact.Facts)
	if err != nil {
		return nil, err
	}
	return &fact, nil
}

func UpdateFact(db *sql.DB, req *arc.Request) (err error) {
	if db == nil {
		return errors.New("Db is nil")
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var agent Agent
	err = tx.QueryRow(ownDb.GetAgentQuery, req.Sender).Scan(&agent.AgentID, &agent.CreatedAt, &agent.UpdatedAt)
	if err == nil {
		log.Infof("Registry for sender %q will be updated.", req.Sender)		
		if _, err = tx.Exec(ownDb.UpdateFact, req.Sender, req.Payload); err != nil {
			return
		}		
	} else {
		log.Infof("New registry for sender %q will be saved.", req.Sender)
		var lastInsertId string
		if err = tx.QueryRow(ownDb.InsertFactQuery, req.Sender, req.Payload, time.Now(), time.Now()).Scan(&lastInsertId); err != nil {
			return
		}
	}

	return
}