package models

import (
	"database/sql"
	"errors"
	"time"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

type LogPart struct {
	JobID     string    `json:"job_id"`
	Number    uint      `json:"number"`
	Content   string    `json:"content"`
	Final     bool      `json:"final"`
	CreatedAt time.Time `json:"created_at"`
}

func (log_part *LogPart) Collect(db *sql.DB) (*string, error) {
	if db == nil {
		return nil, errors.New("db connection is nil")
	}

	//var content string
	var content sql.NullString
	if err := db.QueryRow(ownDb.CollectLogPartsQuery, log_part.JobID).Scan(&content); err != nil {
		return nil, err
	}
	if !content.Valid {
		return nil, sql.ErrNoRows
	}

	return &content.String, nil
}

func (log_part *LogPart) Get(db *sql.DB) error {
	if db == nil {
		return errors.New("db connection is nil")
	}

	return db.QueryRow(ownDb.GetLogPartQuery, log_part.JobID, log_part.Number).Scan(&log_part.JobID, &log_part.Number, &log_part.Content, &log_part.Final, &log_part.CreatedAt)
}

func (log_part *LogPart) Save(db *sql.DB) error {
	if db == nil {
		return errors.New("db connection is nil")
	}

	var lastInsertId string
	return db.QueryRow(ownDb.InsertLogPartQuery, log_part.JobID, log_part.Number, log_part.Content, log_part.Final, log_part.CreatedAt).Scan(&lastInsertId)
}
