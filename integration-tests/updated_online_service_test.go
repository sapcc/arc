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
	client, err := NewTestClient()
	if err != nil {
		t.Fatal(err)
	}

	statusCode, body := client.GetApiV1("/agents", ApiServer)
	if statusCode != "200 OK" {
		t.Errorf("Expected to get 200 response code getting all agents, got %s", statusCode)
		return
	}

	var agents models.Agents
	err = json.Unmarshal(*body, &agents)
	if err != nil {
		t.Error("Expected not to get an error unmarshaling agents")
		return
	}

	// check the version from each agent

	results := make(map[int]bool, len(agents))
	for i := 0; i < *timeout; i++ {

		for idx, agent := range agents {
			if _, ok := results[idx]; ok {
				//skip finished agents
				continue
			}
			statusCode, body = client.GetApiV1(fmt.Sprint("/agents/", agent.AgentID, "/facts"), ApiServer)
			if statusCode != "200 OK" {
				t.Error("Expected to get 200 response code getting facts for agent ", agent)
				results[idx] = false
				continue
			}

			var facts checkFacts
			err = json.Unmarshal(*body, &facts)
			if err != nil {
				t.Error("Error unmarshaling response for ", agent)
				results[idx] = false
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
			results[idx] = true
		}
		if len(results) == len(agents) {
			break
		}
		fmt.Println("Sleeping for 1 second...")
		time.Sleep(1 * time.Second)
	}
	for i, agent := range agents {
		if _, ok := results[i]; !ok {
			t.Errorf("Agent %s failed to update before %ds timeout was reached", agent.AgentID, *timeout)
		}
	}
}
