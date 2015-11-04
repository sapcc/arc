package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/filter"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"
)

var FilterError = fmt.Errorf("Filter query has a syntax error.")
var RegistrationExistsError = fmt.Errorf("Registration message already handeled.")

type Agent struct {
	AgentID      string    `json:"agent_id"`
	Project      string    `json:"project"`
	Organization string    `json:"organization"`
	Facts        string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedWith  string    `json:"updated_with"`
	UpdatedBy    string    `json:"updated_by"`
}

type Agents []Agent

func (agents *Agents) Get(db *sql.DB, filterQuery string) error {
	// select the query
	sqlQuery, err := buildAgentsQuery("", filterQuery)
	if err != nil {
		return err
	}

	return agents.getAllAgents(db, sqlQuery)
}

func (agents *Agents) GetAuthorized(db *sql.DB, filterQuery string, authorization *auth.Authorization) error {
	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// select the query
	sqlQuery, err := buildAgentsQuery(authorization.ProjectId, filterQuery)
	if err != nil {
		return err
	}

	return agents.getAllAgents(db, sqlQuery)
}

func (agent *Agent) Get(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	err := db.QueryRow(ownDb.GetAgentQuery, agent.AgentID).Scan(&agent.AgentID, &agent.Project, &agent.Organization, &agent.Facts, &agent.CreatedAt, &agent.UpdatedAt, &agent.UpdatedWith, &agent.UpdatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (agent *Agent) GetAuthorized(db Db, authorization *auth.Authorization) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// get the agent
	err = agent.Get(db)
	if err != nil {
		return err
	}

	// check project
	if agent.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	return nil
}

func (agent *Agent) DeleteAuthorized(db Db, authorization *auth.Authorization) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// get the agent
	err = agent.Get(db)
	if err != nil {
		return err
	}

	// check project
	if agent.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	res, err := db.Exec(ownDb.DeleteAgentQuery, agent.AgentID)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("Agent with id %q is removed. %v row(s) where updated.", agent.AgentID, affect)

	return nil
}

func (agent *Agent) Save(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	var lastInsertId string
	if err := db.QueryRow(ownDb.InsertAgentQuery, agent.AgentID, agent.Project, agent.Organization, agent.Facts, agent.CreatedAt, agent.UpdatedAt, agent.UpdatedWith, agent.UpdatedBy).Scan(&lastInsertId); err != nil {
		return err
	}

	log.Infof("New agent with id %q and registration id %q was saved.", agent.AgentID, agent.UpdatedWith)

	return nil
}

func (agent *Agent) Update(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	res, err := db.Exec(ownDb.UpdateAgent, agent.AgentID, agent.Project, agent.Organization, agent.Facts, agent.UpdatedAt, agent.UpdatedWith, agent.UpdatedBy)
	if err != nil {
		return err
	}

	log.Infof("Agent with id %q and registration id %q was updated", agent.AgentID, agent.UpdatedWith)

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("%v row(s) where updated for agent id %q and registratrion reply id %q", affect, agent.AgentID, agent.UpdatedWith)

	return nil
}

func (agent *Agent) FromRegistration(reg *arc.Registration, agentId string) {
	if reg == nil {
		return
	}
	agent.AgentID = reg.Sender
	agent.Project = reg.Project
	agent.Organization = reg.Organization
	agent.Facts = reg.Payload
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()
	agent.UpdatedWith = reg.RegistrationID
	agent.UpdatedBy = agentId
	return
}

func ProcessRegistration(db *sql.DB, reg *arc.Registration, agentId string, concurrencySafe bool) (err error) {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// cast registration to agent
	agent := Agent{}
	agent.FromRegistration(reg, agentId)

	// should check concurrency
	if concurrencySafe {
		safe, err := IsConcurrencySafe(db, agent.UpdatedWith, agentId)
		if err != nil {
			return err
		}
		if safe {
			return processRegistration(db, &agent)
		} else {
			return RegistrationExistsError
		}
	} else {
		return processRegistration(db, &agent)
	}

	return
}

// private

func (agents *Agents) getAllAgents(db *sql.DB, query string) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	*agents = make(Agents, 0)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var agent Agent
	for rows.Next() {
		err = rows.Scan(&agent.AgentID, &agent.Project, &agent.Organization, &agent.Facts, &agent.CreatedAt, &agent.UpdatedAt, &agent.UpdatedWith, &agent.UpdatedBy)
		if err != nil {
			log.Errorf("Error scaning agent results. Got ", err.Error())
			continue
		}
		*agents = append(*agents, agent)
	}

	rows.Close()
	return nil
}

func buildAgentsQuery(authProjectId, filterParam string) (string, error) {
	var err error
	resultQuery := fmt.Sprintf(ownDb.GetAgentsQuery, "")
	authQuery := ""
	filterQuery := ""

	if authProjectId != "" {
		authQuery = fmt.Sprintf(`project='%s'`, authProjectId)
	}

	if filterParam != "" {
		// query string to sql query
		filterQuery, err = filter.Postgresql(filterParam)
		if err != nil {
			return "", FilterError
		}
	}

	log.Info("****")
	log.Info(authProjectId)
	log.Info(filterParam)
	log.Info(filterQuery)
	log.Info("****")

	if authQuery != "" {
		resultQuery = fmt.Sprintf(ownDb.GetAgentsQuery, fmt.Sprint("WHERE ", authQuery))
		if filterQuery != "" {
			resultQuery = fmt.Sprintf(ownDb.GetAgentsQuery, fmt.Sprintf(`WHERE %s AND (%s)`, authQuery, filterQuery))
		}
	} else {
		if filterQuery != "" {
			resultQuery = fmt.Sprintf(ownDb.GetAgentsQuery, fmt.Sprint("WHERE ", filterQuery))
		}
	}

	log.Info("####")
	log.Info(resultQuery)
	log.Info("####")

	return resultQuery, nil
}

func processRegistration(db *sql.DB, agent *Agent) (err error) {
	// create transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// check if the agent already exists
	checkAgent := Agent{AgentID: agent.AgentID}
	existAgentError := checkAgent.Get(tx)
	if existAgentError == sql.ErrNoRows { // agent not found, new agent entry
		if err = agent.Save(tx); err != nil {
			return err
		}
	} else if existAgentError != nil { // something wrong happned
		return existAgentError
	} else { // update the agent
		if err = agent.Update(tx); err != nil {
			return err
		}
	}

	return
}
