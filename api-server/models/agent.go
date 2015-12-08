package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"

	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/filter"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"
)

var FilterError = fmt.Errorf("Filter query has a syntax error.")
var RegistrationExistsError = fmt.Errorf("Registration message already handeled.")

type JSONB map[string]interface{}

type Agent struct {
	AgentID      string    `json:"agent_id"`
	Project      string    `json:"project"`
	Organization string    `json:"organization"`
	Facts        JSONB     `json:"facts,omitempty"`
	Tags         JSONB     `json:"tags,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedWith  string    `json:"updated_with"`
	UpdatedBy    string    `json:"updated_by"`
}

type Agents []Agent

func (agents *Agents) Get(db *sql.DB, filterQuery string) error {
	// build the query just with facts filter
	sqlQuery, err := buildAgentsQuery("", filterQuery)
	if err != nil {
		return err
	}

	return agents.getAllAgents(db, sqlQuery, []string{})
}

func (agents *Agents) GetAuthorizedAndShowFacts(db *sql.DB, filterQuery string, authorization *auth.Authorization, showFacts []string) error {
	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// build the query
	sqlQuery, err := buildAgentsQuery(authorization.ProjectId, filterQuery)
	if err != nil {
		return err
	}

	return agents.getAllAgents(db, sqlQuery, showFacts)
}

func (agent *Agent) Get(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	err := db.QueryRow(ownDb.GetAgentQuery, agent.AgentID).Scan(&agent.AgentID, &agent.Project, &agent.Organization, &agent.Facts, &agent.CreatedAt, &agent.UpdatedAt, &agent.UpdatedWith, &agent.UpdatedBy, &agent.Tags)
	if err != nil {
		return err
	}

	return nil
}

func (agent *Agent) GetAuthorizedAndShowFacts(db Db, authorization *auth.Authorization, showFacts []string) error {
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

	// filter facts
	facts, err := filterFacts(agent.Facts, showFacts)
	if err != nil {
		return err
	}
	agent.Facts = facts

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

	//When creating an agent we assume it is online by default
	//This is a little hacky but the altnatives are also not that good
	agent.Facts["online"] = true

	// transform agents JSONB to JSON string
	facts, err := json.Marshal(agent.Facts)
	if err != nil {
		return err
	}
	tags, err := json.Marshal(agent.Tags)
	if err != nil {
		return err
	}

	// we set default empty JSON if no tags are set
	// avoids an error when using the json_set_key function
	if agent.Tags == nil {
		tags = []uint8("{}")
	}

	var lastInsertId string
	if err := db.QueryRow(ownDb.InsertAgentQuery, agent.AgentID, agent.Project, agent.Organization, string(facts), agent.CreatedAt, agent.UpdatedAt, agent.UpdatedWith, agent.UpdatedBy, string(tags)).Scan(&lastInsertId); err != nil {
		return err
	}

	log.Infof("New agent with id %q and registration id %q was saved.", agent.AgentID, agent.UpdatedWith)

	return nil
}

// registration (facts)

func (agent *Agent) Update(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// transform agents JSONB to JSON string
	facts, err := json.Marshal(agent.Facts)
	if err != nil {
		return err
	}

	res, err := db.Exec(ownDb.UpdateAgentWithRegistration, agent.AgentID, agent.Project, agent.Organization, string(facts), agent.UpdatedAt, agent.UpdatedWith, agent.UpdatedBy)
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

func (agent *Agent) FromRegistration(reg *arc.Registration, agentId string) error {
	if reg == nil {
		return errors.New("Registration is nil")
	}
	agent.AgentID = reg.Sender
	agent.Project = reg.Project
	agent.Organization = reg.Organization
	err := json.Unmarshal([]byte(reg.Payload), &agent.Facts)
	if err != nil {
		return err
	}
	agent.CreatedAt = time.Now() // usefull just when saving, update method will ignore this value
	agent.UpdatedAt = time.Now()
	agent.UpdatedWith = reg.RegistrationID
	agent.UpdatedBy = agentId
	return nil
}

func ProcessRegistration(db *sql.DB, reg *arc.Registration, agentId string, concurrencySafe bool) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	var err error

	// cast registration to agent
	agent := Agent{}
	err = agent.FromRegistration(reg, agentId)
	if err != nil {
		return err
	}

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

	return nil
}

// ****
// tags
// ****

func (agent *Agent) AddTagAuthorized(db Db, authorization *auth.Authorization, tagKey string, tagValue string) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// check if the project is empty
	if agent.Project == "" {
		// get the agent
		err = agent.Get(db)
		if err != nil {
			return err
		}
	}

	// check project
	if agent.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	res, err := db.Exec(ownDb.AddAgentTag, agent.AgentID, time.Now(), tagKey, tagValue)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Agent with id %q has added tag with key $q value $q. %v row(s) affected", agent.AgentID, tagKey, tagValue, affect)

	return nil
}

func (agent *Agent) DeleteTagAuthorized(db Db, authorization *auth.Authorization, tagKey string) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// check if the project is empty
	if agent.Project == "" {
		// get the agent
		err = agent.Get(db)
		if err != nil {
			return err
		}
	}

	// check project
	if agent.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	res, err := db.Exec(ownDb.DeleteAgentTagQuery, agent.AgentID, time.Now(), tagKey)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("Tag %q from Agent id %q is removed. %v row(s) where updated.", tagKey, agent.AgentID, affect)

	return nil
}

func ProcessTags(db *sql.DB, authorization *auth.Authorization, agentId string, tags url.Values) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// get agent
	agent := Agent{AgentID: agentId}
	err := agent.GetAuthorizedAndShowFacts(db, authorization, []string{})
	if err != nil {
		return err
	}

	for k, v := range tags {
		keyTag := k
		valueTag := ""
		if len(v) > 0 {
			valueTag = v[0]
		}
		err := agent.AddTagAuthorized(db, authorization, keyTag, valueTag)
		// if something wrong happens we brake the process
		if err != nil {
			return err
		}
	}

	return nil
}

// private

func (agents *Agents) getAllAgents(db *sql.DB, query string, facts []string) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	*agents = make(Agents, 0)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		agent := Agent{}
		err = rows.Scan(&agent.AgentID, &agent.Project, &agent.Organization, &agent.Facts, &agent.CreatedAt, &agent.UpdatedAt, &agent.UpdatedWith, &agent.UpdatedBy, &agent.Tags)
		if err != nil {
			log.Errorf("Error scaning agent results. Got ", err.Error())
			continue
		}

		// filter facts from the agent
		facts, err := filterFacts(agent.Facts, facts)
		if err != nil {
			return err
		}
		agent.Facts = facts

		*agents = append(*agents, agent)
	}

	rows.Close()
	return nil
}

//
// Key word "all", when no other facts are added in the showFacts array, will show
// all existing facts
//
func filterFacts(facts map[string]interface{}, showFacts []string) (JSONB, error) {
	// check if there is facts to filter out
	if len(facts) == 0 {
		return nil, nil
	}

	target := make(JSONB)

	// check if there is just the "all" key word
	if len(showFacts) == 1 && showFacts[0] == "all" {
		for key, value := range facts {
			target[key] = value
		}
	} else {
		for _, item := range showFacts {
			if val, ok := facts[item]; ok {
				target[item] = val
			}
		}
	}

	if len(target) == 0 {
		return nil, nil
	}

	return target, nil
}

//
// Builds a sql query based on the facts filter and the authorization project id
//
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

func (j JSONB) Value() (interface{}, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	// nothing is set yet in the clumn
	if value == nil {
		return nil
	}

	// convert to json
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}

	return nil
}
