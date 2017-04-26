package janitor

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pki"
)

// fail jobs which no heartbeat was send back after created_at + 60 sec
type FailQueuedJobs struct {
	db *sql.DB
}

func (c FailQueuedJobs) Run() {
	affectedJobs, err := models.FailQueuedJobs(c.db)
	if err != nil {
		log.Error("FailQueuedJobs scheduler: ", err.Error())
	}
	log.Infof("FailQueuedJobs scheduler: %v jobs without heartbeat answer.", affectedJobs)
}

// fail jobs which the timeout + 60 sec has exceeded and still in queued or executing status
type FailExpiredJobs struct {
	db *sql.DB
}

func (c FailExpiredJobs) Run() {
	affectedJobs, err := models.FailExpiredJobs(c.db)
	if err != nil {
		log.Error("FailExpiredJobs scheduler: ", err.Error())
	}
	log.Infof("FailExpiredJobs scheduler: %v timeout jobs.", affectedJobs)
}

// delete jobs which are older than 30 days
type PruneJobs struct {
	db *sql.DB
}

func (c PruneJobs) Run() {
	affectedJobs, err := models.PruneJobs(c.db)
	if err != nil {
		log.Error("PruneJobs scheduler: ", err.Error())
	}
	log.Infof("PruneJobs scheduler: %v old jobs.", affectedJobs)
}

// aggregate all log parts when final log_part created_at > 5 min or older the 1 day
type AggregateLogs struct {
	db *sql.DB
}

func (c AggregateLogs) Run() {
	rowsCount, err := models.AggregateLogs(c.db)
	if err != nil {
		log.Error("AggregateLogs scheduler: ", err.Error())
	}
	log.Infof("AggregateLogs scheduler: %v aggregable log parts found", rowsCount)
}

type PruneLocks struct {
	db *sql.DB
}

func (c PruneLocks) Run() {
	affectedLocks, err := models.PruneLocks(c.db)
	if err != nil {
		log.Error("PruneLocks scheduler: ", err.Error())
	}
	log.Infof("PruneLocks scheduler %v old (5 min) locks removed from the db", affectedLocks)
}

type PruneTokens struct {
	db *sql.DB
}

func (c PruneTokens) Run() {
	affectedRows, err := pki.PruneTokens(c.db)
	if err != nil {
		log.Error("PruneTokens scheduler: Failed to get count of prune tokens: ", err.Error())
		return
	}
	log.Infof("PruneTokens scheduler: %v tokens removed from the db", affectedRows)
}

type PruneCertificates struct {
	db *sql.DB
}

func (c PruneCertificates) Run() {
	affectedRows, err := pki.PruneCertificates(c.db)
	if err != nil {
		log.Error("Failed to get count of removed certificates: ", err.Error())
		return
	}
	log.Infof("Clean pki certificates routine: %v certificates removed from the db", affectedRows)
}
