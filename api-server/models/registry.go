package models

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

type Registry struct {
	RegistryID string `json:"registry_id"`
	AgentID    string `json:"agent_id"`
}

func (reg *Registry) Get(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	err := db.QueryRow(ownDb.GetRegistryQuery, reg.RegistryID).Scan(&reg.RegistryID, &reg.AgentID)
	if err != nil {
		return err
	}

	return nil
}

func (reg *Registry) Save(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	var lastInsertId string
	if err := db.QueryRow(ownDb.InsertRegistryQuery, reg.RegistryID, reg.AgentID).Scan(&lastInsertId); err != nil {
		return err
	}

	log.Infof("New registration with id %q was saved.", reg.RegistryID)

	return nil
}
