package models

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"time"
)

type Lock struct {
	LockID    string    `json:"registry_id"`
	AgentID   string    `json:"agent_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *Lock) Get(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	err := db.QueryRow(ownDb.GetLockQuery, l.LockID).Scan(&l.LockID, &l.AgentID, &l.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (l *Lock) Save(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	var lastInsertId string
	if err := db.QueryRow(ownDb.InsertLockQuery, l.LockID, l.AgentID).Scan(&lastInsertId); err != nil {
		return err
	}

	log.Infof("New lock with id %q was saved.", l.LockID)

	return nil
}
