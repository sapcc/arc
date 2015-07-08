// +build integration

package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"	
	"code.google.com/p/go-uuid/uuid"	

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
			
			// save facts for agent
			fact := Fact{Agent: Agent{AgentID: agent_id, Project: proj, Organization: org}, Facts: facts}
			err := fact.Save(db)
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

	Describe("ProcessRequest", func() {

		It("returns an error if no db connection is given", func() {
			newFact := Fact{Agent: Agent{AgentID: uuid.New()}}			
			err := newFact.ProcessRequest(nil, nil)
			Expect(err).To(HaveOccurred())
		})
		
		It("should insert a new entry if facts doesn't exist yet", func() {
			agent_id := uuid.New()
			proj := "huhu project"
			org := "Miau organization"
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)
			
			// build a request
			req := arc.Request{Sender:agent_id, Payload: facts}
			
			// process the request
			fact := Fact{}
			err := fact.ProcessRequest(db, &req)
			Expect(err).NotTo(HaveOccurred())
			
			// check
			newFact := Fact{Agent: Agent{AgentID: agent_id}}	
			err = newFact.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.Facts).To(Equal(fact.Facts))			
			Expect(newFact.Project).To(Equal(proj))
			Expect(newFact.Organization).To(Equal(org))
		})		

		It("should update an existing entry", func() {
			proj := "huhu project"
			org := "Miau organization"
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)
			
			// save agent with facts			
			req := arc.Request{Sender:uuid.New(), Payload: facts}
			fact := Fact{}; fact.FromRequest(&req)
			err := fact.Save(db)
			Expect(err).NotTo(HaveOccurred())
			
			// build a request
			memory_used := 666666
			memory_total := 55555
			newFacts :=  fmt.Sprintf(`{"memory_used": %v, "memory_total": %v}`, memory_used, memory_total)
			newReq := arc.Request{Sender:fact.AgentID, Payload: newFacts}
			
			// process the request
			updateFact := Fact{}			
			err = updateFact.ProcessRequest(db, &newReq)
			Expect(err).NotTo(HaveOccurred())
						
			// check
			dbFacts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": %v, "memory_total": %v, "organization": "%s"}`, proj, memory_used, memory_total, org)
			dbFact := Fact{Agent: Agent{AgentID: fact.AgentID}}
			err = dbFact.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbFact.Facts).To(Equal(dbFacts))
		})		
		
	})

})
