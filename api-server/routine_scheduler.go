package main

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

var routineSchedulerChan *time.Ticker

func routineScheduler(db *sql.DB, duration time.Duration) error {
	if db == nil {
		return fmt.Errorf("Db connection is nil")
	}

	routineSchedulerChan = time.NewTicker(duration)

	for {
		select {
		case <-routineSchedulerChan.C:
			cleanJobs(db)
			cleanLogParts(db)
		}
	}

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
