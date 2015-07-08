package models

import (
	"database/sql"
	"errors"
	"time"
	"encoding/json"

	log "github.com/Sirupsen/logrus"

	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type DbConn interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Fact struct {
	Agent
	Facts string `json:"facts"`
}

func (fact *Fact) Get(dbc DbConn) error {
	if dbc == nil {
		return errors.New("Db connection is nil")
	}

	err := dbc.QueryRow(ownDb.GetFactQuery, fact.AgentID).Scan(&fact.AgentID, &fact.Project, &fact.Organization, &fact.Facts, &fact.CreatedAt, &fact.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (fact *Fact) FromRequest(req *arc.Request) error {	
	proj, org, err := extractProjectAndOrg(req.Payload)
	if err != nil {
		return err
	}
	
	fact.AgentID = req.Sender
	fact.Project = proj 
	fact.Organization = org
	fact.Facts = req.Payload
	fact.CreatedAt = time.Now()
	fact.UpdatedAt = time.Now()
	
	return nil
}

func (fact *Fact) Save(dbc DbConn) error {
	if dbc == nil {
		return errors.New("Db is nil")
	}
	
	var lastInsertId string
	if err := dbc.QueryRow(ownDb.InsertFactQuery, fact.AgentID, fact.Project, fact.Organization, fact.Facts, fact.CreatedAt, fact.UpdatedAt).Scan(&lastInsertId); err != nil {
		return err
	}
	
	log.Infof("New registry for sender %q will be saved.", fact.AgentID)
	
	return nil
}

func (fact *Fact) Update(dbc DbConn) error {
	if dbc == nil {
		return errors.New("Db is nil")
	}
	
	res, err := dbc.Exec(ownDb.UpdateFact, fact.AgentID, fact.Facts); 
	if err != nil {
		return err
	}	
	
	log.Infof("Registry for sender %q will be updated.", fact.AgentID)
	
	affect, err := res.RowsAffected(); 
	if err != nil {
		return err
	}

	log.Infof("%v rows where updated sender %q", affect, fact.AgentID)

	return nil
}

func (fact *Fact) ProcessRequest(db *sql.DB, req *arc.Request) (err error) {
	if db == nil {
		return errors.New("Db is nil")
	}
	
	err = fact.FromRequest(req)
	if err != nil {
		return
	}	

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	
	checkFact := Fact{Agent: Agent{AgentID: fact.AgentID}}
	err = checkFact.Get(tx)
	if err == sql.ErrNoRows { // fact not found		
		if err = fact.Save(tx); err != nil {
			return err
		}
	}	else if err != nil { // something wrong happned
		return
	} else {
		if err = fact.Update(tx); err != nil {
			return
		}
	}
	
	// update object data
	err = tx.QueryRow(ownDb.GetFactQuery, req.Sender).Scan(&fact.AgentID, &fact.Project, &fact.Organization, &fact.Facts, &fact.CreatedAt, &fact.UpdatedAt)
	if err != nil {
		return
	}
	
	return
}

func (fact *Fact) ProcessRequest_old(db *sql.DB, req *arc.Request) (err error) {
	if db == nil {
		return errors.New("Db is nil")
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// insert or update
	err = tx.QueryRow(ownDb.GetFactQuery, req.Sender).Scan(&fact.AgentID, &fact.Project, &fact.Organization, &fact.Facts, &fact.CreatedAt, &fact.UpdatedAt)
	if err == nil {
		log.Infof("Registry for sender %q will be updated.", req.Sender)
		if _, err = tx.Exec(ownDb.UpdateFact, req.Sender, req.Payload); err != nil {
			return
		}
	} else if err == sql.ErrNoRows {
		log.Infof("New registry for sender %q will be saved.", req.Sender)
		proj, org, err := extractProjectAndOrg(req.Payload)
		if err != nil {
			return err
		}
		var lastInsertId string
		if err = tx.QueryRow(ownDb.InsertFactQuery, req.Sender, proj, org, req.Payload, time.Now(), time.Now()).Scan(&lastInsertId); err != nil {
			return err
		}
	} else if err != nil {
		return
	}

	// update object data
	err = tx.QueryRow(ownDb.GetFactQuery, req.Sender).Scan(&fact.AgentID, &fact.Project, &fact.Organization, &fact.Facts, &fact.CreatedAt, &fact.UpdatedAt)
	if err != nil {
		return
	}

	return
}

// private

func extractProjectAndOrg(payload string) (proj string, org string, err error) {
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal([]byte(payload), &objmap)
	if err != nil {
		return 
	}
	
	// get the project
	if objmap["project"] != nil {
		err = json.Unmarshal(*objmap["project"], &proj)
		if err != nil {
			return 
		}	
	}
	
	// get the organization
	if objmap["project"] != nil {
		err = json.Unmarshal(*objmap["organization"], &org)
		if err != nil {
			return 
		}		
	}	
	
	return
}
