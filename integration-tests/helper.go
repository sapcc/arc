package integrationTests

import (
	"bytes"
	"os"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path"
	"flag"
)

var apiServerFlag = flag.String("api-server", "http://localhost:3000", "integration-test")
var updateServerFlag = flag.String("update-server", "http://localhost:3001", "integration-test") 

type ServerType int

const (
	_ ServerType = iota
	ApiServer
	UpdateServer
)

type Client struct {
	Client        *http.Client
	ApiServerUrl  string
	UpdateServerUrl string
}

func NewTestClient() *Client {	
	// override flags if enviroment variable exists
	if os.Getenv("ARC_API_SERVER") != "" {
		apiServerUrl := os.Getenv("ARC_API_SERVER")
		apiServerFlag = &apiServerUrl
	}
	if os.Getenv("ARC_UPDATE_SERVER") != "" {
		updateServerUrl := os.Getenv("ARC_UPDATE_SERVER")
		updateServerFlag = &updateServerUrl
	}
	
	return &Client{
		Client:          &http.Client{},
		ApiServerUrl:    *apiServerFlag,
		UpdateServerUrl: *updateServerFlag,
	}
}

func (c *Client) Get(pathTo string, server ServerType) (string, *[]byte) {
	url := fmt.Sprint(c.serverUrl(server), path.Join("/", pathTo))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return "", nil
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

func (c *Client) Post(pathTo string, server ServerType, headers map[string]string, jsonBody []byte) (string, *[]byte) {
	url := fmt.Sprint(c.serverUrl(server), path.Join("/", pathTo))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
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

func (c *Client) serverUrl(s ServerType) string {
	switch s {
	case ApiServer:
		return c.ApiServerUrl
	case UpdateServer:
		return c.UpdateServerUrl
	}
	return ""
}
