package db

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/kylelemons/go-gypsy/yaml"
	_ "github.com/lib/pq"
)

var db *sql.DB

func NewConnection(dbConfigFile, env string) (*sql.DB, error) {
	if _, err := os.Stat(dbConfigFile); err != nil {
		return nil, fmt.Errorf("Can't load database configuration file %s: %s", dbConfigFile, err)
	}
	f, err := yaml.ReadFile(dbConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse database configuration file %s: %s", dbConfigFile, err)
	}
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

	logDSN := regexp.MustCompile(`password=[^ ]+`).ReplaceAllString(dbDSN, "password=****")
	logDSN = regexp.MustCompile(`:[^/:@]+@`).ReplaceAllString(logDSN, ":****@")

	log.Infof(fmt.Sprintf("Connected to the DB with address %q", logDSN))

	return db, nil
}

// private

func execQuery(db *sql.DB, query string) (sql.Result, error) {
	res, err := db.Exec(query)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	return res, nil
}
