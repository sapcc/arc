// +build integration

package integration_tests

import (
	"testing"
	"os"
	"fmt"
	"encoding/json"
	"strings"	
)

type factArcVersion struct {
	Version string `json:"arc_version"`
}

func TestAgentsAreUpdated(t *testing.T) {
	deployedVersion := os.Getenv("ARC_DEPLOY_VERSION")
	agent1 := os.Getenv("ARC_AGENT_1")
	
	client := NewTestClient()
	statusCode, body := client.Get(fmt.Sprint("/agents/", agent1, "/facts"), ApiServer)

	if statusCode != "200 OK" {
		t.Error(fmt.Sprint("Expected to get 200 response code for agent ", agent1))
	}
	
	var factVersion factArcVersion
	err := json.Unmarshal(*body, &factVersion)
	if err != nil {
		t.Error("Expected not to get an error")
	}
	
	if !strings.Contains(factVersion.Version, deployedVersion) { 
		t.Error(fmt.Sprint("Expected to match versions. Got ", deployedVersion, " and ", factVersion.Version))
	}
}