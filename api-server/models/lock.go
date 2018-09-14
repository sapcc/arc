package models

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
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

	return db.QueryRow(ownDb.GetLockQuery, l.LockID).Scan(&l.LockID, &l.AgentID, &l.CreatedAt)
}

func (l *Lock) Save(db Db) error {
	if db == nil {
		return errors.New("Db connection is nil")
	}

	var lastInsertId string
	if err := db.QueryRow(ownDb.InsertLockQuery, l.LockID, l.AgentID, l.CreatedAt).Scan(&lastInsertId); err != nil {
		return err
	}

	log.Debugf("New lock with id %q was saved.", l.LockID)

	return nil
}

func IsConcurrencySafe(db Db, messageId string, agentId string) (bool, error) {
	lock := Lock{LockID: messageId, AgentID: agentId}
	err := lock.Save(db)
	if pg_err, ok := err.(*pq.Error); ok {
		if pg_err.Code == "23505" { // FOREIGN KEY VIOLATION
			return false, nil
		}
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func PruneLocks(db Db) (int64, error) {
	if db == nil {
		return 0, errors.New("Db connection is nil")
	}

	// get log parts to aggregate
	res, err := db.Exec(ownDb.DeleteLocksQuery, 300) // 5 min
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
