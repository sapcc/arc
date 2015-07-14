package models

import (
	"database/sql"
	"fmt"
	"time"
	"errors"

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

func CleanLogParts(db *sql.DB) error {	
	if db == nil {
		return errors.New("Db connection is nil")
	}

	// get log parts to aggregate
	rows, err := db.Query(ownDb.GetLogPartsToCleanQuery, 600, 84600) // 10 min and 1 day
	if err != nil {
		return err
	}
	defer rows.Close()
	
	rowsCount := 0
	var logPartID string
	for rows.Next() {
		// scan row
		err = rows.Scan(&logPartID)
		if err != nil {
			return err
		}	
			
		// aggregate the log parts found to the log table
		err = aggregateLogParts(db, logPartID)
		if err != nil{
			return err
		}
		rowsCount++
	}

	log.Infof("Clean log parts: %v aggregable log parts found", rowsCount)

	err = rows.Err()
	if err != nil {
		return err
	}

	rows.Close()
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
