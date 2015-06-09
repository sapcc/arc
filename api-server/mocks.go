package main

import (
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

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
