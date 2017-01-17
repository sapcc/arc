package main

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

func routineScheduler(db *sql.DB, duration time.Duration) {

	routineSchedulerChan := time.NewTicker(duration)

	for {
		select {
		case <-routineSchedulerChan.C:
			runRoutineTasks(db)
		}
	}

}

func runRoutineTasks(db *sql.DB) {
	cleanJobs(db)
	cleanLogParts(db)
	cleanLocks(db)
	cleanStaleTokens(db)
	cleanOldCertificates(db)
}

func cleanJobs(db *sql.DB) {
	affectHeartbeatJobs, affectTimeOutJobs, affectOldJobs, err := models.CleanJobs(db)
	if err != nil {
		log.Error("Clean jobs: ", err.Error())
	}
	log.Infof("Clean job routine : %v jobs without heartbeat answer and %v timeout jobs where updated. %v old jobs where deleted", affectHeartbeatJobs, affectTimeOutJobs, affectOldJobs)
}

func cleanLogParts(db *sql.DB) {
	rowsCount, err := models.CleanLogParts(db)
	if err != nil {
		log.Error("Clean log parts routine: ", err.Error())
	}
	log.Infof("Clean log parts routine: %v aggregable log parts found", rowsCount)
}

func cleanLocks(db *sql.DB) {
	affectedLocks, err := models.CleanLocks(db)
	if err != nil {
		log.Error("Clean locks routine: ", err.Error())
	}
	log.Infof("Clean locks routine: %v old (5 min) locks removed from the db", affectedLocks)
}

func cleanStaleTokens(db *sql.DB) {
	res, err := db.Exec(ownDb.CleanPkiTokensQuery, 3600)
	if err != nil {
		log.Error("Failed to prune tokens: ", err.Error())
		return
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		log.Error("Failed to get count of prune tokens: ", err.Error())
		return
	}
	log.Infof("Clean pki tokens routine: %v tokens removed from the db", affectedRows)
}

func cleanOldCertificates(db *sql.DB) {
	res, err := db.Exec(ownDb.CleanPkiCertificatesQuery)
	if err != nil {
		log.Error("Failed to prune tokens: ", err.Error())
		return
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		log.Error("Failed to get count of removed certificates: ", err.Error())
		return
	}
	log.Infof("Clean pki certificates routine: %v certificates removed from the db", affectedRows)
}
