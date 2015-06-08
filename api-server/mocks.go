package main

import (
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

var jobs = models.Jobs{
	models.Job{
		ReqID:   "1234567890",
		Payload: "payload",
		Status:  "Queued",
	},
	models.Job{
		ReqID:   "miauBup",
		Payload: "payload",
		Status:  "Executing",
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
