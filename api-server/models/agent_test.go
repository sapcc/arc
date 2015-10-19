// +build integration

package models_test

import (
	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"

	"fmt"
	"time"
)

var _ = Describe("Agents", func() {

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			agents := Agents{}
			err := agents.Get(nil, "")
			Expect(err).To(HaveOccurred())
		})

		It("should return all agents", func() {
			agents := Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)

			// insert facts / agent
			dbAgents := Agents{}
			err := dbAgents.Get(db, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))
			Expect(dbAgents[2].AgentID).To(Equal(agents[2].AgentID))
		})

		It("should return an error if the filter syntax is wrong", func() {
			// insert facts / agent
			dbAgents := Agents{}
			err := dbAgents.Get(db, `os =`)
			Expect(err).To(HaveOccurred())
		})

		It("should return all agents filtered", func() {
			facts := `{"os": "%s", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
			os := []string{"darwin", "windows", "windows"}
			agents := Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)
			for i := 0; i < len(agents); i++ {
				agents[i].Facts = fmt.Sprintf(facts, os[i])
				err := agents[i].Update(db)
				Expect(err).NotTo(HaveOccurred())
			}

			// insert facts / agent
			dbAgents := Agents{}
			err := dbAgents.Get(db, `os = "darwin"`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(1))
			Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))

			err = dbAgents.Get(db, `os = "windows"`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(2))
			Expect(dbAgents[0].AgentID).To(Equal(agents[1].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[2].AgentID))
		})

	})

	Describe("GetAuthorized", func() {

		var (
			agents        = Agents{}
			authorization = auth.Authorization{}
		)

		JustBeforeEach(func() {
			agents.CreateAndSaveAgentExamples(db, 3)
			authorization.IdentityStatus = "Confirmed"
			authorization.UserId = "userID"
			authorization.ProjectId = "test-project"
		})

		It("returns an error if no db connection is given", func() {
			agents := Agents{}
			err := agents.GetAuthorized(nil, "", &authorization)
			Expect(err).To(HaveOccurred())
		})

		It("should return all agents where with same project", func() {
			// add a new agent
			agent := Agent{}
			agent.Example()
			agent.Project = "miau"
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// change authorization
			authorization.ProjectId = "miau"

			// insert facts / agent
			dbAgents := Agents{}
			err = dbAgents.GetAuthorized(db, "", &authorization)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(1))
			Expect(dbAgents[0].Project).To(Equal(authorization.ProjectId))
		})

		It("should return an error if the filter syntax is wrong", func() {
			// insert facts / agent
			dbAgents := Agents{}
			err := dbAgents.GetAuthorized(db, `os =`, &authorization)
			Expect(err).To(HaveOccurred())
		})

		It("should return all agents filtered with same project", func() {
			facts := `{"os": "%s", "online": true, "project": "miau", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
			os := []string{"darwin", "windows", "windows"}
			agents := Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)
			for i := 0; i < len(agents); i++ {
				agents[i].Facts = fmt.Sprintf(facts, os[i])
				agents[i].Project = "miau"
				err := agents[i].Update(db)
				Expect(err).NotTo(HaveOccurred())
			}

			// change authorization
			authorization.ProjectId = "miau"

			// insert facts / agent
			dbAgents := Agents{}
			err := dbAgents.GetAuthorized(db, `os = "darwin"`, &authorization)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(1))
			Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))

			err = dbAgents.GetAuthorized(db, `os = "windows"`, &authorization)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(2))
			Expect(dbAgents[0].AgentID).To(Equal(agents[1].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[2].AgentID))
		})

		It("should return an identity authorization error", func() {
			authorization.IdentityStatus = "Something different from Confirmed"

			dbAgents := Agents{}
			err := dbAgents.GetAuthorized(db, "", &authorization)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(auth.IdentityStatusInvalid))
		})

		It("should return a project authorization error", func() {
			authorization.ProjectId = "Some other project"

			dbAgents := Agents{}
			err := dbAgents.GetAuthorized(db, "", &authorization)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(0))
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

	Describe("Delete authorized", func() {

		var (
			agent = Agent{}
		)

		JustBeforeEach(func() {
			// insert facts / agent
			agent.Example()
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an identity authorization error", func() {
			authorization := auth.Authorization{
				IdentityStatus: "Something different from Confirmed",
				UserId:         "userID",
				ProjectId:      agent.Project,
			}

			// delete agent
			err := agent.DeleteAuthorized(db, &authorization)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(auth.IdentityStatusInvalid))
		})

		It("should return a project authorization error", func() {
			authorization := auth.Authorization{
				IdentityStatus: "Confirmed",
				UserId:         "userID",
				ProjectId:      "Some other project",
			}

			// delete agent
			err := agent.DeleteAuthorized(db, &authorization)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(auth.NotAuthorized))
		})

		It("should delete an agent", func() {
			authorization := auth.Authorization{
				IdentityStatus: "Confirmed",
				UserId:         "userID",
				ProjectId:      agent.Project,
			}

			// delete agent
			err := agent.DeleteAuthorized(db, &authorization)
			Expect(err).NotTo(HaveOccurred())

			// check agent deleted agent
			newAgent := Agent{AgentID: agent.AgentID}
			err = newAgent.Get(db)
			Expect(err).To(HaveOccurred())
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
			err := ProcessRegistration(nil, nil, "darwin", true)
			Expect(err).To(HaveOccurred())
		})

		It("should insert a new entry if registration doesn't exist yet", func() {
			// build a registration
			reg := Registration{}
			reg.Example()

			// process the registration
			err := ProcessRegistration(db, &reg.Registration, "darwin", true)
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
			agentId := "darwin"
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)

			// save agent with facts
			reg := arc.Registration{RegistrationID: uuid.New(), Sender: uuid.New(), Project: proj, Organization: org, Payload: facts}
			agent := Agent{}
			agent.FromRegistration(&reg, agentId)
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// build a request
			memory_used := 666666
			memory_total := 55555
			newProj := "Miao project"
			newOrg := "Miao organization"
			newFacts := fmt.Sprintf(`{"memory_used": %v, "memory_total": %v}`, memory_used, memory_total)
			newReg := arc.Registration{RegistrationID: uuid.New(), Sender: agent.AgentID, Project: newProj, Organization: newOrg, Payload: newFacts}

			// process the request
			err = ProcessRegistration(db, &newReg, agentId, true)
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

		It("should check concurrency safe", func() {
			proj := "huhu project"
			org := "huhu organization"
			agentId := "darwin"
			registrationId := uuid.New()
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)

			// save agent with facts
			reg := arc.Registration{RegistrationID: registrationId, Sender: uuid.New(), Project: proj, Organization: org, Payload: facts}
			err := ProcessRegistration(db, &reg, agentId, true)
			Expect(err).NotTo(HaveOccurred())

			// build a new registration with same id
			memory_used := 666666
			memory_total := 55555
			newProj := "Miao project"
			newOrg := "Miao organization"
			newFacts := fmt.Sprintf(`{"memory_used": %v, "memory_total": %v}`, memory_used, memory_total)
			newReg := arc.Registration{RegistrationID: registrationId, Sender: agentId, Project: newProj, Organization: newOrg, Payload: newFacts}

			// process the request
			err = ProcessRegistration(db, &newReg, agentId, true)
			Expect(err).To(Equal(RegistrationExistsError))
		})

		It("should NOT check concurrency safe", func() {
			proj := "huhu project"
			org := "huhu organization"
			agentId := "darwin"
			registrationId := uuid.New()
			facts := fmt.Sprintf(`{"os": "darwin", "online": true, "project": "%s", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "%s"}`, proj, org)
			sender := uuid.New()

			// save agent with facts
			reg := arc.Registration{RegistrationID: registrationId, Sender: sender, Project: proj, Organization: org, Payload: facts}
			err := ProcessRegistration(db, &reg, agentId, false)
			Expect(err).NotTo(HaveOccurred())

			// build a new registration with same id
			memory_used := 666666
			memory_total := 55555
			newProj := "Miao project"
			newOrg := "Miao organization"
			newFacts := fmt.Sprintf(`{"memory_used": %v, "memory_total": %v}`, memory_used, memory_total)
			newReg := arc.Registration{RegistrationID: registrationId, Sender: sender, Project: newProj, Organization: newOrg, Payload: newFacts}

			// process the request
			err = ProcessRegistration(db, &newReg, agentId, false)
			Expect(err).NotTo(HaveOccurred()) // it is updated twice
		})

	})

})
