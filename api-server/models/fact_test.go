// +build integration

package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"		
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"	
	"code.google.com/p/go-uuid/uuid"	
	
	"time"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fact", func() {

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			newFact := Fact{Agent: Agent{AgentID: uuid.New()}}			
			err := newFact.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if facts not found", func() {
			newFact := Fact{Agent: Agent{AgentID: uuid.New()}}			
			err := newFact.Get(db)
			Expect(err).To(HaveOccurred())
		})

		It("should return the facts", func() {			
			agent_id := uuid.New()
			facts := `{"os": "darwin", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
			
			// insert facts			
			req := arc.Request{Sender:agent_id, Payload: facts}
			var lastInsertId string
			err := db.QueryRow(InsertFactQuery, req.Sender, req.Payload, time.Now(), time.Now()).Scan(&lastInsertId);
			Expect(err).NotTo(HaveOccurred())
			
			// get the facts
			newFact := Fact{Agent: Agent{AgentID: agent_id}}			
			err = newFact.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.Facts).To(Equal(facts))
		})
		
	})

	Describe("Update", func() {

		It("returns an error if no db connection is given", func() {
			newFact := Fact{Agent: Agent{AgentID: uuid.New()}}			
			err := newFact.Update(nil, nil)
			Expect(err).To(HaveOccurred())
		})
		
		It("should insert a new entry if facts doesn't exist yet", func() {
			agent_id := uuid.New()
			facts := `{"os": "darwin", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
			
			// build a request
			req := arc.Request{Sender:agent_id, Payload: facts}
			
			// facts upate
			fact := Fact{Agent: Agent{AgentID: uuid.New()}}			
			err := fact.Update(db, &req)
			Expect(err).NotTo(HaveOccurred())
			
			// check
			newFact := Fact{}
			err = db.QueryRow(GetFactQuery, agent_id).Scan(&newFact.AgentID, &newFact.Facts, &newFact.CreatedAt, &newFact.UpdatedAt)
			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.Facts).To(Equal(fact.Facts))			
		})		
		
		It("should update an existing entry", func() {
			agent_id := uuid.New()
			facts := `{"os": "darwin", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
			
			// insert facts			
			req := arc.Request{Sender:agent_id, Payload: facts}
			var lastInsertId string
			err := db.QueryRow(InsertFactQuery, req.Sender, req.Payload, time.Now(), time.Now()).Scan(&lastInsertId);
			Expect(err).NotTo(HaveOccurred())
			
			// build a request
			newFacts :=  `{"memory_used": 666666, "memory_total": 55555}`			
			newReq := arc.Request{Sender:agent_id, Payload: newFacts}
			
			// facts upate
			fact := Fact{Agent: Agent{AgentID: agent_id}}			
			err = fact.Update(db, &newReq)
			Expect(err).NotTo(HaveOccurred())
						
			// check
			dbFact := Fact{}
			err = db.QueryRow(GetFactQuery, agent_id).Scan(&dbFact.AgentID, &dbFact.Facts, &dbFact.CreatedAt, &dbFact.UpdatedAt)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbFact.Facts).To(Equal(`{"os": "darwin", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 666666, "memory_total": 55555, "organization": "test-org"}`))			
		})		
		
	})

})
