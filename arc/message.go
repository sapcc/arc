package arc

import (
	"encoding/json"
	"fmt"

	"github.com/pborman/uuid"
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
	Number    uint     `json:"number"`
}

type Registration struct {
	RegistrationID string `json:"registration_id"`
	Version        int    `json:"version"`
	Sender         string `json:"sender"`
	Organization   string `json:"organization"`
	Project        string `json:"project"`
	Payload        string `json:"payload"`
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

func (r *Registration) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		*Registration
		Type string `json:"type"`
	}{r, "registration"})
}

func CreateReply(request *Request, identity string, state JobState, payload string, number uint) *Reply {

	final := state == Complete || state == Failed
	return &Reply{
		Version:   1,
		Agent:     request.Agent,
		Action:    request.Action,
		Payload:   payload,
		RequestID: request.RequestID,
		State:     state,
		Final:     final,
		Number:    number,
		Sender:    identity,
	}

}

func CreateRequest(agent string, action string, identity string, to string, timeout int, payload string) (*Request, error) {
	request := Request{
		Version:   1,
		Agent:     agent,
		Action:    action,
		To:        to,
		Timeout:   timeout,
		Payload:   payload,
		RequestID: uuid.New(),
		Sender:    identity,
	}

	err := ValidateRequest(&request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func CreateRegistration(organization, project, identity, payload string) (*Registration, error) {
	registration := Registration{
		RegistrationID: uuid.New(),
		Version:        1,
		Project:        project,
		Organization:   organization,
		Sender:         identity,
		Payload:        payload,
	}

	return &registration, nil
}

func ParseRegistration(data *[]byte) (*Registration, error) {
	var reg Registration
	err := json.Unmarshal(*data, &reg)
	if err != nil {
		return nil, err
	}

	err = ValidateRegistration(&reg)
	if err != nil {
		return nil, err
	}

	return &reg, err
}

func ParseRequest(data *[]byte) (*Request, error) {
	var request Request
	// unmarshal
	err := json.Unmarshal(*data, &request)
	if err != nil {
		return nil, err
	}

	// validation
	err = ValidateRequest(&request)
	if err != nil {
		return nil, err
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

	// validation
	err = ValidateReply(&reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func ValidateRequest(request *Request) error {
	field_error := "Attribute '%s' is missing or invalid"

	if request.Version < 1 {
		return fmt.Errorf(field_error, "Version")
	}

	if request.Sender == "" {
		return fmt.Errorf(field_error, "Sender")
	}

	if request.To == "" {
		return fmt.Errorf(field_error, "To")
	}

	if request.RequestID == "" {
		return fmt.Errorf(field_error, "RequestID")
	}

	if request.Timeout < 1 {
		return fmt.Errorf(field_error, "Timeout")
	}

	if request.Agent == "" {
		return fmt.Errorf(field_error, "Agent")
	}

	if request.Action == "" {
		return fmt.Errorf(field_error, "Action")
	}

	return nil
}

func ValidateReply(reply *Reply) error {
	field_error := "Attribute '%s' is missing or invalid"

	if reply.Version < 1 {
		return fmt.Errorf(field_error, "Version")
	}

	if reply.RequestID == "" {
		return fmt.Errorf(field_error, "RequestID")
	}

	if reply.Agent == "" {
		return fmt.Errorf(field_error, "Agent")
	}

	if reply.Action == "" {
		return fmt.Errorf(field_error, "Action")
	}

	if reply.State == 0 {
		return fmt.Errorf(field_error, "State")
	}

	if reply.Number < 0 {
		return fmt.Errorf(field_error, "Number")
	}

	return nil
}

func ValidateRegistration(reg *Registration) error {
	field_error := "Attribute '%s' is missing or invalid"

	if reg.Version < 1 {
		return fmt.Errorf(field_error, "Version")
	}

	if reg.Organization == "" {
		return fmt.Errorf(field_error, "Organization")
	}
	if reg.Project == "" {
		return fmt.Errorf(field_error, "Project")
	}
	if reg.Sender == "" {
		return fmt.Errorf(field_error, "Sender")
	}

	return nil

}
