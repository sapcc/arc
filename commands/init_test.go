package commands

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/codegangsta/cli"
)

var mockMetadata = `{"random_seed": "TZutBbpShKDcvbpi60ud25dAP6ZsfzbAqc9PjnGqS03IB1R/AEtaGxQKOv2RR19qmEaNe8cQHENJpJ+2rIRC0WDfkpTvxLRfWsguxEml17jdHc+JQ/4mqTN4PeLnVmoV1ciuzdT0P+sxfoL7TCJcPEZnLrcICFPTs7D3QYeSrsjYvryKMFCAPX4ZnFppsULgUDUf45k16ji7GFslH1k1QIJUKNLxDnhyZITtk++zOo/5t8pNrZTWgbD9EBERlh1MvggmZVokHyEIYX3RRII5f2o5vL9dxdZ1TYoWBZ7gWKmbb3phJOYyLCCfPzMxJyvP/OYytegI5xqG7iYLrB6MXpgrwWPGZQL7UeweBVnXcZpwPYSFxxVfshIPG3c/uqqv43gQHBcUS+Zs2jwNsscQipDcqFhbn/WbcfNYnr0Zr3PfdgDo8TNMjRDSNNzGhvLhAbMFGuhk2CLC4LqUYtyQLjASJeAffV8QYreFGG523QVYDX3JJR0+r3Oh3qNY4DDylB8Mwpr1+89wbLZVb32OQbYnjZG3LcIEkPpMTIPXvb34OnaV8JgjA9ILrrVeMy8swQENY1HNEF2O6BCt+cpAGPxuZ3ierEBdJdrhJbpN+PfsR9800UTmB+InWE+U5qIP0ke1o7ElHEayRi2rv9qh1uTk7kaqSTOxcLnghJc8FUo=", "uuid": "16d2ef3a-da47-43bb-8204-e4fd6a5db689", "availability_zone": "eu-de-1a", "keys": [{"data": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDeevePzO0cd4sV/iIvKsaVw08s8gC0BuzyJu7/hSt10tt4O46nZpqLLybB/K/telZ2vpgvVSxfpOMWngVCIegv72jVnEMpD1WIL40ackN7TRRfT6JLrSYwvoUYgE3CIk8TBZYf9OalXWWXVgYdHp/12u7NMOEvwBEdor9aWKB39ojnuA3s5guZt4fqBuOoaYE/32W4sL4TL7QqLBBqdGKjOZKvVZITKr4IPn4EDUVoGJ2hKS8f89kNSvmDe4tFgWSu7mohc9V7M8N0TkNz9bXIFe+9tZtpM55ZJIhlejvqHgn0yXX/evtvdwjjZv0aShaqDZkWkGfsmjXbMlmWs2+J a.reuschenbach.puncernau@sap.com", "type": "ssh", "name": "Arturo_std"}], "hostname": "rel7.novalocal", "launch_index": 0, "public_keys": {"Arturo_std": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDeevePzO0cd4sV/iIvKsaVw08s8gC0BuzyJu7/hSt10tt4O46nZpqLLybB/K/telZ2vpgvVSxfpOMWngVCIegv72jVnEMpD1WIL40ackN7TRRfT6JLrSYwvoUYgE3CIk8TBZYf9OalXWWXVgYdHp/12u7NMOEvwBEdor9aWKB39ojnuA3s5guZt4fqBuOoaYE/32W4sL4TL7QqLBBqdGKjOZKvVZITKr4IPn4EDUVoGJ2hKS8f89kNSvmDe4tFgWSu7mohc9V7M8N0TkNz9bXIFe+9tZtpM55ZJIhlejvqHgn0yXX/evtvdwjjZv0aShaqDZkWkGfsmjXbMlmWs2+J a.reuschenbach.puncernau@sap.com"}, "project_id": "1d1ad583e98c4913a0226feac0f010f9", "name": "rel7"}`

func TestInstanceId(t *testing.T) {
	server := testTools(200, mockMetadata)
	defer server.Close()
	bakcupMetadataURL := metadataURL
	metadataURL = server.URL

	res := instanceID()
	if res != "16d2ef3a-da47-43bb-8204-e4fd6a5db689" {
		t.Error(fmt.Sprint("Expected to get the metadata id. ", res, " =! ", "16d2ef3a-da47-43bb-8204-e4fd6a5db689"))
	}

	// set back the metadatURL
	metadataURL = bakcupMetadataURL
}

func TestCommonNameFromFlag(t *testing.T) {
	server := testTools(200, mockMetadata)
	defer server.Close()
	bakcupMetadataURL := metadataURL
	metadataURL = server.URL

	// set a context
	flagSet := flag.NewFlagSet("local", 0)
	flagSet.String("common-name", "cnameTest", "local")
	ctx := cli.NewContext(nil, flagSet, getParentCtx())

	name := commonName(ctx)
	if name != "cnameTest" {
		t.Error(fmt.Sprint("Expected to get the cname from the flag. ", name, " != ", "cnameTest"))
	}

	// set back the metadatURL
	metadataURL = bakcupMetadataURL
}

func TestCommonNameFromMetadata(t *testing.T) {
	server := testTools(200, mockMetadata)
	defer server.Close()
	bakcupMetadataURL := metadataURL
	metadataURL = server.URL

	// set a context
	flagSet := flag.NewFlagSet("local", 0)
	ctx := cli.NewContext(nil, flagSet, getParentCtx())

	name := commonName(ctx)
	if name != "16d2ef3a-da47-43bb-8204-e4fd6a5db689" {
		t.Error(fmt.Sprint("Expected to get the cname from the metadata. ", name, " != ", "16d2ef3a-da47-43bb-8204-e4fd6a5db689"))
	}

	// set back the metadatURL
	metadataURL = bakcupMetadataURL
}

func TestCommonNameFromHostName(t *testing.T) {
	server := testTools(200, "no metadata given")
	defer server.Close()
	bakcupMetadataURL := metadataURL
	metadataURL = server.URL

	// set a context
	flagSet := flag.NewFlagSet("local", 0)
	ctx := cli.NewContext(nil, flagSet, getParentCtx())

	name := commonName(ctx)
	if name == "16d2ef3a-da47-43bb-8204-e4fd6a5db689" || name == "cnameTest" {
		t.Error(fmt.Sprint("Expected to get the cname from hostname. ", name, " == 16d2ef3a-da47-43bb-8204-e4fd6a5db689 || == cnameTest"))
	}

	hostname, _ := os.Hostname()
	if name != hostname {
		t.Error(fmt.Sprint("Expected to get the cname from hostname. ", name, " != ", hostname))
	}

	// set back the metadatURL
	metadataURL = bakcupMetadataURL
}

// private

func testTools(code int, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	return server
}
