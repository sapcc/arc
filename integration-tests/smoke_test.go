// +build integration

package integrationTests

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

var updateSiteVersion = flag.String("update-site-version", "2015.01.01", "Expected version of Update Site")
var apiVersion = flag.String("api-version", "2015.01.01", "Expected version of Api Server")

func TestApiServerIsUp(t *testing.T) {
	// override flags if enviroment variable exists
	if e := os.Getenv("API_VERSION"); e != "" {
		apiVersion = &e
	}
	if e := os.Getenv("UPDATE_SITE_VERSION"); e != "" {
		updateSiteVersion = &e
	}

	client, err := NewTestClient()
	if err != nil {
		t.Fatal(err)
	}
	statusCode, body := client.Get("/healthcheck", ApiServer)

	expected := "200 OK"
	if statusCode != expected {
		t.Errorf("Expected to get %#v code for the ApiServer. Got %#v", expected, statusCode)
	}

	bodystring := bytes.NewBuffer(*body).String()
	if !strings.Contains(bodystring, *apiVersion) {
		t.Errorf("ApiServer is running version %#v, expected %#v", bodystring, *apiVersion)
	}

}

func TestUpdateServerIsUp(t *testing.T) {
	client, err := NewTestClient()
	if err != nil {
		t.Fatal(err)
	}
	statusCode, body := client.Get("/healthcheck", UpdateServer)

	expected := "200 OK"
	if statusCode != expected {
		t.Errorf("Expected to get %#v code for the UpdateServer. Got %#v", expected, statusCode)
	}

	bodystring := bytes.NewBuffer(*body).String()
	if !strings.Contains(bodystring, *updateSiteVersion) {
		t.Errorf("UpdateServer is running version %#v, expected %#v", bodystring, *updateSiteVersion)
	}
}
