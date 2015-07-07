// +build integration

package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"		
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"	
	"code.google.com/p/go-uuid/uuid"	
	
	"time"
	"fmt"
	
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
			proj := "huhu project"
			org := "Miau organization"
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)
			
			// insert facts			
			req := arc.Request{Sender:agent_id, Payload: facts}
			var lastInsertId string
			err := db.QueryRow(InsertFactQuery, req.Sender, proj, org, req.Payload, time.Now(), time.Now()).Scan(&lastInsertId);
			Expect(err).NotTo(HaveOccurred())
			
			// get the facts
			newFact := Fact{Agent: Agent{AgentID: agent_id}}			
			err = newFact.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.Facts).To(Equal(facts))
			Expect(newFact.Project).To(Equal(proj))
			Expect(newFact.Organization).To(Equal(org))
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
			proj := "huhu project"
			org := "Miau organization"
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)
			
			// build a request
			req := arc.Request{Sender:agent_id, Payload: facts}
			
			// facts upate
			fact := Fact{Agent: Agent{AgentID: uuid.New()}}			
			err := fact.Update(db, &req)
			Expect(err).NotTo(HaveOccurred())
			
			// check
			newFact := Fact{}
			err = db.QueryRow(GetFactQuery, agent_id).Scan(&newFact.AgentID, &newFact.Project, &newFact.Organization, &newFact.Facts, &newFact.CreatedAt, &newFact.UpdatedAt)
			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.Facts).To(Equal(fact.Facts))			
			Expect(newFact.Project).To(Equal(proj))
			Expect(newFact.Organization).To(Equal(org))
		})		
		
		It("should update an existing entry", func() {
			agent_id := uuid.New()
			proj := "huhu project"
			org := "Miau organization"
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)
			
			// insert facts			
			req := arc.Request{Sender:agent_id, Payload: facts}
			var lastInsertId string
			err := db.QueryRow(InsertFactQuery, req.Sender, proj, org, req.Payload, time.Now(), time.Now()).Scan(&lastInsertId);
			Expect(err).NotTo(HaveOccurred())
			
			// build a request
			memory_used := 666666
			memory_total := 55555
			newFacts :=  fmt.Sprintf(`{"memory_used": %v, "memory_total": %v}`, memory_used, memory_total)
			newReq := arc.Request{Sender:agent_id, Payload: newFacts}
			
			// facts upate
			fact := Fact{Agent: Agent{AgentID: agent_id}}			
			err = fact.Update(db, &newReq)
			Expect(err).NotTo(HaveOccurred())
						
			// check
			dbFact := Fact{}
			dbFacts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": %v, "memory_total": %v, "organization": "%s"}`, proj, memory_used, memory_total, org)
			err = db.QueryRow(GetFactQuery, agent_id).Scan(&dbFact.AgentID, &dbFact.Project, &dbFact.Organization, &dbFact.Facts, &dbFact.CreatedAt, &dbFact.UpdatedAt)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbFact.Facts).To(Equal(dbFacts))
		})		
		
	})

})
