package models

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type Log struct {
	JobID 			string		`json:"job_id"`
	Content			string		`json:"content"`
	CreatedAt   time.Time	`json:"created_at"`
	UpdatedAt   time.Time	`json:"updated_at"`
}

func (log *Log) Get(db *sql.DB, id string) error {
	err := db.QueryRow(ownDb.GetLogQuery, id).Scan(&log.JobID, &log.Content, &log.CreatedAt, &log.UpdatedAt)	
	
	// check if the entry already exists
	if err != nil {
		// if no log entry collect all log parts
		log_part := LogPart{JobID:id}		
		content, err := log_part.Collect(db)
		if err != nil {
			return err
		}
		log.JobID = id
		log.Content = *content
	}	
	return nil
}

func ProcessLogReply(db *sql.DB, reply *arc.Reply) error {
	if db == nil {
		return fmt.Errorf("Db is nil")
	}

	// save log part
	if reply.Payload != "" {
		log.Infof("Saving payload for reply with id %q, number %v, payload %q", reply.RequestID, reply.Number, reply.Payload)
		
		logPart := LogPart{reply.RequestID, reply.Number, reply.Payload, reply.Final, time.Now()}
		err := logPart.Save(db)
		if err != nil {
			return fmt.Errorf("Error saving log for request id %q. Got %q", reply.RequestID, err.Error())
		}
	}

	// collect log parts and save an entire log text
	if reply.Final == true {
		err := aggregateLogParts(db, reply.RequestID)
		if err != nil {
			return fmt.Errorf("Error aggregating log parts to log. Got %q", err.Error())
		}
		log.Infof("Aggregated log parts to log with id %q", reply.RequestID)
	}

	return nil
}

// private

func aggregateLogParts(db *sql.DB, id string) (err error) {
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

	var content string
	if err = tx.QueryRow(ownDb.CollectLogPartsQuery, id).Scan(&content); err != nil {
		return
	}
	if _, err = tx.Exec(ownDb.InsertLogQuery, id, content, time.Now(), time.Now()); err != nil {
		return
	}
	if _, err = tx.Exec(ownDb.DeleteLogPartsQuery, id); err != nil {
		return
	}

	return
}
