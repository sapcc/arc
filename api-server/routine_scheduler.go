package main

import (
	log "github.com/Sirupsen/logrus"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"time"
)

func routineScheduler() {
	cleanJobsChan := time.NewTicker(time.Second * 60)

	for {
		select {
		case <-cleanJobsChan.C:
			go cleanJobs()
		}
	}

}

func cleanJobs() {
	log.Info("Clean job routine started")

	if db == nil {
		log.Error("Db is nil")
	}

	// add an interval in seconds to the timeout
	res, err := db.Exec(ownDb.CleanJobsQuery, 60)
	if err != nil {
		log.Error(err.Error())
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Error(err.Error())
	}

	log.Infof("%v jobs where updated", affect)
}

func cleanLogParts() {
}
