// +build integration

package models_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/api-server/db"			
	. "gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"code.google.com/p/go-uuid/uuid"
	
	"time"

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
			agent1 := Agent{AgentID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
			agent2 := Agent{AgentID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
			agent3 := Agent{AgentID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()}						
			agents := []Agent{agent1, agent2, agent3}
			
			// insert 3 agents
			for i := 0; i < len(agents); i++ {
				agent := agents[i]
				var lastInsertId string
				err := db.QueryRow(InsertFactQuery, agent.AgentID, "{}", agent.CreatedAt, agent.UpdatedAt).Scan(&lastInsertId);
				Expect(err).NotTo(HaveOccurred())			
			}
			
			// insert facts / agent
			dbAgents := Agents{}
			err := dbAgents.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbAgents[0].AgentID).To(Equal(agent1.AgentID))
			Expect(dbAgents[1].AgentID).To(Equal(agent2.AgentID))
			Expect(dbAgents[2].AgentID).To(Equal(agent3.AgentID))			
		})
		
	})
	
})

var _ = Describe("Agent", func() {

	Describe("Get", func() {

		It("returns an error if no db connection is given", func() {
			newAgent := Agent{AgentID: uuid.New()}
			err := newAgent.Get(nil)
			Expect(err).To(HaveOccurred())
		})
		
		It("returns an error if no agent found", func() {
			newAgent := Agent{AgentID: uuid.New()}
			err := newAgent.Get(db)
			Expect(err).To(HaveOccurred())
		})
		
		It("should return the agent", func() {
			agent_id := uuid.New()
			
			// insert facts / agent
			agent := Agent{AgentID: agent_id, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			var lastInsertId string
			err := db.QueryRow(InsertFactQuery, agent_id, "{}", agent.CreatedAt, agent.UpdatedAt).Scan(&lastInsertId);
			Expect(err).NotTo(HaveOccurred())
			
			// get agent
			newAgent := Agent{AgentID: agent_id}
			err = newAgent.Get(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(agent.AgentID).To(Equal(newAgent.AgentID))
			Expect(agent.CreatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(newAgent.CreatedAt.Format("2006-01-02 15:04:05.99")))
			Expect(agent.UpdatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(newAgent.UpdatedAt.Format("2006-01-02 15:04:05.99")))
		})
		
	})

})
