package onos

import (
	"encoding/json"
)

type Request struct {
	Version   int    `json:"version"`
	Sender    string `json:"sender"`
	RequestID string `json:"request_id"`
	To        string `json:"to"`
	Timeout   uint64 `json:"timeout"`
	Agent     string `json:"agent"`
	Action    string `json:"action"`
	Payload   string `json:"payload"`
}

type Reply struct {
	Version   int    `json:"version"`
	Sender    string `json:"sender"`
	RequestID string `json:"request_id"`
	Agent     string `json:"agent"`
	Action    string `json:"action"`
	State     string `json:"state"`
	Final     bool   `json:"final"`
	Payload   string `json:"payload"`
}

func (r *Request) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		*Request
		Type string `json:"type"`
	}{r, "request"})
}

func (r *Reply) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		*Reply
		Type string `json:"type"`
	}{r, "reply"})
}

func CreateReply(request *Request, payload string) *Reply {

	return &Reply{
		Version:   1,
		Agent:     request.Agent,
		Action:    request.Action,
		Payload:   payload,
		RequestID: request.RequestID,
	}

}
