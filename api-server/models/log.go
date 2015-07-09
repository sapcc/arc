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

func (log *Log) Get(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("Db connection is nil")
	}
	
	err := db.QueryRow(ownDb.GetLogQuery, log.JobID).Scan(&log.JobID, &log.Content, &log.CreatedAt, &log.UpdatedAt)		
	if err != nil {
		return err
	}
	
	return nil
}

func (log *Log) Save(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("Db connection is nil")
	}	
	
	var lastInsertId string
	err := db.QueryRow(ownDb.InsertLogQuery, log.JobID, log.Content, log.CreatedAt, log.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	
	return nil
}

func (log *Log) GetOrCollect(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("Db is nil")
	}
	
	err := log.Get(db)
	if err == sql.ErrNoRows {	
		// if no log entry collect all log parts
		log_part := LogPart{JobID:log.JobID}		
		content, err := log_part.Collect(db)
		if err != nil {
			return err
		}
		log.Content = *content
	}	else if err != nil {
		return err
	}
	
	return nil
}

func ProcessLogReply(db *sql.DB, reply *arc.Reply) error {
	if db == nil {
		return fmt.Errorf("Db connection is nil")
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
