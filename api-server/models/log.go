package models

import (
	"database/sql"
	"errors"
	"time"
	
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type Log struct {
	RequestID string   `json:"request_id"`
	Id				string 	 `json:"id"`
	Payload   string   `json:"payload"`
}

type Logs []Log

func CollectLogs(db *sql.DB, requestID string) (*string, error) {
	
	var content string
	err := db.QueryRow(ownDb.GetLogsQuery, requestID).Scan(&content)
	if err != nil {
		return nil, err
	}
	
	return &content, nil
}

func SaveLog(db *sql.DB, reply *arc.Reply) error {
	if db == nil {
		return errors.New("Db is nil")
	}
	
	var lastInsertId string
	err := db.QueryRow(ownDb.InsertLogQuery, reply.RequestID, reply.Number, reply.Payload, time.Now().Unix()).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	
	return nil
}