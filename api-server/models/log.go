package models

import (
	"database/sql"
	"time"
	"fmt"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

func GetLog(db *sql.DB, id string) (*string, error) {
	// check if log already exists (should no be the case)
	var content string
	db.QueryRow(ownDb.GetLogQuery, id).Scan(&content)
		
	if content != "" {
		return &content, nil		
	} else {
			content, err := CollectLogParts(db, id)
			if err != nil {
				return nil, err
			}
			return content, nil
	}
	
	return &content, nil
}

func SaveLog(db *sql.DB, reply *arc.Reply) error {
	if db == nil {
		return fmt.Errorf("Db is nil")
	}
	
	// save log part
	if reply.Payload != "" {
		log.Infof("Saving payload for reply with id %q, number %v, payload %q", reply.RequestID, reply.Number, reply.Payload)
		err := SaveLogPart(db, reply)
		if err != nil {
			return fmt.Errorf("Error saving log for request id %q. Got %q", reply.RequestID, err.Error())
		}
	}

	// collect log parts and save an entire log text
	if reply.Final == true {
		err := aggregateLogParts(db,reply.RequestID)
		if err != nil {
			return fmt.Errorf("Error aggregating log parts to log. Got %q", err.Error())
		}
		log.Infof("Aggregated log parts to log with id %q", reply.RequestID)
	}

	return nil
}

// private

func aggregateLogParts(db *sql.DB, id string) error {
	// check if log already exists (should no be the case)
	var content string
	db.QueryRow(ownDb.GetLogQuery, id).Scan(&content)
		
	if content != "" {
		// update log entry
		err := aggregateAndUpdate(db, id, content)
		if err != nil {
			return err
		}		
	} else {
		// insert new log entry as a transaction
		err := aggregateAndInsert(db, id)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func aggregateAndUpdate(db *sql.DB, id string, content string) (err error) {
	// collect log parts
	addContent, err := CollectLogParts(db, id)
	if err != nil {
		return
	}	
	
	// add content to the existing content
	newContent := fmt.Sprintf(content, addContent)
	
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
	
  if _, err = tx.Exec(ownDb.UpdateLogQuery, newContent, time.Now().Unix(), id); err != nil {
      return
  }
  if _, err = tx.Exec(ownDb.DeleteLogPartsQuery, id); err != nil {
      return
  }
	
	return
}

func aggregateAndInsert(db *sql.DB, id string) (err error) {
	// collect log parts
	content, err := CollectLogParts(db, id)
	if err != nil {
		return
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
	
  if _, err = tx.Exec(ownDb.InsertLogQuery, id, content, time.Now().Unix(), time.Now().Unix()); err != nil {
      return
  }
  if _, err = tx.Exec(ownDb.DeleteLogPartsQuery, id); err != nil {
      return
  }
	
	return
}