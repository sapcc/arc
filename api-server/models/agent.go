package models

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"
)

type Db interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Agent struct {
	AgentID      string    `json:"agent_id"`
	Project      string    `json:"project"`
	Organization string    `json:"organization"`
	Facts        string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Agents []Agent

func (agents *Agents) Get(db *sql.DB) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	*agents = make(Agents, 0)
	rows, err := db.Query(ownDb.GetAgentsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var agent Agent
	for rows.Next() {
		err = rows.Scan(&agent.AgentID, &agent.Project, &agent.Organization, &agent.Facts, &agent.CreatedAt, &agent.UpdatedAt)
		if err != nil {
			log.Errorf("Error scaning agent results. Got ", err.Error())
			continue
		}
		*agents = append(*agents, agent)
	}

	rows.Close()
	return nil
}

func (agent *Agent) Get(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	err := db.QueryRow(ownDb.GetAgentQuery, agent.AgentID).Scan(&agent.AgentID, &agent.Project, &agent.Organization, &agent.Facts, &agent.CreatedAt, &agent.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (agent *Agent) Save(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	var lastInsertId string
	if err := db.QueryRow(ownDb.InsertAgentQuery, agent.AgentID, agent.Project, agent.Organization, agent.Facts, agent.CreatedAt, agent.UpdatedAt).Scan(&lastInsertId); err != nil {
		return err
	}

	log.Infof("New agent with id %q saved.", agent.AgentID)

	return nil
}

func (agent *Agent) Update(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	res, err := db.Exec(ownDb.UpdateAgent, agent.AgentID, agent.Project, agent.Organization, agent.Facts, agent.UpdatedAt)
	if err != nil {
		return err
	}

	log.Infof("Agent with id %q updated.", agent.AgentID)

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("%v rows where updated agent id %q", affect, agent.AgentID)

	return nil
}

func (agent *Agent) FromRegistration(reg *arc.Registration) {
	if reg == nil {
		return
	}
	agent.AgentID = reg.Sender
	agent.Project = reg.Project
	agent.Organization = reg.Organization
	agent.Facts = reg.Payload
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()
	return
}

func (agent *Agent) ProcessRegistration(db *sql.DB, reg *arc.Registration) (err error) {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// cast to agent
	agent.FromRegistration(reg)

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

	checkAgent := Agent{AgentID: agent.AgentID}
	err = checkAgent.Get(tx)
	if err == sql.ErrNoRows { // fact not found
		if err = agent.Save(tx); err != nil {
			return err
		}
	} else if err != nil { // something wrong happned
		return
	} else {
		if err = agent.Update(tx); err != nil {
			return
		}
	}

	return
}
