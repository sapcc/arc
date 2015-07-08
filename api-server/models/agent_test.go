// +build integration

package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"time"
)

var _ = Describe("Agents", func() {

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			agents := Agents{}
			err := agents.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should return all agents", func() {
			agents := Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)

			// insert facts / agent
			dbAgents := Agents{}
			err := dbAgents.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))
			Expect(dbAgents[2].AgentID).To(Equal(agents[2].AgentID))
		})

	})

})

var _ = Describe("Agent", func() {

	Describe("Get and Save", func() {

		It("returns an error if no db connection is given", func() {
			agent := Agent{}
			agent.Example()
			err := agent.Get(nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if no agent found", func() {
			agent := Agent{}
			agent.Example()
			err := agent.Get(db)
			Expect(err).To(HaveOccurred())
		})

		It("should return the agent", func() {
			// insert facts / agent
			agent := Agent{}
			agent.Example()
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// get agent
			newAgent := Agent{AgentID: agent.AgentID}
			err = newAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(agent.AgentID).To(Equal(newAgent.AgentID))
			Expect(agent.Project).To(Equal(newAgent.Project))
			Expect(agent.Organization).To(Equal(newAgent.Organization))
			Expect(agent.Facts).To(Equal(newAgent.Facts))
			Expect(agent.CreatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(newAgent.CreatedAt.Format("2006-01-02 15:04:05.99")))
			Expect(agent.UpdatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(newAgent.UpdatedAt.Format("2006-01-02 15:04:05.99")))
		})

	})

	Describe("Update", func() {

		It("returns an error if no db connection is given", func() {
			newAgent := Agent{AgentID: uuid.New()}
			err := newAgent.Update(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should update project, organization, facts and updated_at", func() {
			facts := `{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": %v, "memory_total": %v, "organization": "%s"}`
			// save an agent
			agent := Agent{
				AgentID:      uuid.New(),
				Project:      "huhu project",
				Organization: "huhu organization",
				Facts:        fmt.Sprintf(facts, "huhu project", 123456789, 987654321, "huhu organization"),
				CreatedAt:    time.Now().Add((-5) * time.Minute),
				UpdatedAt:    time.Now().Add((-5) * time.Minute),
			}
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// update same agent
			newProj := "Miau project"
			newOrg := "Miau organization"
			memory_used := 666666
			memory_total := 55555
			newFacts := fmt.Sprintf(`{"memory_used": %v, "memory_total": %v, "project": "%s", "organization": "%s"}`, memory_used, memory_total, newProj, newOrg)
			updateAgent := Agent{
				AgentID:      agent.AgentID,
				Project:      newProj,
				Organization: newOrg,
				Facts:        newFacts,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			err = updateAgent.Update(db)
			Expect(err).NotTo(HaveOccurred())

			// check
			dbAgent := Agent{AgentID: agent.AgentID}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgent.Facts).To(Equal(fmt.Sprintf(facts, newProj, memory_used, memory_total, newOrg)))
			Expect(dbAgent.Project).To(Equal(newProj))
			Expect(dbAgent.Organization).To(Equal(newOrg))
			Expect(dbAgent.CreatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(agent.CreatedAt.Format("2006-01-02 15:04:05.99")))
			Expect(dbAgent.UpdatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(updateAgent.UpdatedAt.Format("2006-01-02 15:04:05.99")))
		})

	})

	Describe("ProcessRegistry", func() {

		It("returns an error if no db connection is given", func() {
			newAgent := Agent{AgentID: uuid.New()}
			err := newAgent.ProcessRegistration(nil, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should insert a new entry if registration doesn't exist yet", func() {
			// build a registration
			reg := Registration{}; reg.Example()

			// process the registration
			agent := Agent{}
			err := agent.ProcessRegistration(db, &reg.Registration)
			Expect(err).NotTo(HaveOccurred())

			// check
			dbAgent := Agent{AgentID: reg.Sender}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgent.Facts).To(Equal(reg.Payload))
			Expect(dbAgent.Project).To(Equal(reg.Project))
			Expect(dbAgent.Organization).To(Equal(reg.Organization))
		})

		It("should update an existing entry", func() {
			proj := "huhu project"
			org := "huhu organization"
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)

			// save agent with facts
			reg := arc.Registration{Sender: uuid.New(), Project: proj, Organization: org, Payload: facts}
			agent := Agent{}
			agent.FromRegistration(&reg)
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// build a request
			memory_used := 666666
			memory_total := 55555
			newProj := "Miao project"
			newOrg := "Miao organization"
			newFacts := fmt.Sprintf(`{"memory_used": %v, "memory_total": %v}`, memory_used, memory_total)
			newReg := arc.Registration{Sender: agent.AgentID, Project: newProj, Organization: newOrg, Payload: newFacts}

			// process the request
			updateAgent := Agent{}
			err = updateAgent.ProcessRegistration(db, &newReg)
			Expect(err).NotTo(HaveOccurred())

			// check
			dbFacts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": %v, "memory_total": %v, "organization": "%s"}`, proj, memory_used, memory_total, org)
			dbAgent := Agent{AgentID: agent.AgentID}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgent.Facts).To(Equal(dbFacts))
			Expect(dbAgent.Project).To(Equal(newProj))
			Expect(dbAgent.Organization).To(Equal(newOrg))
		})

	})

})
