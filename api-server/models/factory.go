package models

import (
	"code.google.com/p/go-uuid/uuid"
	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"		

	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Request struct {
	arc.Request
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
	job.CustomExecuteScriptExample(arc.Queued, time.Now().Add(-1 * time.Minute), 60)
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

func (request *Request) RegistryExample() {
	request.Sender = uuid.New()
	request.Version = 1
	request.Agent = "registration"
	request.Action = "register"
	request.To = "registry"
	request.Timeout = 5
	request.Payload = `{"os": "darwin", "online": true, "project": "test-project", "hostname": "BERM32186999A", "identity": "darwin", "platform": "mac_os_x", "arc_version": "0.1.0-dev(69f43fd)", "memory_used": 9206046720, "memory_total": 17179869184, "organization": "test-org"}`
	request.RequestID = uuid.New()
}

func (agent *Agent) Example() {
	agent.AgentID = uuid.New()
	agent.Project = "test project"
	agent.Organization = "test organization"
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
		var lastInsertId string
		err := db.QueryRow(ownDb.InsertFactQuery, agent.AgentID, agent.Project, agent.Organization, "{}", agent.CreatedAt, agent.UpdatedAt).Scan(&lastInsertId);				
		if err != nil {
			log.Error(err)
		}
		*agents = append(*agents, agent)
	}
}

func (agents *Agents) CreateAndSaveRegistryExamples(db *sql.DB, number int) {
	for i := 0; i < number; i++ {
		// build a request
		req := Request{}
		req.RegistryExample()
		// save a job
		fact := Fact{}
		err := fact.ProcessRequest(db, &req.Request)
		if err != nil {
			log.Error(err)
		}
		agent := fact.Agent
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
