// +build integration

package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"			
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			agent := Agent{}; agent.Example()
			err := agent.Get(nil)
			Expect(err).To(HaveOccurred())
		})
		
		It("returns an error if no agent found", func() {
			agent := Agent{}; agent.Example()
			err := agent.Get(db)
			Expect(err).To(HaveOccurred())
		})
		
		It("should return the agent", func() {			
			// insert facts / agent
			agent := Agent{}; agent.Example()
			var lastInsertId string
			err := db.QueryRow(InsertFactQuery, agent.AgentID, agent.Project, agent.Organization, "{}", agent.CreatedAt, agent.UpdatedAt).Scan(&lastInsertId);
			Expect(err).NotTo(HaveOccurred())
			
			// get agent
			newAgent := Agent{AgentID: agent.AgentID}
			err = newAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(agent.AgentID).To(Equal(newAgent.AgentID))
			Expect(agent.CreatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(newAgent.CreatedAt.Format("2006-01-02 15:04:05.99")))
			Expect(agent.UpdatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(newAgent.UpdatedAt.Format("2006-01-02 15:04:05.99")))
		})
		
	})

})
