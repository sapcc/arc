// +build integration

package integrationTests

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

var arcLatestVersion = flag.String("latest-version", "2015.01.01", "integration-test")

func TestApiServerIsUp(t *testing.T) {
	// override flags if enviroment variable exists
	if e := os.Getenv("LATEST_VERSION"); e != "" {
		arcLatestVersion = &e
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
	if !strings.Contains(bodystring, *arcLatestVersion) {
		t.Errorf("ApiServer is running version %#v, expected %#v", bodystring, *arcLatestVersion)
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
	if !strings.Contains(bodystring, *arcLatestVersion) {
		fmt.Printf("UpdateServer is running version %#v, expected %#v\n", bodystring, *arcLatestVersion)
	}
}
