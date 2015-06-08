package main

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

const (
	DB_AGENTS_TABLE = "Agents"
	DB_FACTS_TABLE  = "Facts"
	DB_JOBS_TABLE   = "Jobs"
)

type Agents struct {
	uid sql.NullInt64
}

type Jobs struct {
	reqID   sql.NullInt64
	payload sql.NullString
	status  sql.NullInt64
}

type Facts struct {
	factID   sql.NullInt64
	name     sql.NullString
	value    sql.NullString
	agent_id sql.NullInt64
}

var db *sql.DB

func NewDb(dbAddreess string) {
	var err error

	// conect to the db
	db, err = sql.Open("postgres", dbAddreess)
	if err != nil {
		panic(err)
	} else {
		log.Infof(fmt.Sprintf("Connected to the DB with address %q", dbAddreess))
	}

	createTables()
}

func CreateJob(id string, payload string) error {
	_, err := db.Exec(fmt.Sprintf("Insert into %s (reqID,payload,status) values (%s,%s,%s)", DB_JOBS_TABLE, id, payload, models.Queued))
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// private

func createTables() {
	var err error

	// create agents table
	if _, err = execQuery(db, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( reqID varchar(255) PRIMARY KEY NOT null, payload text Not null, status varchar(255) NOT null, CONSTRAINT uc_reqID UNIQUE (reqID))", DB_JOBS_TABLE)); err != nil {
		panic(err)
	}

	// create agents table
	/*if _, err = execQuery(db, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( uid integer PRIMARY KEY NOT null, CONSTRAINT uc_AgentsID UNIQUE (uid))", DB_AGENTS_TABLE)); err != nil {
		panic(err)
	}*/

	// create facts table
	/*if _, err = execQuery(db, fmt.Sprintf("create table IF NOT EXISTS %s (factID integer PRIMARY KEY NOT null,name varchar(255) NOT null,value varchar(255) NOT null,agent_id int, FOREIGN KEY (agent_id) REFERENCES agents(uid), CONSTRAINT uc_factID UNIQUE (factID, name))", DB_FACTS_TABLE)); err != nil {
		panic(err)
	}*/
}

func fillWithExamples() {
	// instert examples to the agents table
	_, err := db.Exec(fmt.Sprintf("Insert into %s (uid) values (1)", DB_AGENTS_TABLE))
	if err != nil {
		log.Fatal(err)
	}
	// instert examples to the facts table
	_, err = db.Exec(fmt.Sprintf("Insert into %s (factID,name,value,agent_id) values (2,'os','windows',1)", DB_FACTS_TABLE))
	if err != nil {
		log.Fatal(err)
	}
}

func execQuery(db *sql.DB, query string) (sql.Result, error) {
	res, err := db.Exec(query)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	return res, nil
}
