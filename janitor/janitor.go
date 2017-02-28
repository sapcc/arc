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
	log.Printf("Listening on %s...", j.Conf.BindAddress)
	err := http.ListenAndServe(j.Conf.BindAddress, nil)
	fatalfOnError(err, "Failed to bind on %s: ", j.Conf.BindAddress)
}

func (j *Janitor) InitScheduler() {
	log.Printf("Starting arc janitor scheduler %s", VersionString())
	db, err := dbConnection(j.Conf.DbConfigFile, j.Conf.Environment)
	fatalfOnError(err, "Failed to bind to database ")

	// start clean jobs
	jobrunner.Schedule(SCHEDULE_TIME, CleanJobs{db: db})
	jobrunner.Schedule(SCHEDULE_TIME, CleanLogParts{db: db})
	jobrunner.Schedule(SCHEDULE_TIME, CleanLocks{db: db})
	jobrunner.Schedule(SCHEDULE_TIME, CleanTokens{db: db})
	jobrunner.Schedule(SCHEDULE_TIME, CleanCertificates{db: db})
}

func dbConnection(dbConfigFile, env string) (*sql.DB, error) {
	// check and load config file
	if _, err := os.Stat(dbConfigFile); err != nil {
		return nil, fmt.Errorf("Can't load database configuration file %s: %s", dbConfigFile, err)
	}
	f, err := yaml.ReadFile(dbConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse database configuration file %s: %s", dbConfigFile, err)
	}
	// read the environment
	open, err := f.Get(fmt.Sprintf("%s.open", env))
	if err != nil {
		return nil, fmt.Errorf("Can't find 'open' key for %s environment ", env)
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
