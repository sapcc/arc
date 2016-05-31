package execute

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"gitHub.***REMOVED***/monsoon/arc/arc"
)

func createTarball(files map[string]string) []byte {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)

	// Create a new tar archive.
	tw := tar.NewWriter(gw)
	for name, body := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0755,
			Size: int64(len(body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			panic(err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			panic(err)
		}
	}
	tw.Close()
	gw.Close()

	return buf.Bytes()

}

func TestTarballAction(t *testing.T) {

	testserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write(createTarball(map[string]string{
			"run.sh": `#!/bin/sh
			printf "arg:%s,env:%s" "$1" "$VAR"`,
		}))
	}))
	defer testserver.Close()
	payload := fmt.Sprintf(`{"path": "run.sh", "url":"%s", "arguments": ["arg1"], "environment": {"VAR":"env1"}}`, testserver.URL)
	req, err := arc.CreateRequest("execute", "tarball", "sender", "identity", 60, payload)
	if err != nil {
		t.Fatal(err)
	}
	out := make(chan *arc.Reply, 10)
	job := arc.NewJob("identity", req, out)

	agent := &executeAgent{}

	_, err = agent.TarballAction(context.Background(), job)
	if err != nil {
		t.Fatal(err)
	}

	// inital empty heartbeat
	<-out
	//command output
	reply := <-out
	expected := fmt.Sprintf("arg:%s,env:%s", "arg1", "env1")
	if reply.Payload != expected {
		t.Fatalf("Expected %v, got %v", expected, reply.Payload)
	}
}

func TestTarballHTTPError(t *testing.T) {
	testserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))

	defer testserver.Close()
	payload := fmt.Sprintf(`{"url":"%s"}`, testserver.URL)
	req, err := arc.CreateRequest("execute", "tarball", "sender", "identity", 60, payload)
	if err != nil {
		t.Fatal(err)
	}
	out := make(chan *arc.Reply, 10)
	job := arc.NewJob("identity", req, out)
	agent := &executeAgent{}
	_, err = agent.TarballAction(context.Background(), job)
	if err == nil {
		t.Fatal("Action should fail when the download errors.")
	}
	if !strings.Contains(err.Error(), testserver.URL) {
		t.Fatalf("Expected error message containing %v, got %v", testserver.URL, err.Error())
	}

}
