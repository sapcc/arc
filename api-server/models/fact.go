package models

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"		
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type Fact struct {
	Agent			
	Facts  	string 	`json:"facts"`
}

func (fact *Fact) Get(db *sql.DB) error {
	if db == nil {
		return errors.New("Db is nil")
	}
	
	err := db.QueryRow(ownDb.GetFactQuery, fact.AgentID).Scan(&fact.AgentID, &fact.Facts, &fact.CreatedAt, &fact.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (fact *Fact) Update(db *sql.DB, req *arc.Request) (err error) {
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

	// insert or update
	err = tx.QueryRow(ownDb.GetAgentQuery, req.Sender).Scan(&fact.AgentID, &fact.CreatedAt, &fact.UpdatedAt)
	if err == nil {
		log.Infof("Registry for sender %q will be updated.", req.Sender)		
		if _, err = tx.Exec(ownDb.UpdateFact, req.Sender, req.Payload); err != nil {
			return
		}		
	} else if err == sql.ErrNoRows{
		log.Infof("New registry for sender %q will be saved.", req.Sender)
		var lastInsertId string
		if err = tx.QueryRow(ownDb.InsertFactQuery, req.Sender, req.Payload, time.Now(), time.Now()).Scan(&lastInsertId); err != nil {
			return
		}
	} else if err != nil {
		return
	}

	// update object data
	err = tx.QueryRow(ownDb.GetFactQuery, req.Sender).Scan(&fact.AgentID, &fact.Facts, &fact.CreatedAt, &fact.UpdatedAt)
	if err != nil {
		return
	}

	return
}