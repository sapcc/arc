//lint:file-ignore SA1019 need to be synchronize with the grafana dashboards.
// TODO
// warning: prometheus.InstrumentHandler is deprecated: InstrumentHandler has several issues. Use the tooling provided in package promhttp instead.
// The issues are the following:
// (1) It uses Summaries rather than Histograms. Summaries are not useful if aggregation across multiple instances is required.
// (2) It uses microseconds as unit, which is deprecated and should be replaced by seconds.
// (3) The size of the request is calculated in a separate goroutine. Since this calculator requires access to the request header,
//     it creates a race with any writes to the header performed during request handling.  httputil.ReverseProxy is a prominent example for a handler performing such writes.
// (4) It has additional issues with HTTP/2, cf. https://github.com/prometheus/client_golang/issues/272.  (SA1019) (staticcheck)
// FIX proposal: https://gitlab.cncf.ci/prometheus/prometheus/commit/83325c8d822d022fec74d21e2efd15e3b6b6a0af

package janitor

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/bamzi/jobrunner"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/prometheus/client_golang/prometheus"
)

var SCHEDULE_TIME = "@every 60s"

type Janitor struct {
	Conf JanitorConf
}

type JanitorConf struct {
	BindAddress  string
	DbConfigFile string
	Environment  string
}

func InitJanitor(conf JanitorConf) *Janitor {
	jobrunner.Start()

	return &Janitor{
		Conf: conf,
	}
}

func (j *Janitor) Stop() {
	jobrunner.Stop()
}

func (j *Janitor) Start() {
	jobrunner.Start()
}

func (j *Janitor) InitServer() {
	log.Printf("Starting arc janitor server %s", VersionString())

	// live monitoring
	http.HandleFunc("/", serveVersion)
	http.HandleFunc("/jobrunner", serveJobrunner)
	http.Handle("/metrics", prometheus.Handler())

	log.Printf("Listening on %s...", j.Conf.BindAddress)
	err := http.ListenAndServe(j.Conf.BindAddress, nil)
	fatalfOnError(err, "failed to bind on %s: ", j.Conf.BindAddress)
}

func (j *Janitor) InitScheduler() {
	log.Printf("Starting arc janitor scheduler %s", VersionString())
	db, err := dbConnection(j.Conf.DbConfigFile, j.Conf.Environment)
	fatalfOnError(err, "failed to bind to database ")

	// start clean jobs
	if err = jobrunner.Schedule(SCHEDULE_TIME, FailQueuedJobs{db: db}); err != nil {
		log.Errorf(fmt.Sprintf("janitor job 'FailQueuedJobs' scheduler failed: %s", err))
	}
	if err = jobrunner.Schedule(SCHEDULE_TIME, FailExpiredJobs{db: db}); err != nil {
		log.Errorf(fmt.Sprintf("janitor job 'FailExpiredJobs' scheduler failed: %s", err))
	}
	if err = jobrunner.Schedule(SCHEDULE_TIME, PruneJobs{db: db}); err != nil {
		log.Errorf(fmt.Sprintf("janitor job 'PruneJobs' scheduler failed: %s", err))
	}
	if err = jobrunner.Schedule(SCHEDULE_TIME, AggregateLogs{db: db}); err != nil {
		log.Errorf(fmt.Sprintf("janitor job 'AggregateLogs' scheduler failed: %s", err))
	}
	if err = jobrunner.Schedule(SCHEDULE_TIME, PruneLocks{db: db}); err != nil {
		log.Errorf(fmt.Sprintf("janitor job 'PruneLocks' scheduler failed: %s", err))
	}
	if err = jobrunner.Schedule(SCHEDULE_TIME, PruneCertificates{db: db}); err != nil {
		log.Errorf(fmt.Sprintf("janitor job 'PruneCertificates' scheduler failed: %s", err))
	}
}

func dbConnection(dbConfigFile, env string) (*sql.DB, error) {
	// check and load config file
	if _, err := os.Stat(dbConfigFile); err != nil {
		return nil, fmt.Errorf("can't load database configuration file %s: %s", dbConfigFile, err)
	}
	f, err := yaml.ReadFile(dbConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration file %s: %s", dbConfigFile, err)
	}
	// read the environment
	open, err := f.Get(fmt.Sprintf("%s.open", env))
	if err != nil {
		return nil, fmt.Errorf("can't find 'open' key for %s environment ", env)
	}
	dbDSN := os.ExpandEnv(open)

	// conect to the db
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		return nil, err
	}

	//connection is defered until the first query unless we ping
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// hide user data
	logDSN := regexp.MustCompile(`password=[^ ]+`).ReplaceAllString(dbDSN, "password=****")
	logDSN = regexp.MustCompile(`:[^/:@]+@`).ReplaceAllString(logDSN, ":****@")

	log.Infof(fmt.Sprintf("Connected to the DB with address %q", logDSN))

	return db, nil
}

func fatalfOnError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalf(msg, args...)
	}
}
