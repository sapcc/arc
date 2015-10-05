// +build integration

package integrationTests

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
)

var arcLatestVersion = flag.String("latest-version", "2015.01.01", "integration-test")
var timeout = flag.Int("timeout", 20, "timeout waiting for agents to update")

type checkFacts struct {
	Version string `json:"arc_version"`
	Online  bool   `json:"online"`
}

func TestAgentsAreUpdatedAndOnline(t *testing.T) {
	// override flags if enviroment variable exists
	if e := os.Getenv("LATEST_VERSION"); e != "" {
		arcLatestVersion = &e
	}
	if os.Getenv("TIMEOUT") != "" {
		i, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
		timeout = &i
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

	for i := 0; i < *timeout; i++ {

		for i, agent := range agents {
			statusCode, body = client.GetApiV1(fmt.Sprint("/agents/", agent.AgentID, "/facts"), ApiServer)
			if statusCode != "200 OK" {
				t.Error("Expected to get 200 response code getting facts for agent ", agent)
				agents = append(agents[:i], agents[:i+1]...)
				continue
			}

			var facts checkFacts
			err = json.Unmarshal(*body, &facts)
			if err != nil {
				t.Error("Error unmarshaling response for ", agent)
				agents = append(agents[:i], agents[:i+1]...)
				continue
			}

			// check version
			if !strings.Contains(facts.Version, *arcLatestVersion) {
				fmt.Printf("Agent %s is running version %#v, expected %#v\n", agent.AgentID, facts.Version, *arcLatestVersion)
				continue
			}

			// check online
			if facts.Online == false {
				t.Error(fmt.Sprint("Expected agent ", agent.AgentID, " to be online."))
			}
			fmt.Printf("Agent %s is online and updated\n", agent.AgentID)
			agents = append(agents[:i], agents[:i+1]...)
		}
		fmt.Println("Sleeping for 1 second...")
		time.Sleep(1 * time.Second)
	}
	for _, agent := range agents {
		t.Errorf("Agent %s failed to update before %ds timeout was reached", agent.AgentID, *timeout)
	}
}
