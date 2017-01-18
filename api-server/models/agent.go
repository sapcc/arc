package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"

	"gitHub.***REMOVED***/monsoon/arc/api-server/auth"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/filter"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pagination"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"
)

type RegistrationExistsError struct {
	Msg string
}

func (e RegistrationExistsError) Error() string {
	return e.Msg
}

type FilterError struct {
	Msg string
}

func (e FilterError) Error() string {
	return e.Msg
}

type TagError struct {
	Messages map[string][]string
}

func (e *TagError) Error() string {
	return fmt.Sprintf("Tag key is not alphanumeric ([a-z0-9A-Z]) or key value is empty.")
}

func (e *TagError) MessagesToJson() (string, error) {
	jsonBytes, err := json.Marshal(e.Messages)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"errors":{"tags":%s}}`, string(jsonBytes)), nil
}

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
	sqlQuery, err := buildAgentsQuery(ownDb.GetAgentsQuery, "", filterQuery, nil)
	if err != nil {
		return err
	}

	return agents.getAllAgents(db, sqlQuery, []string{})
}

func (agents *Agents) GetAuthorizedAndShowFacts(db *sql.DB, filterQuery string, authorization *auth.Authorization, showFacts []string, pag *pagination.Pagination) error {
	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// count agents and set total pages
	countQuery, err := buildAgentsQuery(ownDb.CountAgentsQuery, authorization.ProjectId, filterQuery, nil) // no use pagination
	if err != nil {
		return err
	}
	countAgents, err := countAgents(db, countQuery)
	if err != nil {
		return err
	}
	err = pag.SetTotalElements(countAgents)
	if err != nil {
		return err
	}

	// build the query
	sqlQuery, err := buildAgentsQuery(ownDb.GetAgentsQuery, authorization.ProjectId, filterQuery, pag)
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
		return auth.NotAuthorized{Msg: fmt.Sprintf("%s is not project %s", agent.Project, authorization.ProjectId)}
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
		return auth.NotAuthorized{Msg: fmt.Sprintf("%s is not project %s", agent.Project, authorization.ProjectId)}
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
			return RegistrationExistsError{Msg: fmt.Sprint("IsConcurrencySafe returns false by ProcessRegistration. Agent update with: ", agent.UpdatedWith)}
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
		return auth.NotAuthorized{Msg: fmt.Sprintf("%s is not project %s", agent.Project, authorization.ProjectId)}
	}

	res, err := db.Exec(ownDb.AddAgentTag, agent.AgentID, time.Now(), tagKey, tagValue)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Agent with id %q has added tag with key %q value %q. %v row(s) affected", agent.AgentID, tagKey, tagValue, affect)

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
		return auth.NotAuthorized{Msg: fmt.Sprintf("%s is not project %s", agent.Project, authorization.ProjectId)}
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

	// check for aphanumeric keys or empty values
	tagsErrorMessages := make(map[string][]string)
	for k, v := range tags {
		// check for aphanumeric
		match, _ := regexp.MatchString("^\\w+$", k)
		if !match {
			tagsErrorMessages[k] = append(tagsErrorMessages[k], fmt.Sprintf("Tag key %s is not alphanumeric [a-z0-9A-Z].", k))
			continue
		}
		// check for empty values
		if len(v) == 0 || len(v[0]) == 0 {
			tagsErrorMessages[k] = append(tagsErrorMessages[k], fmt.Sprintf("Tag key %s is empty.", k))
			continue
		}
	}

	if len(tagsErrorMessages) > 0 {
		return &TagError{Messages: tagsErrorMessages}
	}

	for k, v := range tags {
		keyTag := k
		valueTag := v[0]

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
func buildAgentsQuery(baseQuery string, authProjectId, filterParam string, pag *pagination.Pagination) (string, error) {
	var err error
	resultQuery := fmt.Sprintf(baseQuery, "", "")
	authQuery := ""
	filterQuery := ""
	paginationQuery := ""

	// check pagination
	if pag != nil {
		paginationQuery = fmt.Sprintf(`OFFSET %v LIMIT %v`, pag.Offset, pag.Limit)
		resultQuery = fmt.Sprintf(baseQuery, "", paginationQuery)
	}

	if authProjectId != "" {
		authQuery = fmt.Sprintf(`project='%s'`, authProjectId)
	}

	if filterParam != "" {
		// query string to sql query
		filterQuery, err = filter.Postgresql(filterParam)
		if err != nil {
			return "", FilterError{Msg: err.Error()}
		}
	}

	if authQuery != "" {
		resultQuery = fmt.Sprintf(baseQuery, fmt.Sprint("WHERE ", authQuery), paginationQuery)
		if filterQuery != "" {
			resultQuery = fmt.Sprintf(baseQuery, fmt.Sprintf(`WHERE %s AND (%s)`, authQuery, filterQuery), paginationQuery)
		}
	} else {
		if filterQuery != "" {
			resultQuery = fmt.Sprintf(baseQuery, fmt.Sprint("WHERE ", filterQuery), paginationQuery)
		}
	}

	return resultQuery, nil
}

func countAgents(db *sql.DB, query string) (int, error) {
	if db == nil {
		return 0, errors.New("Db is nil")
	}

	var countAgents int
	err := db.QueryRow(query).Scan(&countAgents)
	if err != nil {
		return 0, err
	}

	return countAgents, nil
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
