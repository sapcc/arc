// +build integration

package integrationTests

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

var agentIdentityFlag = flag.String("agent-identity", "", "integration-test")

type systemFact struct {
	Os             string `json:"os"`
	PlatformFamily string `json:"platform_family"`
}

func TestRunJob(t *testing.T) {
	// override flags if enviroment variable exists
	if os.Getenv("AGENT_IDENTITY") != "" {
		agentIdentity := os.Getenv("AGENT_IDENTITY")
		agentIdentityFlag = &agentIdentity
	}

	client, err := NewTestClient()
	if err != nil {
		t.Fatal(err)
	}

	// get info about the agent
	statusCode, body := client.GetApiV1(fmt.Sprint("/agents/", *agentIdentityFlag, "/facts"), ApiServer)
	if statusCode != "200 OK" {
		t.Errorf("Expected to get 200 response code getting facts for agent %s got %s ", *agentIdentityFlag, statusCode)
		return
	}

	// transform the body to system facts struct
	var sysFact systemFact
	if err := json.Unmarshal(*body, &sysFact); err != nil {
		t.Error("Expected not to get an error unmarshaling: ", err)
		return
	}

	// build the post request
	payload := `echo Start; for i in {1..2}; do echo $i; sleep 1s; done; echo Done`
	if sysFact.Os == "windows" || sysFact.PlatformFamily == "windows" {
		payload = `echo Start; for($i=1;$i -le 2;$i++){echo $i; sleep -seconds 1}; echo Done`
	}
	to := *agentIdentityFlag
	timeout := 60
	agent := "execute"
	action := "script"
	data := fmt.Sprintf(`{"to":%q,"timeout":%v,"agent":%q,"action":%q,"payload":%q}`, to, timeout, agent, action, payload)
	jsonStr := []byte(data)

	// post the job
	statusCode, body = client.PostApiV1("/jobs", ApiServer, nil, jsonStr)
	if statusCode != "200 OK" {
		t.Errorf("Expected to get 200 response code posting the job, got %s", statusCode)
		return
	}

	var jobId models.JobID
	err = json.Unmarshal(*body, &jobId)
	if err != nil {
		t.Error("Error unmarshalling jobs reply: ", err)
		return
	}

	err = checkStatus(client, jobId, arc.Executing, 3000)
	if err != nil {
		t.Error(err)
		return
	}

	err = checkStatus(client, jobId, arc.Complete, 8000)
	if err != nil {
		t.Error(err)
		return
	}

	// check log
	statusCode, body = client.GetApiV1(fmt.Sprint("/jobs/", jobId.RequestID, "/log"), ApiServer)
	if statusCode != "200 OK" {
		t.Errorf("Expected to get 200 response code for job %s, got %s", jobId.RequestID, statusCode)
	}
	if len(string(*body)) == 0 {
		t.Error("Expected to get a log")
	}
}

// private

func checkStatus(client *Client, jobId models.JobID, status arc.JobState, timeout int) error {
	var job *models.Job
	var err error
	for {
		job, err = getJobStatus(client, jobId)
		if err != nil {
			break
		}
		if job.Status == status {
			break
		}
		if timeout < 0 {
			err = fmt.Errorf(fmt.Sprint("Timeout: Expected to get status ", status, " for job ", jobId.RequestID, ". Got status ", job.Status))
			break
		}

		timeout = timeout - 100
		time.Sleep(time.Millisecond * 100)
	}

	return err
}

func getJobStatus(client *Client, jobId models.JobID) (*models.Job, error) {
	var job models.Job
	statusCode, body := client.GetApiV1(fmt.Sprint("/jobs/", jobId.RequestID), ApiServer)
	if statusCode != "200 OK" {
		return nil, fmt.Errorf("Expected to get 200 response code getting job %s", jobId.RequestID)
	}
	err := json.Unmarshal(*body, &job)
	if err != nil {
		return nil, fmt.Errorf("Expected not to get an error unmarshaling body from job %s", jobId.RequestID)
	}
	return &job, nil
}
