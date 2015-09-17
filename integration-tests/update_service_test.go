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

var arcDeployVersionFlag = flag.String("arc-last-deployed-version", "2015.01.01", "integration-test")

type factArcVersion struct {
	Version string `json:"arc_version"`
}

func TestAgentsAreUpdated(t *testing.T) {
	// override flags if enviroment variable exists
	if os.Getenv("ARC_LAST_DEPLOY_VERSION") != "" {
		deployVersion := os.Getenv("ARC_LAST_DEPLOY_VERSION")
		arcDeployVersionFlag = &deployVersion
	}

	// get the logged agents
	client := NewTestClient()
	statusCode, body := client.Get("/agents", ApiServer)
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
		statusCode, body = client.Get(fmt.Sprint("/agents/", agents[i].AgentID, "/facts"), ApiServer)
		if statusCode != "200 OK" {
			t.Error(fmt.Sprint("Expected to get 200 response code getting facts for agent ", agents[i]))
			continue
		}

		var factVersion factArcVersion
		err = json.Unmarshal(*body, &factVersion)
		if err != nil {
			t.Error(fmt.Sprint("Expected not to get an error unmarshaling for agent ", agents[i]))
			continue
		}

		if !strings.Contains(factVersion.Version, *arcDeployVersionFlag) {
			t.Error(fmt.Sprint("Expected to match versions for agent ", agents[i].AgentID, ". Got ", *arcDeployVersionFlag, " and ", factVersion.Version))
		}
	}
}
