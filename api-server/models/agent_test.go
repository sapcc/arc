// +build integration

package models_test

import (
	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	arc "gitHub.***REMOVED***/monsoon/arc/arc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"

	"encoding/json"
	"fmt"
	"reflect"
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
			// check that the agents are sorted descending
			Expect(dbAgents[0].AgentID).To(Equal(agents[2].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))
			Expect(dbAgents[2].AgentID).To(Equal(agents[0].AgentID))
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

			// create 3 examples
			agents := Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)
			// update facts in the examples
			for i := 0; i < len(agents); i++ {
				currentAgent := agents[i]
				if err := json.Unmarshal([]byte(fmt.Sprintf(facts, os[i])), &currentAgent.Facts); err != nil {
					Expect(err).NotTo(HaveOccurred())
				}
				err := currentAgent.Update(db)
				Expect(err).NotTo(HaveOccurred())
			}

			// get agents with os darwin
			dbAgents := Agents{}
			err := dbAgents.Get(db, `os = "darwin"`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(1))
			Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))

			// get agents with os windows
			err = dbAgents.Get(db, `os = "windows"`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(2))
			Expect(dbAgents[0].AgentID).To(Equal(agents[2].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))
		})

	})

	Describe("Authorized and show facts", func() {

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
			err := agents.GetAuthorizedAndShowFacts(nil, "", &authorization, []string{})
			Expect(err).To(HaveOccurred())
		})

		It("should return all agents with same authorization project id", func() {
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
			err = dbAgents.GetAuthorizedAndShowFacts(db, "", &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(1))
			Expect(dbAgents[0].Project).To(Equal(authorization.ProjectId))
		})

		It("should return an error if the filter syntax is wrong", func() {
			// insert facts / agent
			dbAgents := Agents{}
			err := dbAgents.GetAuthorizedAndShowFacts(db, `os =`, &authorization, []string{})
			Expect(err).To(HaveOccurred())
		})

		It("should return all agents filtered by os", func() {
			facts := `{"os": "%s", "online": true, "project": "miau", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
			os := []string{"darwin", "windows", "windows"}
			agents := Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)
			for i := 0; i < len(agents); i++ {
				currentAgent := agents[i]
				if err := json.Unmarshal([]byte(fmt.Sprintf(facts, os[i])), &currentAgent.Facts); err != nil {
					Expect(err).NotTo(HaveOccurred())
				}
				currentAgent.Project = "miau"
				err := currentAgent.Update(db)
				Expect(err).NotTo(HaveOccurred())
			}

			// change authorization
			authorization.ProjectId = "miau"

			// agent with darwin os
			dbAgents := Agents{}
			err := dbAgents.GetAuthorizedAndShowFacts(db, `os = "darwin"`, &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(1))
			Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))

			// agent with windows os
			err = dbAgents.GetAuthorizedAndShowFacts(db, `os = "windows"`, &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(2))
			Expect(dbAgents[0].AgentID).To(Equal(agents[2].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))
		})

		It("should return an identity authorization error", func() {
			authorization.IdentityStatus = "Something different from Confirmed"

			dbAgents := Agents{}
			err := dbAgents.GetAuthorizedAndShowFacts(db, "", &authorization, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(auth.IdentityStatusInvalid))
		})

		It("should return a project authorization error", func() {
			authorization.ProjectId = "Some other project"

			dbAgents := Agents{}
			err := dbAgents.GetAuthorizedAndShowFacts(db, "", &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(0))
		})

		It("should return all agents with the given facts", func() {
			// change authorization
			authorization.ProjectId = "test-project"

			// get agents with existing facts
			dbAgents := Agents{}
			err := dbAgents.GetAuthorizedAndShowFacts(db, "", &authorization, []string{"os", "online"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(3))
			for i := 0; i < len(dbAgents); i++ {
				currentAgent := dbAgents[i]
				Expect(len(currentAgent.Facts)).To(Equal(2))
				_, ok := currentAgent.Facts["os"]
				Expect(ok).To(Equal(true))
				_, ok = currentAgent.Facts["online"]
				Expect(ok).To(Equal(true))
			}

			// get agents with non existing facts
			err = dbAgents.GetAuthorizedAndShowFacts(db, "", &authorization, []string{"os", "bup"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(3))
			for i := 0; i < len(dbAgents); i++ {
				currentAgent := dbAgents[i]
				Expect(len(currentAgent.Facts)).To(Equal(1))
			}

			// get agent with all existing agents
			err = dbAgents.GetAuthorizedAndShowFacts(db, "", &authorization, []string{"all"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(3))
			for i := 0; i < len(dbAgents); i++ {
				currentAgent := dbAgents[i]
				Expect(len(currentAgent.Facts)).To(BeNumerically(">=", 10))
			}

			// get agents with no facts
			err = dbAgents.GetAuthorizedAndShowFacts(db, "", &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(3))
			for i := 0; i < len(dbAgents); i++ {
				currentAgent := dbAgents[i]
				Expect(len(currentAgent.Facts)).To(Equal(0))
			}
		})

		It("should return all agents sorted descending", func() {
			facts := `{"os": "%s", "online": true, "project": "miau", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
			agents := Agents{}
			agents.CreateAndSaveAgentExamples(db, 3)

			agent1 := agents[0]
			json.Unmarshal([]byte(fmt.Sprintf(facts, "windows")), &agent1.Facts)
			agent1.Project = "miau"
			agent1.UpdatedAt = time.Now().Add(-5 * time.Minute)
			err := agent1.Update(db)
			Expect(err).NotTo(HaveOccurred())

			agent2 := agents[1]
			agent2.Project = "miau"
			json.Unmarshal([]byte(fmt.Sprintf(facts, "darwin")), &agent2.Facts)
			agent2.UpdatedAt = time.Now().Add(-30 * time.Minute)
			err = agent2.Update(db)
			Expect(err).NotTo(HaveOccurred())

			agent3 := agents[1]
			agent3.Project = "miau"
			json.Unmarshal([]byte(fmt.Sprintf(facts, "windows")), &agent3.Facts)
			agent3.UpdatedAt = time.Now().Add(-20 * time.Minute)
			err = agent3.Update(db)
			Expect(err).NotTo(HaveOccurred())

			// change authorization
			authorization.ProjectId = "miau"

			dbAgents := Agents{}
			err = dbAgents.GetAuthorizedAndShowFacts(db, `os = "windows"`, &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(2))
			Expect(dbAgents[0].AgentID).To(Equal(agent1.AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agent3.AgentID))
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

	Describe("Get authorized and show facts", func() {

		var (
			agent         = Agent{}
			authorization = auth.Authorization{}
		)

		JustBeforeEach(func() {
			// insert facts / agent
			agent.Example()
			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())
			// reset authorization
			authorization.IdentityStatus = "Confirmed"
			authorization.UserId = "userID"
			authorization.ProjectId = "test-project"
		})

		It("should return an agent with same authorization project id", func() {
			// add a new agent
			newAgent := Agent{}
			newAgent.Example()
			newAgent.Project = "miau"
			err := newAgent.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// change authorization
			authorization.ProjectId = "miau"

			// insert facts / agent
			dbAgent := Agent{AgentID: newAgent.AgentID}
			err = dbAgent.GetAuthorizedAndShowFacts(db, &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgent.Project).To(Equal(authorization.ProjectId))
		})

		It("should return an identity authorization error", func() {
			authorization.IdentityStatus = "Something different from Confirmed"

			dbAgent := Agent{AgentID: agent.AgentID}
			err := dbAgent.GetAuthorizedAndShowFacts(db, &authorization, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(auth.IdentityStatusInvalid))
		})

		It("should return a project authorization error", func() {
			authorization.ProjectId = "Some other project"

			dbAgent := Agent{AgentID: agent.AgentID}
			err := dbAgent.GetAuthorizedAndShowFacts(db, &authorization, []string{})
			Expect(err).To(HaveOccurred())
		})

		It("should return all agents with the given facts", func() {
			// authorization
			authorization.ProjectId = "test-project"

			// get agent with existing facts
			dbAgent := Agent{AgentID: agent.AgentID}
			err := dbAgent.GetAuthorizedAndShowFacts(db, &authorization, []string{"os", "online"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgent.Facts)).To(Equal(2))
			_, ok := dbAgent.Facts["os"]
			Expect(ok).To(Equal(true))
			_, ok = dbAgent.Facts["online"]
			Expect(ok).To(Equal(true))

			// get agent with non existing facts
			err = dbAgent.GetAuthorizedAndShowFacts(db, &authorization, []string{"os", "bup"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgent.Facts)).To(Equal(1))

			// get agent with all existing agents
			err = dbAgent.GetAuthorizedAndShowFacts(db, &authorization, []string{"all"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgent.Facts)).To(BeNumerically(">=", 10))

			// get agent with no facts
			err = dbAgent.GetAuthorizedAndShowFacts(db, &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgent.Facts)).To(Equal(0))
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
				CreatedAt:    time.Now().Add((-5) * time.Minute),
				UpdatedAt:    time.Now().Add((-5) * time.Minute),
			}
			if err := json.Unmarshal([]byte(fmt.Sprintf(facts, "huhu project", 123456789, 987654321, "huhu organization")), &agent.Facts); err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			err := agent.Save(db)
			Expect(err).NotTo(HaveOccurred())

			// update same agent
			newProj := "Miau project"
			newOrg := "Miau organization"
			memory_used := 666666
			memory_total := 55555
			updateAgent := Agent{
				AgentID:      agent.AgentID,
				Project:      newProj,
				Organization: newOrg,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if err := json.Unmarshal([]byte(fmt.Sprintf(`{"memory_used": %v, "memory_total": %v, "project": "%s", "organization": "%s"}`, memory_used, memory_total, newProj, newOrg)), &updateAgent.Facts); err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			err = updateAgent.Update(db)
			Expect(err).NotTo(HaveOccurred())

			// check
			dbAgent := Agent{AgentID: agent.AgentID}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgent.Facts["project"]).To(Equal(newProj))
			Expect(dbAgent.Facts["memory_used"]).To(Equal(float64(memory_used)))
			Expect(dbAgent.Facts["memory_total"]).To(Equal(float64(memory_total)))
			Expect(dbAgent.Facts["organization"]).To(Equal(newOrg))
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

			// build test facts
			checkFacts := JSONBfromString(reg.Payload)

			// check
			dbAgent := Agent{AgentID: reg.Sender}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			eq := reflect.DeepEqual(dbAgent.Facts, checkFacts)
			Expect(eq).To(Equal(true))
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
			checkFacts := JSONBfromString(dbFacts)

			dbAgent := Agent{AgentID: agent.AgentID}
			err = dbAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			eq := reflect.DeepEqual(dbAgent.Facts, checkFacts)
			Expect(eq).To(Equal(true))
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
