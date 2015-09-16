// +build integration

package integrationTests

import (
	"testing"
	"fmt"
	"flag"
	"encoding/json"
)

var serverIdentityFlag = flag.String("arc-server-identity", "darwin", "integration-test")

type Facts struct {
	Version            string `json:"arc_version"`
	DefaultGateway     string `json:"default_gateway"`
	DefaultInterface   string `json:"default_interface"`
	Domain					   string `json:"domain"`
	FQDN						   string `json:"fqdn"`
	Hostname				   string `json:"hostname"`
	IpAddress				   string `json:"ipaddress"`
	Platform					 string `json:"platform"`
	PlatformFamily		 string `json:"platform_family"`
	PlatformVersion    string `json:"platform_version"`
}

func TestRunFacts(t *testing.T) {
	client := NewTestClient()	
	statusCode, body := client.Get(fmt.Sprint("/agents/", *serverIdentityFlag, "/facts"), ApiServer)
	if statusCode != "200 OK" {
		t.Error(fmt.Sprint("Expected to get 200 response code for agent ", *serverIdentityFlag))
	}
	
	var facts Facts
	err := json.Unmarshal(*body, &facts)
	if err != nil {
		t.Error("Expected not to get an error unmarshaling")
	}

	if facts.Version == "" {
		t.Error(fmt.Sprintf("Expected version to not be empty. Got %q", facts.Version))
	}	
	if facts.DefaultGateway == "" {
		t.Error(fmt.Sprintf("Expected default gateway to not be empty. Got %q", facts.DefaultGateway))
	}	
	if facts.DefaultInterface == "" {
		t.Error(fmt.Sprintf("Expected default interface to not be empty. Got %q", facts.DefaultInterface))
	}
	if facts.Domain == "" {
		t.Error(fmt.Sprintf("Expected domain to not be empty. Got %q", facts.Domain))
	}
	if facts.FQDN == "" {
		t.Error(fmt.Sprintf("Expected fqdn to not be empty. Got %q", facts.FQDN))
	}
	if facts.Hostname == "" {
		t.Error(fmt.Sprintf("Expected hostname to not be empty. Got %q", facts.Hostname))
	}
	if facts.IpAddress == "" {
		t.Error(fmt.Sprintf("Expected ip address to not be empty. Got %q", facts.IpAddress))
	}
	if facts.Platform == "" {
		t.Error(fmt.Sprintf("Expected platform to not be empty. Got %q", facts.Platform))
	}
	if facts.PlatformFamily == "" {
		t.Error(fmt.Sprintf("Expected platform family to not be empty. Got %q", facts.PlatformFamily))
	}
	if facts.PlatformVersion == "" {
		t.Error(fmt.Sprintf("Expected platform version to not be empty. Got %q", facts.PlatformVersion))
	}	
}