package main

import (
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

var jobs = models.Jobs{
	models.Job{
		Request: arc.Request{
			Version: 666,
			Sender: "me",
			RequestID: "123456789",
			To: "you",
			Timeout: 333,
			Agent: "007",
			Action: "hhmm",
			Payload: "payload",
		},
		Status:  "Queued",
	},
}

var agents = models.Agents{
	models.Agent{
		Id:   "miau",
		Name: "Miau",
		Facts: []models.Fact{
			models.Fact{
				Id:    "os",
				Name:  "Platform",
				Value: "windows"},
			models.Fact{
				Id:    "arch",
				Name:  "Architecture",
				Value: "amd64"},
		},
	},
	models.Agent{
		Id:   "bup",
		Name: "Bup",
		Facts: []models.Fact{
			models.Fact{
				Id:    "os",
				Name:  "Platform",
				Value: "darwin"},
			models.Fact{
				Id:    "arch",
				Name:  "Architecture",
				Value: "amd64"},
		},
	},
}
