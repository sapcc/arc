// +build integration

package integration_tests

import (
	"testing"
)

func TestApiServerIsUp(t *testing.T) {
	client := NewTestClient()
	statusCode, _ := client.Get("/", ApiServer)

	if statusCode != "200 OK" {
		t.Error("Expected to get 200 response code for the ApiServer")
	}
}

func TestUpdateServerIsUp(t *testing.T) {
	client := NewTestClient()
	statusCode, _ := client.Get("/", UpdateServer)

	if statusCode != "200 OK" {
		t.Error("Expected to get 200 response code for the UpdateServer")
	}
}
