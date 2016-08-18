package integrationTests

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
)

var apiServer = flag.String("api-server", "http://localhost:3000", "integration-test")
var updateServer = flag.String("update-server", "http://localhost:3001", "integration-test")
var token = flag.String("token", "", "Valid auth token")
var keystoneEndpoint = flag.String("keystone-endpoint", "", "Authentication endpoint for aquiring auth tokens")
var username = flag.String("username", "arc_test", "Username for keystone authentication")
var password = flag.String("password", "", "Password")
var project = flag.String("project", "", "(keystone) project name")
var domain = flag.String("domain", "", "(keystone) domain name")

var GITCOMMIT = "HEAD"

func TestMain(m *testing.M) {
	fmt.Printf("Git Revision of tests: %s\n", GITCOMMIT)
	os.Exit(m.Run())
}

type ServerType int

const (
	_ ServerType = iota
	ApiServer
)

type Client struct {
	Client       *http.Client
	ApiServerUrl string
	Token        string
}

func NewTestClient() (*Client, error) {
	// override flags if enviroment variable exists
	if e := os.Getenv("API_SERVER"); e != "" {
		apiServer = &e
	}
	if e := os.Getenv("KEYSTONE_ENDPOINT"); e != "" {
		keystoneEndpoint = &e
	}
	if e := os.Getenv("USERNAME"); e != "" {
		username = &e
	}
	if e := os.Getenv("PASSWORD"); e != "" {
		password = &e
	}
	if e := os.Getenv("PROJECT"); e != "" {
		project = &e
	}
	if e := os.Getenv("DOMAIN"); e != "" {
		domain = &e
	}
	if e := os.Getenv("TOKEN"); e != "" {
		token = &e
	}

	if *token == "" && *keystoneEndpoint != "" {
		var err error
		var authToken string
		if authToken, err = getToken(); err != nil {
			return nil, fmt.Errorf("Failed to get token from keystone: %s ", err)
		}
		token = &authToken
	}

	c := &Client{
		Client:       &http.Client{},
		ApiServerUrl: *apiServer,
		Token:        *token,
	}
	return c, nil
}

func (c *Client) Get(pathTo string, server ServerType) (string, *[]byte) {
	url := fmt.Sprint(c.serverUrl(server), pathTo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return "", nil
	}
	if c.Token != "" {
		req.Header.Add("X-Auth-Token", c.Token)
	} else {
		req.Header.Add("X-Identity-Status", `Confirmed`)
		req.Header.Add("X-Project-Id", *project)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Error(err)
		return "", nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", nil
	}

	return resp.Status, &body
}

func (c *Client) GetApiV1(pathTo string, server ServerType) (string, *[]byte) {
	return c.Get(path.Join("/api/v1/", pathTo), server)
}

func (c *Client) Post(pathTo string, server ServerType, headers map[string]string, jsonBody []byte) (string, *[]byte) {
	url := fmt.Sprint(c.serverUrl(server), pathTo)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	if c.Token != "" {
		req.Header.Add("X-Auth-Token", c.Token)
	} else {
		req.Header.Add("X-Identity-Status", `Confirmed`)
		req.Header.Add("X-Project-Id", *project)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Error(err)
		return "", nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", nil
	}

	return resp.Status, &body
}

func (c *Client) PostApiV1(pathTo string, server ServerType, headers map[string]string, jsonBody []byte) (string, *[]byte) {
	return c.Post(path.Join("/api/v1/", pathTo), server, headers, jsonBody)
}

func (c *Client) serverUrl(s ServerType) string {
	switch s {
	case ApiServer:
		return c.ApiServerUrl
	}
	return ""
}

func getToken() (string, error) {
	requestBody := `
	{ "auth": {
    "identity": {
      "methods": ["password"],
      "password": {
        "user": {
          "name": {{ .username }},
          "domain": { "name": {{ .domain }} },
          "password": {{ .password }} 
        }
      }
    },
    "scope": {
      "project": {
        "name": {{ .project }},
        "domain": { "name": {{ .domain }} }
      }
    }
  }
}
	`
	jsonEscape := func(s string) string { r, _ := json.Marshal(s); return string(r) }

	requestTemplate := template.Must(template.New("auth").Parse(requestBody))

	var buf bytes.Buffer
	err := requestTemplate.Execute(
		&buf,
		map[string]string{
			"username": jsonEscape(*username),
			"password": jsonEscape(*password),
			"project":  jsonEscape(*project),
			"domain":   jsonEscape(*domain),
		},
	)
	if err != nil {
		return "", err
	}

	c := http.Client{Timeout: 5 * time.Second}
	resp, err := c.Post(*keystoneEndpoint+"/auth/tokens?nocatalog", "application/json", &buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", errors.New(resp.Status)
	}
	return resp.Header.Get("X-Subject-Token"), nil
}
