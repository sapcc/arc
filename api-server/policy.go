package main

import "github.com/ory/ladon"

var policies = []ladon.Policy{
	&ladon.DefaultPolicy{
		ID:          "1",
		Description: `admins can do anything in any resource`,
		Subjects:    []string{"automation_admin"},
		Actions:     []string{"<.*>"},
		Resources:   []string{"<.*>"},
		Effect:      ladon.AllowAccess,
	},
	&ladon.DefaultPolicy{
		ID:          "2",
		Description: "viewers can call get in any resource",
		Subjects:    []string{"automation_viewer"},
		Actions:     []string{"get"},
		Resources:   []string{"<.*>"},
		Effect:      ladon.AllowAccess,
	},
}
