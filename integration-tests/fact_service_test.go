// +build integration

package integrationTests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type Facts struct {
	Version          string `json:"arc_version"`
	DefaultGateway   string `json:"default_gateway"`
	DefaultInterface string `json:"default_interface"`
	Domain           string `json:"domain"`
	FQDN             string `json:"fqdn"`
	Hostname         string `json:"hostname"`
	IpAddress        string `json:"ipaddress"`
	Platform         string `json:"platform"`
	PlatformFamily   string `json:"platform_family"`
	PlatformVersion  string `json:"platform_version"`
}

func TestRunFacts(t *testing.T) {
	// override flags if enviroment variable exists
	if os.Getenv("AGENT_IDENTITY") != "" {
		agentIdentity := os.Getenv("AGENT_IDENTITY")
		agentIdentityFlag = &agentIdentity
	}

	client, err := NewTestClient()
	if err != nil {
		t.Fatal(err)
	}

	// get the facts for the given agent
	statusCode, body := client.GetApiV1(fmt.Sprint("/agents/", *agentIdentityFlag, "/facts"), ApiServer)
	if statusCode != "200 OK" {
		t.Error(fmt.Sprint("Expected to get 200 response code for agent ", *agentIdentityFlag))
		return
	}

	// transform the body to facts struct
	var facts Facts
	if err := json.Unmarshal(*body, &facts); err != nil {
		t.Error("Expected not to get an error unmarshaling: ", err)
		return
	}

	// check
	if facts.Version == "" {
		t.Errorf("Expected version to not be empty. Got %q", facts.Version)
	}
	if facts.DefaultGateway == "" {
		t.Errorf("Expected default gateway to not be empty. Got %q", facts.DefaultGateway)
	}
	if facts.DefaultInterface == "" {
		t.Errorf("Expected default interface to not be empty. Got %q", facts.DefaultInterface)
	}
	if facts.Domain == "" {
		t.Errorf("Expected domain to not be empty. Got %q", facts.Domain)
	}
	if facts.FQDN == "" {
		t.Errorf("Expected fqdn to not be empty. Got %q", facts.FQDN)
	}
	if facts.Hostname == "" {
		t.Errorf("Expected hostname to not be empty. Got %q", facts.Hostname)
	}
	if facts.IpAddress == "" {
		t.Errorf("Expected ip address to not be empty. Got %q", facts.IpAddress)
	}
	if facts.Platform == "" {
		t.Errorf("Expected platform to not be empty. Got %q", facts.Platform)
	}
	if facts.PlatformFamily == "" {
		t.Errorf("Expected platform family to not be empty. Got %q", facts.PlatformFamily)
	}
	if facts.PlatformVersion == "" {
		t.Errorf("Expected platform version to not be empty. Got %q", facts.PlatformVersion)
	}
}
