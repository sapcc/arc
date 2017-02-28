package janitor

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pki"
)

type CleanJobs struct {
	db *sql.DB
}

func (c CleanJobs) Run() {
	affectHeartbeatJobs, affectTimeOutJobs, affectOldJobs, err := models.CleanJobs(c.db)
	if err != nil {
		log.Error("Clean jobs: ", err.Error())
	}
	log.Infof("Clean job routine : %v jobs without heartbeat answer and %v timeout jobs where updated. %v old jobs where deleted", affectHeartbeatJobs, affectTimeOutJobs, affectOldJobs)
}

type CleanLogParts struct {
	db *sql.DB
}

func (c CleanLogParts) Run() {
	rowsCount, err := models.CleanLogParts(c.db)
	if err != nil {
		log.Error("Clean log parts routine: ", err.Error())
	}
	log.Infof("Clean log parts routine: %v aggregable log parts found", rowsCount)
}

type CleanLocks struct {
	db *sql.DB
}

func (c CleanLocks) Run() {
	affectedLocks, err := models.CleanLocks(c.db)
	if err != nil {
		log.Error("Clean locks routine: ", err.Error())
	}
	log.Infof("Clean locks routine: %v old (5 min) locks removed from the db", affectedLocks)
}

type CleanTokens struct {
	db *sql.DB
}

func (c CleanTokens) Run() {
	affectedRows, err := pki.CleanOldTokens(c.db)
	if err != nil {
		log.Error("Failed to get count of prune tokens: ", err.Error())
		return
	}
	log.Infof("Clean pki tokens routine: %v tokens removed from the db", affectedRows)
}

type CleanCertificates struct {
	db *sql.DB
}

func (c CleanCertificates) Run() {
	affectedRows, err := pki.CleanOldCertificates(c.db)
	if err != nil {
		log.Error("Failed to get count of removed certificates: ", err.Error())
		return
	}
	log.Infof("Clean pki certificates routine: %v certificates removed from the db", affectedRows)
}
