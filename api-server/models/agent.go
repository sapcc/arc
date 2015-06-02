package models

import ()

type Agent struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Facts []Fact `json:"facts"`
}

type Agents []Agent
