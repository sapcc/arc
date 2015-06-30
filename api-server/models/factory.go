package models

import (
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"code.google.com/p/go-uuid/uuid"
		
	"time"
)

func RpcVersionJob() Job {
	return Job{
		Request: arc.Request{
			Version:   1,
			Agent:     "rpc",
			Action:    "version",
			To:        "darwin",
			Timeout:   60,
			RequestID: uuid.New(),
			Sender:    "darwin",
		},
		Status:    arc.Queued,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}	
}

func ExecuteSctiptJob() Job {
	return Job{
		Request: arc.Request{
			Version:   1,
			Agent:     "execute",
			Action:    "script",
			To:        "darwin",
			Timeout:   60,
			Payload:   "echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\"",
			RequestID: uuid.New(),
			Sender:    "darwin",
		},
		Status:    arc.Queued,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
