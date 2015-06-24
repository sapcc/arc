package db

import (
	"database/sql"
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

var db *sql.DB

func NewConnection(dbAddress string) (*sql.DB, error) {
	var err error

	// conect to the db
	db, err = sql.Open("postgres", dbAddress)
	if err != nil {
		return nil, err
	}

	dbAddress = regexp.MustCompile(`password=[^ ]+`).ReplaceAllString(dbAddress, "password=****")
	dbAddress = regexp.MustCompile(`:[^/:@]+@`).ReplaceAllString(dbAddress, ":****@")

	log.Infof(fmt.Sprintf("Connected to the DB with address %q", dbAddress))

	// create tables if needed
	err = createTables()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// private

func createTables() error {
	var err error
	for _, t := range Tables {
		if _, err = execQuery(db, t); err != nil {
			break
		}
	}
	return err
}

func execQuery(db *sql.DB, query string) (sql.Result, error) {
	res, err := db.Exec(query)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	return res, nil
}
