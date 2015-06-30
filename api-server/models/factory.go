package models

import (
	"code.google.com/p/go-uuid/uuid"
	"gitHub.***REMOVED***/monsoon/arc/arc"

	"database/sql"
	"time"
)

func (jobs *Jobs) CreateAndSaveRpcVersionExamples(db *sql.DB, number int) {
	for i := 0; i < number; i++ {
		job := Job{}
		job.RpcVersionExample()
		job.Save(db)
		*jobs = append(*jobs, job)
	}
}

func (job *Job) RpcVersionExample() {
	job.Version = 1
	job.Agent = "rpc"
	job.Action = "version"
	job.To = "darwin"
	job.Timeout = 60
	job.RequestID = uuid.New()
	job.Status = arc.Queued
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
}

func (job *Job) ExecuteScriptExample() {
	job.Version = 1
	job.Agent = "execute"
	job.Action = "script"
	job.To = "darwin"
	job.Timeout = 60
	job.Payload = "echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""
	job.RequestID = uuid.New()
	job.Status = arc.Queued
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
}
