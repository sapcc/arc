// +build integration

package models_test

import (
// . "github.com/onsi/ginkgo"
// . "github.com/onsi/gomega"
// "time"
//
// auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
// . "gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

// var _ = Describe("Tag", func() {
//
// 	var (
// 		agent         = Agent{}
// 		tag           = Tag{}
// 		authorization = auth.Authorization{}
// 	)
//
// 	JustBeforeEach(func() {
// 		agent.Example()
// 		err := agent.Save(db)
// 		Expect(err).NotTo(HaveOccurred())
//
// 		authorization.IdentityStatus = "Confirmed"
// 		authorization.UserId = "userID"
// 		authorization.ProjectId = agent.Project
//
// 		tag.AgentID = agent.AgentID
// 		tag.Project = agent.Project
// 		tag.Value = "tag_awesome"
// 		tag.CreatedAt = time.Now()
// 		err = tag.SaveAuthorized(db, &authorization)
// 		Expect(err).NotTo(HaveOccurred())
// 	})
//
// 	Describe("Get", func() {
//
// 		It("returns an error if no db connection is given", func() {
// 			tag := Tag{AgentID: agent.AgentID, Project: agent.Project, Value: "some_tag", CreatedAt: time.Now()}
// 			err := tag.SaveAuthorized(nil, &authorization)
// 			Expect(err).To(HaveOccurred())
// 		})
//
// 		It("should not save tag if agent does not exist", func() {
// 			tag := Tag{AgentID: "not_existing_agent", Project: agent.Project, Value: "some_tag", CreatedAt: time.Now()}
// 			err := tag.SaveAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 		})
//
// 		It("should remove the tag if the agent gets deleted", func() {
// 			agent.DeleteAuthorized(db, &authorization)
//
// 			dbTag := Tag{AgentID: agent.AgentID, Value: "tag_awesome"}
// 			err := dbTag.GetAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 		})
//
// 		It("should save and get the tag the tag with same project id", func() {
// 			newTag := Tag{AgentID: agent.AgentID, Project: agent.Project, Value: "tag_cool", CreatedAt: time.Now()}
// 			err := newTag.SaveAuthorized(db, &authorization)
// 			Expect(err).NotTo(HaveOccurred())
//
// 			dbTag := Tag{AgentID: agent.AgentID, Value: "tag_cool"}
// 			err = dbTag.GetAuthorized(db, &authorization)
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(dbTag.AgentID).To(Equal(newTag.AgentID))
// 			Expect(dbTag.Value).To(Equal(newTag.Value))
// 			Expect(dbTag.CreatedAt.Format("2006-01-02 15:04:05.99")).To(Equal(newTag.CreatedAt.Format("2006-01-02 15:04:05.99")))
// 		})
//
// 		It("should return an identity authorization error", func() {
// 			authorization.IdentityStatus = "Something different from Confirmed"
//
// 			dbTag := Tag{AgentID: agent.AgentID, Value: "tag_awesome"}
// 			err := dbTag.GetAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.IdentityStatusInvalid))
// 		})
//
// 		It("should return a project authorization error", func() {
// 			authorization.ProjectId = "Some other project"
//
// 			dbTag := Tag{AgentID: agent.AgentID, Value: "tag_awesome"}
// 			err := dbTag.GetAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.NotAuthorized))
// 		})
//
// 	})
//
// 	Describe("Save authorized errors", func() {
//
// 		It("should return an identity authorization error", func() {
// 			authorization.IdentityStatus = "Something different from Confirmed"
//
// 			newTag := Tag{AgentID: agent.AgentID, Project: agent.Project, Value: "tag_cool", CreatedAt: time.Now()}
// 			err := newTag.SaveAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.IdentityStatusInvalid))
// 		})
//
// 		It("should return a project authorization error", func() {
// 			authorization.ProjectId = "Some other project"
//
// 			newTag := Tag{AgentID: agent.AgentID, Project: agent.Project, Value: "tag_cool", CreatedAt: time.Now()}
// 			err := newTag.SaveAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.NotAuthorized))
// 		})
//
// 	})
//
// 	Describe("GetByAgentIdAuthorized", func() {
//
// 		var (
// 			newAgent = Agent{}
// 		)
//
// 		JustBeforeEach(func() {
// 			newAgent.Example()
// 			err := newAgent.Save(db)
// 			Expect(err).NotTo(HaveOccurred())
// 		})
//
// 		It("returns an error if no db connection is given", func() {
// 			dbTags := Tags{}
// 			err := dbTags.GetByAgentIdAuthorized(nil, &authorization, newAgent.AgentID)
// 			Expect(err).To(HaveOccurred())
// 		})
//
// 		It("returns an error if agent does not exist", func() {
// 			dbTags := Tags{}
// 			err := dbTags.GetByAgentIdAuthorized(db, &authorization, "non_existing_id")
// 			Expect(err).To(HaveOccurred())
// 		})
//
// 		It("should return all tags from an agent with same project id", func() {
// 			newTag := Tag{AgentID: newAgent.AgentID, Project: agent.Project, Value: "tag_miau", CreatedAt: time.Now()}
// 			err := newTag.SaveAuthorized(db, &authorization)
// 			Expect(err).NotTo(HaveOccurred())
//
// 			newTag2 := Tag{AgentID: newAgent.AgentID, Project: agent.Project, Value: "tag_bup", CreatedAt: time.Now()}
// 			err = newTag2.SaveAuthorized(db, &authorization)
// 			Expect(err).NotTo(HaveOccurred())
//
// 			dbTags := Tags{}
// 			err = dbTags.GetByAgentIdAuthorized(db, &authorization, newAgent.AgentID)
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(len(dbTags)).To(Equal(2))
// 			Expect(dbTags[0].AgentID).To(Equal(newTag2.AgentID))
// 			Expect(dbTags[0].Value).To(Equal(newTag2.Value))
// 			Expect(dbTags[1].AgentID).To(Equal(newTag.AgentID))
// 			Expect(dbTags[1].Value).To(Equal(newTag.Value))
// 		})
//
// 		It("should return an identity authorization error", func() {
// 			authorization.IdentityStatus = "Something different from Confirmed"
//
// 			dbTags := Tags{}
// 			err := dbTags.GetByAgentIdAuthorized(db, &authorization, newAgent.AgentID)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.IdentityStatusInvalid))
// 		})
//
// 		It("should return a project authorization error", func() {
// 			authorization.ProjectId = "Some other project"
//
// 			dbTags := Tags{}
// 			err := dbTags.GetByAgentIdAuthorized(db, &authorization, newAgent.AgentID)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.NotAuthorized))
// 		})
//
// 	})
//
// 	Describe("GetByValueAuthorized", func() {
//
// 		var (
// 			newAgent = Agent{}
// 		)
//
// 		JustBeforeEach(func() {
// 			newAgent.Example()
// 			err := newAgent.Save(db)
// 			Expect(err).NotTo(HaveOccurred())
//
// 			newTag := Tag{AgentID: newAgent.AgentID, Project: agent.Project, Value: "tag_awesome", CreatedAt: time.Now()}
// 			err = newTag.SaveAuthorized(db, &authorization)
// 			Expect(err).NotTo(HaveOccurred())
// 		})
//
// 		It("returns an error if no db connection is given", func() {
// 			dbTags := Tags{}
// 			err := dbTags.GetByValueAuthorized(nil, &authorization, "tag_awesome")
// 			Expect(err).To(HaveOccurred())
// 		})
//
// 		It("returns an empty array if tag value does not exist", func() {
// 			dbTags := Tags{}
// 			err := dbTags.GetByValueAuthorized(db, &authorization, "non_existing_tag")
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(len(dbTags)).To(Equal(0))
// 		})
//
// 		It("should return all tags from with same value and agents project id", func() {
// 			dbTags := Tags{}
// 			err := dbTags.GetByValueAuthorized(db, &authorization, "tag_awesome")
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(len(dbTags)).To(Equal(2))
// 			Expect(dbTags[0].AgentID).To(Equal(newAgent.AgentID))
// 			Expect(dbTags[0].Value).To(Equal("tag_awesome"))
// 			Expect(dbTags[1].AgentID).To(Equal(agent.AgentID))
// 			Expect(dbTags[1].Value).To(Equal("tag_awesome"))
// 		})
//
// 		It("should return an identity authorization error", func() {
// 			authorization.IdentityStatus = "Something different from Confirmed"
//
// 			dbTags := Tags{}
// 			err := dbTags.GetByValueAuthorized(db, &authorization, "tag_awesome")
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.IdentityStatusInvalid))
// 		})
//
// 		It("should return a project authorization error", func() {
// 			authorization.ProjectId = "Some other project"
//
// 			dbTags := Tags{}
// 			err := dbTags.GetByValueAuthorized(db, &authorization, "tag_awesome")
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(len(dbTags)).To(Equal(0))
// 		})
//
// 	})
//
// 	Describe("DeleteAuthorized", func() {
//
// 		It("should return an error if tag not found", func() {
// 			newTag := Tag{AgentID: "fake_agent", Project: "fake_project", Value: "tag_awesome", CreatedAt: time.Now()}
// 			err := newTag.DeleteAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 		})
//
// 		It("should remove a tag", func() {
// 			err := tag.DeleteAuthorized(db, &authorization)
// 			Expect(err).NotTo(HaveOccurred())
//
// 			// check if we have 0 tags
// 			dbTags := Tags{}
// 			err = dbTags.GetByAgentIdAuthorized(db, &authorization, agent.AgentID)
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(len(dbTags)).To(Equal(0))
// 		})
//
// 		It("should return an identity authorization error", func() {
// 			authorization.IdentityStatus = "Something different from Confirmed"
//
// 			err := tag.DeleteAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.IdentityStatusInvalid))
// 		})
//
// 		It("should return a project authorization error", func() {
// 			authorization.ProjectId = "Some other project"
//
// 			err := tag.DeleteAuthorized(db, &authorization)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.NotAuthorized))
// 		})
//
// 	})
//
// 	Describe("ProcessAgentTagsData", func() {
//
// 		var (
// 			newAgent = Agent{}
// 		)
//
// 		JustBeforeEach(func() {
// 			newAgent.Example()
// 			err := newAgent.Save(db)
// 			Expect(err).NotTo(HaveOccurred())
// 		})
//
// 		It("returns an error if no db connection is given", func() {
// 			err := ProcessAgentTagsData(nil, &authorization, newAgent, []byte("test"))
// 			Expect(err).To(HaveOccurred())
// 		})
//
// 		It("should save all tags for the given agent", func() {
// 			err := ProcessAgentTagsData(db, &authorization, newAgent, []byte("test , test2,test3"))
// 			Expect(err).NotTo(HaveOccurred())
//
// 			dbTags := Tags{}
// 			err = dbTags.GetByAgentIdAuthorized(db, &authorization, newAgent.AgentID)
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(len(dbTags)).To(Equal(3))
// 			Expect(dbTags[0].AgentID).To(Equal(newAgent.AgentID))
// 			Expect(dbTags[0].Value).To(Equal("test3"))
// 			Expect(dbTags[1].AgentID).To(Equal(newAgent.AgentID))
// 			Expect(dbTags[1].Value).To(Equal("test2"))
// 			Expect(dbTags[2].AgentID).To(Equal(newAgent.AgentID))
// 			Expect(dbTags[2].Value).To(Equal("test"))
// 		})
//
// 		It("should do nothing when primary key violation happens", func() {
// 			err := ProcessAgentTagsData(db, &authorization, newAgent, []byte("test , test2,test3"))
// 			Expect(err).NotTo(HaveOccurred())
//
// 			// we save the same + 1 new
// 			err = ProcessAgentTagsData(db, &authorization, newAgent, []byte("test , test2,test3, test4"))
// 			Expect(err).NotTo(HaveOccurred())
// 			dbTags := Tags{}
// 			err = dbTags.GetByAgentIdAuthorized(db, &authorization, newAgent.AgentID)
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(len(dbTags)).To(Equal(4))
// 			Expect(dbTags[0].AgentID).To(Equal(newAgent.AgentID))
// 			Expect(dbTags[0].Value).To(Equal("test4"))
// 		})
//
// 		It("should return an identity authorization error", func() {
// 			authorization.IdentityStatus = "Something different from Confirmed"
//
// 			err := ProcessAgentTagsData(db, &authorization, newAgent, []byte("test , test2,test3"))
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.IdentityStatusInvalid))
// 		})
//
// 		It("should return a project authorization error", func() {
// 			authorization.ProjectId = "Some other project"
//
// 			err := ProcessAgentTagsData(db, &authorization, newAgent, []byte("test , test2,test3"))
// 			Expect(err).To(HaveOccurred())
// 			Expect(err).To(Equal(auth.NotAuthorized))
// 		})
//
// 	})
//
// })
