package integration_tests

import (
	"bytes"
	"os"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path"
)

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
	apiServerUrl := "http://localhost:3000"
	updateServerUrl := "http://localhost:3001"
	
	if os.Getenv("ARC_API_SERVER") != "" {
		apiServerUrl = os.Getenv("ARC_API_SERVER")
	}
	if os.Getenv("ARC_UPDATE_SERVER") != "" {
		updateServerUrl = os.Getenv("ARC_UPDATE_SERVER")
	}
	
	return &Client{
		Client:          &http.Client{},
		ApiServerUrl:    apiServerUrl,
		UpdateServerUrl: updateServerUrl,
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

// func main() {
// 	client := NewTestClient()
//
// 	to := "darwin"
// 	timeout := 60
// 	agent := "execute"
// 	action := "script"
// 	payload := `"payload":"echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""`
// 	data := fmt.Sprintf(`{"to":%q,"timeout":%v,"agent":%q,"action":%q,"payload":%q}`, to, timeout, agent, action, payload)
// 	jsonStr := []byte(data)
// 	statusCode, body = client.Post("/jobs", ApiServer, nil, jsonStr)
// 	fmt.Println(statusCode)
// 	fmt.Println(string(body))
// }
