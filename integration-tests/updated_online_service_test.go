// +build integration

package integrationTests

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

var arcDeployVersionFlag = flag.String("latest-version", "2015.01.01", "integration-test")

type checkFacts struct {
	Version string `json:"arc_version"`
	Online  bool   `json:"online"`
}

func TestAgentsAreUpdatedAndOnline(t *testing.T) {
	// override flags if enviroment variable exists
	if os.Getenv("LATEST_VERSION") != "" {
		deployVersion := os.Getenv("LATEST_VERSION")
		arcDeployVersionFlag = &deployVersion
	}

	// get the logged agents
	client := NewTestClient()
	statusCode, body := client.GetApiV1("/agents", ApiServer)
	if statusCode != "200 OK" {
		t.Error("Expected to get 200 response code getting all agents")
		return
	}

	var agents models.Agents
	err := json.Unmarshal(*body, &agents)
	if err != nil {
		t.Error("Expected not to get an error unmarshaling agents")
		return
	}

	// check the version from each agent
	for i := 0; i < len(agents); i++ {
		statusCode, body = client.GetApiV1(fmt.Sprint("/agents/", agents[i].AgentID, "/facts"), ApiServer)
		if statusCode != "200 OK" {
			t.Error(fmt.Sprint("Expected to get 200 response code getting facts for agent ", agents[i]))
			continue
		}

		var facts checkFacts
		err = json.Unmarshal(*body, &facts)
		if err != nil {
			t.Error(fmt.Sprint("Expected not to get an error unmarshaling for agent ", agents[i]))
			continue
		}

		// check version
		if !strings.Contains(facts.Version, *arcDeployVersionFlag) {
			t.Error(fmt.Sprint("Expected to match versions for agent ", agents[i].AgentID, ". Got environment version ", *arcDeployVersionFlag, " and fact version ", facts.Version))
		}

		// check online
		if facts.Online == false {
			t.Error(fmt.Sprint("Expected agent ", agents[i].AgentID, " to be online."))
		}
	}
}
