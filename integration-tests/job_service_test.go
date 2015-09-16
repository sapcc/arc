// +build integration

package integrationTests

import (
	"testing"
	"os"
	"fmt"
	"flag"
	"encoding/json"
	"time"
	
	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

var serverIdentityFlag = flag.String("arc-server-identity", "darwin", "integration-test")

type factArcVersion struct {
	Version string `json:"arc_version"`
}

func TestRunJob(t *testing.T) {	
	// override flags if enviroment variable exists
	if os.Getenv("ARC_SERVER_IDENTITY") != "" {
		serverIdentity := os.Getenv("ARC_SERVER_IDENTITY")
		serverIdentityFlag = &serverIdentity
	}

	client := NewTestClient()	
	to := *serverIdentityFlag
	timeout := 60
	agent := "execute"
	action := "script"
	payload := `echo Script start; for i in {1..2}; do echo \$i; sleep 1s; done; echo Script done`
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
		
	job, err := checkStatus(client, jobId)
	if err != nil {
		t.Error(fmt.Sprint("Expected not to get an error. Got ", err.Error()))
	}
	if job.Status != arc.Queued {
		t.Error("Expected to get job in queued mode")
	}
	
	time.Sleep(time.Second * 1)
	
	job, err = checkStatus(client, jobId)
	if err != nil {
		t.Error(fmt.Sprint("Expected not to get an error. Got ", err.Error()))
	}
	if job.Status != arc.Executing {
		t.Error("Expected to get job in execution mode")
	}

	time.Sleep(time.Second * 2)
	
	job, err = checkStatus(client, jobId)
	if err != nil {
		t.Error(fmt.Sprint("Expected not to get an error. Got ", err.Error()))
	}
	if job.Status != arc.Complete {
		t.Error("Expected to get job in complete mode")
	}
}

// private

func checkStatus(client *Client, jobId models.JobID) (*models.Job, error){
	var job models.Job
	statusCode, body := client.Get(fmt.Sprint("/jobs/", jobId.RequestID), ApiServer)
	if statusCode != "200 OK" {
		return nil, fmt.Errorf("Expected to get 200 response code")
	}
	err := json.Unmarshal(*body, &job)
	if err != nil {
		return nil, fmt.Errorf("Expected not to get an error unmarshaling")
	}		
	return &job, nil
}
