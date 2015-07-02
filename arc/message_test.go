// +build !integration

package arc

import (
	"fmt"
	"strings"
	"testing"
)

var validJson = `{
	"version":1,
	"sender":"some_sender",
	"request_id":"133B0939-76F4-4C9B-99AB-7A6A873E8C9E",
	"to":"somewhere",
	"timeout":1,
	"agent":"rpc",
	"action": "ping"
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
	request, err := ParseRequest(&data)
	if request == nil {
		t.Error("Expected request not nil, got ", request)
	}
	if err != nil {
		t.Error("Expected get one error, got ", err)
	}
}

func TestParseRequestMalformedJson(t *testing.T) {
	data := []byte("some text instead of JSON")
	request, err := ParseRequest(&data)
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
		request, err := ParseRequest(&data)
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
		request, err := ParseRequest(&data)
		if request != nil {
			t.Error("Expected request nil, got ", request)
		}
		if err == nil {
			t.Error("Expected get one error, got ", err)
		}
	}
}
