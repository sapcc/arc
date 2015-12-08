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
			err := dbAgents.Get(db, `@os = "darwin"`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(1))
			Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))

			// get agents with os windows
			err = dbAgents.Get(db, `@os = "windows"`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(2))
			Expect(dbAgents[0].AgentID).To(Equal(agents[2].AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))
		})

	})

	Describe("Get Agents with Authorized and show facts", func() {

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

		It("should return all agents sorted by update and descending", func() {
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
			err = dbAgents.GetAuthorizedAndShowFacts(db, `@os = "windows"`, &authorization, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbAgents)).To(Equal(2))
			Expect(dbAgents[0].AgentID).To(Equal(agent1.AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agent3.AgentID))
		})

		Describe("filter", func() {

			It("should return an error if the filter syntax is wrong", func() {
				// insert facts / agent
				dbAgents := Agents{}
				err := dbAgents.GetAuthorizedAndShowFacts(db, `os =`, &authorization, []string{})
				Expect(err).To(HaveOccurred())
			})

			It("should filter agents", func() {
				// change authorization
				authorization.ProjectId = "miau"

				// create new custom agents
				var (
					facts      = `{"os": "%s", "online": "%s", "project": "miau", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
					os         = []string{"darwin", "windows", "windows"}
					online     = []string{"true", "false", "true"}
					tagsKey    = []string{"landscape", "landscape", "landscape"}
					tagsValue  = []string{"development", "staging", "production"}
					tagsKey2   = []string{"pool", "pool", "pool"}
					tagsValue2 = []string{"green", "green", "blue"}
				)

				agents := Agents{}
				agents.CreateAndSaveAgentExamples(db, 3)
				for i := 0; i < len(agents); i++ {
					currentAgent := agents[i]
					// change facts
					err := json.Unmarshal([]byte(fmt.Sprintf(facts, os[i], online[i])), &currentAgent.Facts)
					Expect(err).NotTo(HaveOccurred())
					currentAgent.Project = "miau"
					err = currentAgent.Update(db)
					Expect(err).NotTo(HaveOccurred())
					// add tags
					err = currentAgent.AddTagAuthorized(db, &authorization, tagsKey[i], tagsValue[i])
					Expect(err).NotTo(HaveOccurred())
					err = currentAgent.AddTagAuthorized(db, &authorization, tagsKey2[i], tagsValue2[i])
					Expect(err).NotTo(HaveOccurred())
				}

				// agent with darwin os
				dbAgents := Agents{}
				err := dbAgents.GetAuthorizedAndShowFacts(db, `@os = "darwin"`, &authorization, []string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(dbAgents)).To(Equal(1))
				Expect(dbAgents[0].AgentID).To(Equal(agents[0].AgentID))

				// agent with windows os
				err = dbAgents.GetAuthorizedAndShowFacts(db, `@os = "windows"`, &authorization, []string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(dbAgents)).To(Equal(2))
				Expect(dbAgents[0].AgentID).To(Equal(agents[2].AgentID))
				Expect(dbAgents[1].AgentID).To(Equal(agents[1].AgentID))

				// agent with windows os and online
				err = dbAgents.GetAuthorizedAndShowFacts(db, `@os = "windows" AND @online = "false"`, &authorization, []string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(dbAgents)).To(Equal(1))
				Expect(dbAgents[0].AgentID).To(Equal(agents[1].AgentID))

				// agent with tag production
				err = dbAgents.GetAuthorizedAndShowFacts(db, `landscape = "staging"`, &authorization, []string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(dbAgents)).To(Equal(1))
				Expect(dbAgents[0].AgentID).To(Equal(agents[1].AgentID))

				// agent more tags
				err = dbAgents.GetAuthorizedAndShowFacts(db, `landscape = "staging" AND (pool = "green" OR pool = "blue")`, &authorization, []string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(dbAgents)).To(Equal(1))
				Expect(dbAgents[0].AgentID).To(Equal(agents[1].AgentID))

				// agent mix tags and facts
				err = dbAgents.GetAuthorizedAndShowFacts(db, `@os = "darwin" OR (landscape = "staging" AND pool = "green")`, &authorization, []string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(dbAgents)).To(Equal(2))
				Expect(dbAgents[0].AgentID).To(Equal(agents[1].AgentID))
				Expect(dbAgents[1].AgentID).To(Equal(agents[0].AgentID))
			})

		})

		Describe("Show facts", func() {

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

		})

		Describe("Authorized", func() {

			It("should return all agents with same authorization project id", func() {
				// add a new agent with a different project
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

	Describe("Get an agent wiht authorization and show facts", func() {

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

		Describe("Show facts", func() {

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

		Describe("Authorized", func() {

			It("should return same project id", func() {
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

		})

	})

	Describe("Delete an agent", func() {

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
			authorization.ProjectId = agent.Project
		})

		It("should delete an agent", func() {
			// delete agent
			err := agent.DeleteAuthorized(db, &authorization)
			Expect(err).NotTo(HaveOccurred())

			// check agent deleted agent
			newAgent := Agent{AgentID: agent.AgentID}
			err = newAgent.Get(db)
			Expect(err).To(HaveOccurred())
		})

		Describe("authorization", func() {

			It("should return an identity authorization error", func() {
				authorization.IdentityStatus = "Something different from Confirmed"

				// delete agent
				err := agent.DeleteAuthorized(db, &authorization)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(auth.IdentityStatusInvalid))
			})

			It("should return a project authorization error", func() {
				authorization.ProjectId = "Some other project"

				// delete agent
				err := agent.DeleteAuthorized(db, &authorization)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(auth.NotAuthorized))
			})

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
				Tags:         JSONB{"cat": "miau"},
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
			// check the tags have not been updated
			Expect(dbAgent.Tags).To(Equal(agent.Tags))
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

	Describe("Tags", func() {

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
			authorization.ProjectId = agent.Project
		})

		Describe("add tag", func() {

			It("returns an error if no db connection is given", func() {
				err := agent.AddTagAuthorized(nil, &authorization, "cat", "miau")
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if agent does not exist", func() {
				newAgent := Agent{AgentID: "non_existing_id"}
				err := newAgent.AddTagAuthorized(db, &authorization, "test", "miau")
				Expect(err).To(HaveOccurred())
			})

			It("should get the agent if the project id is empty", func() {
				// set project empty
				agent.Project = ""
				// should get the data from the agent again
				err := agent.AddTagAuthorized(db, &authorization, "test", "miau")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should save the tag", func() {
				// add tag
				err := agent.AddTagAuthorized(db, &authorization, "test", "miau")
				Expect(err).NotTo(HaveOccurred())

				// check
				dbAgent := Agent{AgentID: agent.AgentID}
				err = dbAgent.Get(db)
				// conver to JSON string
				tags, err := json.Marshal(dbAgent.Tags)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(tags)).To(Equal(`{"test":"miau"}`))
			})

			It("should override the tag if exists", func() {
				// add tag
				err := agent.AddTagAuthorized(db, &authorization, "test", "miau")
				Expect(err).NotTo(HaveOccurred())

				// override tag
				err = agent.AddTagAuthorized(db, &authorization, "test", "miau, miau")
				Expect(err).NotTo(HaveOccurred())

				// check
				dbAgent := Agent{AgentID: agent.AgentID}
				err = dbAgent.Get(db)
				// conver to JSON string
				tags, err := json.Marshal(dbAgent.Tags)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(tags)).To(Equal(`{"test":"miau, miau"}`))
			})

			Describe("authorization", func() {

				It("should return an identity authorization error", func() {
					authorization.IdentityStatus = "Something different from Confirmed"

					// add tag
					err := agent.AddTagAuthorized(db, &authorization, "test", "miau")
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(auth.IdentityStatusInvalid))
				})

				It("should return a project authorization error", func() {
					authorization.ProjectId = "Some other project"

					err := agent.AddTagAuthorized(db, &authorization, "test", "miau")
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(auth.NotAuthorized))
				})

			})

		})

		Describe("remove tag", func() {

			It("returns an error if no db connection is given", func() {
				err := agent.DeleteTagAuthorized(nil, &authorization, "cat")
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if agent does not exist", func() {
				newAgent := Agent{AgentID: "non_existing_id"}
				err := newAgent.DeleteTagAuthorized(db, &authorization, "dog")
				Expect(err).To(HaveOccurred())
			})

			// no error handling possible because of the row will allways updated
			It("should not return an error if the tag does not exist", func() {
				err := agent.DeleteTagAuthorized(db, &authorization, "non_existing_tag")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should get the agent if the project id is empty", func() {
				// set project empty
				agent.Project = ""
				// should get the data from the agent again
				err := agent.DeleteTagAuthorized(db, &authorization, "dog")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should save the tag", func() {
				// add tag
				err := agent.AddTagAuthorized(db, &authorization, "cat", "miau")
				Expect(err).NotTo(HaveOccurred())

				err = agent.AddTagAuthorized(db, &authorization, "dog", "bup")
				Expect(err).NotTo(HaveOccurred())

				// remove tag
				err = agent.DeleteTagAuthorized(db, &authorization, "dog")
				Expect(err).NotTo(HaveOccurred())

				// check
				dbAgent := Agent{AgentID: agent.AgentID}
				err = dbAgent.Get(db)
				// conver to JSON string
				tags, err := json.Marshal(dbAgent.Tags)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(tags)).To(Equal(`{"cat":"miau"}`))
			})

			Describe("authorization", func() {

				It("should return an identity authorization error", func() {
					authorization.IdentityStatus = "Something different from Confirmed"

					// add tag
					err := agent.DeleteTagAuthorized(db, &authorization, "dog")
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(auth.IdentityStatusInvalid))
				})

				It("should return a project authorization error", func() {
					authorization.ProjectId = "Some other project"

					err := agent.DeleteTagAuthorized(db, &authorization, "dog")
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(auth.NotAuthorized))
				})

			})

		})

	})

})
