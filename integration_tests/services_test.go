// +build integration

package integration_tests

import (
	"testing"
	"os"
	"fmt"
	"encoding/json"
	"strings"
	"time"
	
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type factArcVersion struct {
	Version string `json:"arc_version"`
}

func TestAgentsAreUpdated(t *testing.T) {
	deployedVersion := os.Getenv("ARC_DEPLOY_VERSION")
	agentIdentity1 := os.Getenv("ARC_AGENT_IDENTITY_1")
	
	client := NewTestClient()
	statusCode, body := client.Get(fmt.Sprint("/agents/", agentIdentity1, "/facts"), ApiServer)

	if statusCode != "200 OK" {
		t.Error(fmt.Sprint("Expected to get 200 response code for agent ", agentIdentity1))
	}
	
	var factVersion factArcVersion
	err := json.Unmarshal(*body, &factVersion)
	if err != nil {
		t.Error("Expected not to get an error unmarshaling")
	}
	
	if !strings.Contains(factVersion.Version, deployedVersion) { 
		t.Error(fmt.Sprint("Expected to match versions. Got ", deployedVersion, " and ", factVersion.Version))
	}
}

func TestRunJob(t *testing.T) {	
	client := NewTestClient()
	
	to := os.Getenv("ARC_AGENT_IDENTITY_1")
	timeout := 60
	agent := "execute"
	action := "script"
	payload := `echo Script start; for i in {1..5}; do echo \$i; sleep 1s; done; echo Script done`
	data := fmt.Sprintf(`{"to":%q,"timeout":%v,"agent":%q,"action":%q,"payload":%q}`, to, timeout, agent, action, payload)
	jsonStr := []byte(data)
	
	statusCode, body := client.Post("/jobs", ApiServer, nil, jsonStr)
	if statusCode != "200 OK" {
		t.Error(fmt.Sprint("Expected to get 200 response code"))
	}
	var jobId models.JobID
	err := json.Unmarshal(*body, &jobId)
	if err != nil {
		t.Error("Expected not to get an error unmarshaling")
	}	
	
	statusCode, body = client.Get(fmt.Sprint("/jobs/", jobId.RequestID), ApiServer)
	if statusCode != "200 OK" {
		t.Error(fmt.Sprint("Expected to get 200 response code"))
	}
	
	var job models.Job
	err = json.Unmarshal(*body, &job)
	if err != nil {
		t.Error("Expected not to get an error unmarshaling")
	}	
	if job.Status != arc.Queued {
		t.Error("Expected not to get job in execution mode")
	}
	
	time.Sleep(time.Second * 3)
	
	statusCode, body = client.Get(fmt.Sprint("/jobs/", jobId.RequestID), ApiServer)
	if statusCode != "200 OK" {
		t.Error(fmt.Sprint("Expected to get 200 response code"))
	}
	err = json.Unmarshal(*body, &job)
	if err != nil {
		t.Error("Expected not to get an error unmarshaling")
	}	
	if job.Status != arc.Executing {
		t.Error("Expected not to get job in execution mode")
	}

	time.Sleep(time.Second * 3)

	statusCode, body = client.Get(fmt.Sprint("/jobs/", jobId.RequestID), ApiServer)
	if statusCode != "200 OK" {
		t.Error(fmt.Sprint("Expected to get 200 response code"))
	}
	err = json.Unmarshal(*body, &job)
	if err != nil {
		t.Error("Expected not to get an error unmarshaling")
	}	
	if job.Status != arc.Complete {
		t.Error("Expected not to get job in complete mode")
	}
}

