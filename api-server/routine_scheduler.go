package main

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"

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
}

func cleanJobs(db *sql.DB) {
	log.Info("Clean job routine started")

	err := models.CleanJobs(db)
	if err != nil {
		log.Error("Clean jobs: ", err.Error())
	}
}

func cleanLogParts(db *sql.DB) {
	log.Info("Clean log parts routine started")

	err := models.CleanLogParts(db)
	if err != nil {
		log.Error("Clean log parts: ", err.Error())
	}
}
