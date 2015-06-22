package main

import (
	"time"
	"fmt"
	"database/sql"
	
	log "github.com/Sirupsen/logrus"
	
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

func routineScheduler(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("Db connection is nil")
	}
	
	cleanJobsChan := time.NewTicker(time.Second * 60)

	for {
		select {
		case <-cleanJobsChan.C:
			go cleanJobs(db)
		}
	}

}

func cleanJobs(db *sql.DB) {
	log.Info("Clean job routine started")
	
	// clean jobs which no heartbeat was send back after created_at + 60 sec
	res, err := db.Exec(ownDb.CleanJobsNonHeartbeatQuery, 60)
	if err != nil {
		log.Error(err.Error())
	}
	
	affect, err := res.RowsAffected()
	if err != nil {
		log.Error(err.Error())
	}	
	
	log.Infof("%v jobs without heartbeat answer where updated", affect)
	
	// clean jobs which the timeout + 60 sec has exceeded and still in queued or executing status
	res, err = db.Exec(ownDb.CleanJobsTimeoutQuery, 60)
	if err != nil {
		log.Error(err.Error())
	}

	affect, err = res.RowsAffected()
	if err != nil {
		log.Error(err.Error())
	}

	log.Infof("%v timeout jobs where updated", affect)
}

func cleanLogParts() {
}
