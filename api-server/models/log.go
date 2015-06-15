package models

import (
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type Log struct {
	RequestID string `json:"request_id"`
	Id        string `json:"id"`
	Payload   string `json:"payload"`
}

type Logs []Log

func GetLog(db *sql.DB, requestID string) (*string, error) {
	return nil, nil
}

func SaveLog(db *sql.DB, reply *arc.Reply) error {

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

	}

	return nil
}
