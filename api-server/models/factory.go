package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"

	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type Request struct {
	arc.Request
}

type Registration struct {
	arc.Registration
}

type Reply struct {
	arc.Reply
}

func (jobs *Jobs) CreateAndSaveRpcVersionExamples(db *sql.DB, number int) {
	now := time.Now()
	for i := 0; i < number; i++ {
		job := Job{}
		job.RpcVersionExample()
		job.UpdatedAt = now.Add(time.Duration(i) * time.Minute)
		err := job.Save(db)
		if err != nil {
			log.Error(err)
		}
		*jobs = append(*jobs, job)
	}
}

func (job *Job) RpcVersionExample() {
	job.Sender = "windows"
	job.Version = 1
	job.Agent = "rpc"
	job.Action = "version"
	job.To = "darwin"
	job.Timeout = 60
	job.RequestID = uuid.New()
	job.Status = arc.Queued
	job.CreatedAt = time.Now().Add(-1 * time.Minute)
	job.UpdatedAt = time.Now().Add(-1 * time.Minute)
}

func (job *Job) CustomExecuteScriptExample(status arc.JobState, createdAt time.Time, timeout int) {
	job.Sender = "windows"
	job.Version = 1
	job.Agent = "execute"
	job.Action = "script"
	job.To = "darwin"
	job.Timeout = timeout
	job.Payload = "echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""
	job.RequestID = uuid.New()
	job.Status = status
	job.CreatedAt = createdAt
	job.UpdatedAt = createdAt
}

func (job *Job) ExecuteScriptExample() {
	job.CustomExecuteScriptExample(arc.Queued, time.Now().Add(-1*time.Minute), 60)
}

func (reply *Reply) ExecuteScriptExample(id string, final bool, payload string, number uint) {
	reply.Version = 1
	reply.Sender = "darwin"
	reply.RequestID = id
	reply.Agent = "execute"
	reply.Action = "script"
	if final == true {
		reply.State = arc.Complete
	} else {
		reply.State = arc.Executing
	}
	reply.Final = final
	reply.Payload = payload
	reply.Number = number
}

func (req *Request) Example() {
	req.Version = 1
	req.Sender = "windows"
	req.RequestID = uuid.New()
	req.To = "darwin"
	req.Timeout = 60
	req.Agent = "execute"
	req.Action = "script"
	req.Payload = "echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""
}

func (reg *Registration) Example() {
	reg.Sender = uuid.New()
	reg.Version = 1
	reg.Project = "test-proj"
	reg.Organization = "test-org"
	reg.Payload = `{"os": "darwin", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
}

func (reg *Registry) Example() {
	reg.RegistryID = uuid.New()
	reg.AgentID = "darwin"
}

func (agent *Agent) Example() {
	agent.AgentID = uuid.New()
	agent.Project = "test-project"
	agent.Organization = "test-org"
	agent.Facts = `{"os": "darwin", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()
}

func (agents *Agents) CreateAndSaveAgentExamples(db *sql.DB, number int) {
	now := time.Now()
	for i := 0; i < number; i++ {
		agent := Agent{}
		agent.Example()
		agent.CreatedAt = now.Add(time.Duration(i) * time.Minute)
		agent.UpdatedAt = now.Add(time.Duration(i) * time.Minute)
		err := agent.Save(db)
		if err != nil {
			log.Error(err)
		}
		*agents = append(*agents, agent)
	}
}

func (logpart *LogPart) SaveLogPartExamples(db *sql.DB, id string) string {
	var contentSlice [3]string
	reply := Reply{}
	for i := 0; i < 3; i++ {
		chunk := fmt.Sprint("Log chunk ", i)
		reply.ExecuteScriptExample(id, false, chunk, uint(i))
		err := ProcessLogReply(db, &reply.Reply)
		if err != nil {
			log.Error(err)
		}
		contentSlice[i] = chunk
	}
	return strings.Join(contentSlice[:], "")
}
