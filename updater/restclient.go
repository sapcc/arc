package updater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"gitHub.***REMOVED***/monsoon/arc/version"
)

var NoUpdateAvailable error = fmt.Errorf("No update available")

type Client struct {
	Endpoint string
}

type CheckParams struct {
	AppVersion string `json:"app_version"`
	AppId      string `json:"app_id"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
}

type CheckResult struct {
	Url      string `json:"url"`
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) CheckForUpdate(params CheckParams) (*CheckResult, error) {
	// prepare body with params
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// make request
	resp, err := restCall(c.Endpoint, "updates", "POST", url.Values{}, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// status 204 means no update available
	if resp.StatusCode == 204 {
		return nil, NoUpdateAvailable
	}

	// read body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%v - %s", resp.StatusCode, respBody)
	}

	// response body to struct
	var res CheckResult
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetUpdate(url string) (*io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return &resp.Body, nil
}

// private

func restCall(endpoint string, pathAction string, method string, params url.Values, body *bytes.Buffer) (*http.Response, error) {
	// set up the rest url
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, pathAction)
	u.RawQuery = params.Encode()

	// set up body
	var reqBody io.Reader
	if body != nil && body.Len() > 0 {
		reqBody = body
	}

	// set up the request
	httpclient := &http.Client{}
	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", fmt.Sprint("arc-updater/", version.String()))
	req.Header.Add("Content-Type", "application/json")

	// send the request
	resp, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
