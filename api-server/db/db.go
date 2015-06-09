package db

import (
	"database/sql"
	"fmt"
	"errors"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gopkg.in/gorp.v1"
)

var db *sql.DB
var dbmap *gorp.DbMap

func NewConnection(dbAddreess string) (*sql.DB, error) {
	var err error

	// conect to the db
	db, err = sql.Open("postgres", dbAddreess)
	if err != nil {
		return nil, err
	}
	
	log.Infof(fmt.Sprintf("Connected to the DB with address %q", dbAddreess))
	
	// create tables if needed
	err = createTables()
	if err != nil {
		return nil, err
	}
	
	// create a struct mapper
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(models.Job{}, "jobs")
	
	return db, nil
}

func CloseConnection() {
	if db != nil {
		db.Close()
	}
}

func SaveJob(job *models.Job) error {
	if db == nil || db == nil {		
		return errors.New("db connection or db mapper is nil")
	}
	
	err := dbmap.Insert(job)
	if err != nil {
		return err
	}
	
	return nil
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
