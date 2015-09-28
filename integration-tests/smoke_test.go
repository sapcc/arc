// +build integration

package integrationTests

import (
	"testing"
)

func TestApiServerIsUp(t *testing.T) {
	client := NewTestClient()
	statusCode, _ := client.Get("/", ApiServer)

	expected := "200 OK"
	if statusCode != expected {
		t.Errorf("Expected to get %#v code for the ApiServer. Got %#v", expected, statusCode)
	}
}

func TestUpdateServerIsUp(t *testing.T) {
	client := NewTestClient()
	statusCode, _ := client.Get("/", UpdateServer)

	expected := "200 OK"
	if statusCode != expected {
		t.Errorf("Expected to get %#v code for the UpdateServer. Got %#v", expected, statusCode)
	}
}
