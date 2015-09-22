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

func NewConnection(dbConfigFile, env string) (*sql.DB, error) {
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

func DeleteAllRowsFromTable(db *sql.DB, table string) (sql.Result, error) {
	res, err := db.Exec(fmt.Sprint(DeleteQuery, table))
	if err != nil {
		return nil, err
	}
	return res, nil
}
