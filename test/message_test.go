package test

import (
	"fmt"
	"gitHub.***REMOVED***/monsoon/onos/onos"
	"testing"
	"encoding/json"
	"errors"
	"strings"
)

var validJson = `{
	"version":1,
	"sender":"some_sender",
	"request_id":"133B0939-76F4-4C9B-99AB-7A6A873E8C9E",
	"to":"somewhere",
	"timeout":1,
	"agent":"rpc",
	"action": "ping",
	"payload": "ping"
}`

var typeValidation = []string{
    `{"version":"version",}`,
	`{"sender":987654321,}`,
    `{"request_id":123,}`,
	`{"to":123,}`,
	`{"timeout":"timeout",}`,
	`{"agent":12345,}`,
	`{"agent":12345,}`,
	`{"payload":12345,}`,
}

func TestParseRequestValidJson(t *testing.T) {	
	data := []byte(validJson)
	request, err := parseRequest(&data)
	if request == nil {
		t.Error("Expected request not nil, got ", request)
	}
	if err != nil {
		t.Error("Expected get one error, got ", err)
	}	
}

func TestParseRequestMalformedJson(t *testing.T) {	
	data := []byte("some text instead of JSON")
	request, err := parseRequest(&data)
	if request != nil {
		t.Error("Expected request nil, got ", request)
	}
	if err == nil {
		t.Error("Expected get one error, got ", err)
	}
}

func TestParseRequestTypeValidation(t *testing.T) {		
	for _, str := range typeValidation {
		data := []byte(str)
		request, err := parseRequest(&data)
		if request != nil {
			t.Error("Expected request nil, got ", request)
		}
		if err == nil {
			t.Error("Expected get one error, got ", err)
		}
	}
}

func TestParseRequestAttrValidation(t *testing.T) {
	str := validJson
	// remove {}
	str = strings.Replace(str, "{", "", -1)
	str = strings.Replace(str, "}", "", -1)
	
	// create arrays removing each time a different attribute
	res := strings.Split(str, ",")
    for i := range res {
		n := make([]string, len(res[:i])+len(res[i+1:])) 		
		copy(n[:], res[:i])
		copy(n[len(res[:i]):], res[i+1:])
		
		// build a json from the string arrays
		jsonString := fmt.Sprint("{", strings.Join(n, ","), "}")
		data := []byte(jsonString)		
		
		//test errors
		request, err := parseRequest(&data)
		if request != nil {
			t.Error("Expected request nil, got ", request)
		}
		if err == nil {
			t.Error("Expected get one error, got ", err)
		}
     }
}

// private

func parseRequest(data *[]byte) (request *onos.Request, err error) {
	// unmarshal    
	err = json.Unmarshal(*data, &request)
	if err != nil {		
		request = nil
		return
	}
	
	// check all attr are given	
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
