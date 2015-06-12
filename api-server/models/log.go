package models

import (
	"database/sql"
	"errors"
	
	log "github.com/Sirupsen/logrus"
	
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type Log struct {
	RequestID string   `json:"request_id"`
	Id				string 	 `json:"id"`
	Payload   string   `json:"payload"`
}

type Logs []Log

func GetLogs(db *sql.DB, requestID string) (*Logs, error) {
	var resLogs Logs
	rows, err := db.Query(ownDb.GetLogsQuery, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var resLog Log
	for rows.Next() {
		err = rows.Scan(&resLog.RequestID, &resLog.Id, &resLog.Payload)
		if err != nil {
			log.Errorf("Error scaning log results. Got ", err.Error())
			continue
		}
		resLogs = append(resLogs, resLog)
	}
	
	return &resLogs, nil
}

func SaveLog(db *sql.DB, reply *arc.Reply) error {
	if db == nil {
		return errors.New("Db is nil")
	}
	
	var lastInsertId string
	err := db.QueryRow(ownDb.InsertLogQuery, reply.RequestID, reply.Payload).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	
	return nil
}