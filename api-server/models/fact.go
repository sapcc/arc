package models

import ()

type Fact struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Facts []Fact
