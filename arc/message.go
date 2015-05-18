package arc

import (
	"encoding/json"
	"errors"
	"fmt"	
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

func ParseRequest(data *[]byte) (request *Request, err error) {
	// unmarshal    
	err = json.Unmarshal(*data, &request)
	if err != nil {		
		request = nil
		return
	}
	
	// validation
	if request.Version == 0 {
		err = errors.New(fmt.Sprint("Attribute 'Version' is missing or has no valid value. Got ", request.Version))
		request = nil
		return
	}
	
	if len(request.Sender) == 0 {
		err = errors.New(fmt.Sprint("Attribute 'Sender' is missing or empty. Got ", request.Sender))
		request = nil
		return
	}
	
	if len(request.To) == 0 {
		err = errors.New(fmt.Sprint("Attribute 'To' is missing or empty. Got ", request.To))
		request = nil
		return
	}
	
	if len(request.RequestID) == 0 {
		err = errors.New(fmt.Sprint("Attribute 'RequestID' is missing or empty. Got ", request.RequestID))
		request = nil
		return
	}
	
	if request.Timeout == 0 {
		err = errors.New(fmt.Sprint("Attribute 'Timeout' is missing or has no valid value. Got ", request.Timeout))
		request = nil
		return
	}
	
	if len(request.Agent) == 0 {
		err = errors.New(fmt.Sprint("Attribute 'Agent' is missing or empty. Got ", request.Agent))
		request = nil
		return
	}
	
	if len(request.Action) == 0 {
		err = errors.New(fmt.Sprint("Attribute 'Action' is missing or empty. Got ", request.Action))
		request = nil
		return
	}
	
	if len(request.Payload) == 0 {
		err = errors.New(fmt.Sprint("Attribute 'Payload' is missing or empty. Got ", request.Payload))
		request = nil
		return
	}
	
	return
}