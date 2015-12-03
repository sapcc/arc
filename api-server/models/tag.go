package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

type Tag struct {
	AgentID   string    `json:"agent_id"`
	Project   string    `json:"project"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type Tags []Tag

func (tags *Tags) GetByAgentIdAuthorized(db *sql.DB, authorization *auth.Authorization, agent_id string) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// initialize tags
	*tags = make(Tags, 0)

	// get the agent to check the project
	agent := Agent{AgentID: agent_id}
	err = agent.Get(db)
	if err != nil {
		return err
	}

	// check project
	if agent.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	// save the results
	rows, err := db.Query(ownDb.GetTagsByAgentIdQuery, agent_id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		tag := Tag{}
		err = rows.Scan(&tag.AgentID, &tag.Project, &tag.Value, &tag.CreatedAt)
		if err != nil {
			log.Errorf("Error scaning tag results. Got ", err.Error())
			continue
		}

		*tags = append(*tags, tag)
	}
	rows.Close()

	return nil
}

func (tags *Tags) GetByValueAuthorized(db *sql.DB, authorization *auth.Authorization, value string) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	*tags = make(Tags, 0)
	// save the results
	rows, err := db.Query(ownDb.GetTagsByValueQuery, authorization.ProjectId, value)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		tag := Tag{}
		err = rows.Scan(&tag.AgentID, &tag.Project, &tag.Value, &tag.CreatedAt)
		if err != nil {
			log.Errorf("Error scaning tag results. Got ", err.Error())
			continue
		}

		*tags = append(*tags, tag)
	}
	rows.Close()

	return nil
}

func (tag *Tag) GetAuthorized(db *sql.DB, authorization *auth.Authorization) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// get the agent to check the project
	agent := Agent{AgentID: tag.AgentID}
	err = agent.Get(db)
	if err != nil {
		return err
	}

	// check project
	if agent.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	err = db.QueryRow(ownDb.GetTagQuery, tag.AgentID, tag.Value).Scan(&tag.AgentID, &tag.Project, &tag.Value, &tag.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (tag *Tag) SaveAuthorized(db *sql.DB, authorization *auth.Authorization) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// check project
	if tag.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	var lastInsertId string
	if err := db.QueryRow(ownDb.InsertTagQuery, tag.AgentID, tag.Project, tag.Value, tag.CreatedAt).Scan(&lastInsertId); err != nil {
		return err
	}

	log.Infof("New tag with for agent %q and value %q was saved.", tag.AgentID, tag.Value)

	return nil
}

func (tag *Tag) DeleteAuthorized(db *sql.DB, authorization *auth.Authorization) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return err
	}

	// check project
	if tag.Project != authorization.ProjectId {
		return auth.NotAuthorized
	}

	// check if tag exists
	err = tag.GetAuthorized(db, authorization)
	if err != nil {
		return err
	}

	// remove the tag
	res, err := db.Exec(ownDb.DeleteTagQuery, tag.AgentID, tag.Value)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("Tag with agent id %q and value %q is removed. %v row(s) where updated.", tag.AgentID, tag.Value, affect)

	return nil
}

func ProcessAgentTagsData(db *sql.DB, authorization *auth.Authorization, agent Agent, data []byte) error {
	// split data into an array coma separated
	tagsString := strings.Split(string(data), ",")

	for _, element := range tagsString {
		// remove empty spaces
		tagString := strings.TrimSpace(element)

		// create tags
		tag := Tag{AgentID: agent.AgentID, Project: agent.Project, Value: tagString, CreatedAt: time.Now()}
		err := tag.SaveAuthorized(db, authorization)

		if err != nil {
			if pg_err, ok := err.(*pq.Error); ok {
				if pg_err.Code == "23505" { // FOREIGN KEY VIOLATION, do nothing already saved
					continue
				} else if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
