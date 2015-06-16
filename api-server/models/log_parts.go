package models

import (
	"database/sql"
	"errors"
	"time"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

func CollectLogParts(db *sql.DB, id string) (*string, error) {
	if db == nil {
		return nil, errors.New("Db is nil")
	}

	var content string
	err := db.QueryRow(ownDb.CollectLogPartsQuery, id).Scan(&content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func SaveLogPart(db *sql.DB, reply *arc.Reply) error {
	if db == nil {
		return errors.New("Db is nil")
	}

	var lastInsertId string
	err := db.QueryRow(ownDb.InsertLogPartQuery, reply.RequestID, reply.Number, reply.Payload, reply.Final, time.Now()).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return nil
}
