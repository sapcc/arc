package arc

import (
	"encoding/json"
	"fmt"

	"code.google.com/p/go-uuid/uuid"
)

type Request struct {
	Version   int    `json:"version"`
	Sender    string `json:"sender"`
	RequestID string `json:"request_id"`
	To        string `json:"to"`
	Timeout   int    `json:"timeout"`
	Agent     string `json:"agent"`
	Action    string `json:"action"`
	Payload   string `json:"payload"`
}

type Reply struct {
	Version   int      `json:"version"`
	Sender    string   `json:"sender"`
	RequestID string   `json:"request_id"`
	Agent     string   `json:"agent"`
	Action    string   `json:"action"`
	State     JobState `json:"state"`
	Final     bool     `json:"final"`
	Payload   string   `json:"payload"`
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

func CreateReply(request *Request, state JobState, payload string) *Reply {

	final := state == Complete || state == Failed
	return &Reply{
		Version:   1,
		Agent:     request.Agent,
		Action:    request.Action,
		Payload:   payload,
		RequestID: request.RequestID,
		State:     state,
		Final:     final,
	}

}

func CreateRequest(agent string, action string, to string, timeout int, payload string) *Request {
	return &Request{
		Version:   1,
		Agent:     agent,
		Action:    action,
		To:        to,
		Timeout:   timeout,
		Payload:   payload,
		RequestID: uuid.New(),
		Sender:    "me",
	}
}
func ParseRequest(data *[]byte) (*Request, error) {
	var request Request
	// unmarshal
	err := json.Unmarshal(*data, &request)
	if err != nil {
		return nil, err
	}

	field_error := "Attribute '%s' is missing or invalid"

	// validation
	if request.Version < 1 {
		return nil, fmt.Errorf(field_error, "Version")
	}

	if request.Sender == "" {
		return nil, fmt.Errorf(field_error, "Sender")
	}

	if request.To == "" {
		return nil, fmt.Errorf(field_error, "To")
	}

	if request.RequestID == "" {
		return nil, fmt.Errorf(field_error, "RequestID")
	}

	if request.Timeout < 1 {
		return nil, fmt.Errorf(field_error, "Timeout")
	}

	if request.Agent == "" {
		return nil, fmt.Errorf(field_error, "Agent")
	}

	if request.Action == "" {
		return nil, fmt.Errorf(field_error, "Action")
	}

	return &request, nil
}

func ParseReply(data *[]byte) (*Reply, error) {
	var reply Reply
	// unmarshal
	err := json.Unmarshal(*data, &reply)
	if err != nil {
		return nil, err
	}

	field_error := "Attribute '%s' is missing or invalid"

	if reply.Version < 1 {
		return nil, fmt.Errorf(field_error, "Version")
	}

	if reply.RequestID == "" {
		return nil, fmt.Errorf(field_error, "RequestID")
	}

	if reply.Agent == "" {
		return nil, fmt.Errorf(field_error, "Agent")
	}

	if reply.Action == "" {
		return nil, fmt.Errorf(field_error, "Action")
	}

	if reply.State == 0 {
		return nil, fmt.Errorf(field_error, "State")
	}

	return &reply, nil
}
