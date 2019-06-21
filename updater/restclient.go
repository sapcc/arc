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
	"regexp"

	"github.com/sapcc/arc/version"
)

var ErrorNoUpdateAvailable error = fmt.Errorf("no update available")

type Client struct {
	Endpoint string
}

type CheckParams struct {
	AppId string `json:"app_id"`
	OS    string `json:"os"`
	Arch  string `json:"arch"`
}

type CheckResult struct {
	CheckParams
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
	jsonUrl, err := c.buildPathJson(params)
	if err != nil {
		return nil, err
	}

	// make request
	resp, err := restCall(jsonUrl, "GET", url.Values{}, bytes.NewBufferString(""))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 404 {
		return nil, fmt.Errorf("%s - %d - %s", jsonUrl, resp.StatusCode, respBody)
	}

	// response body to struct
	var res CheckResult
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetUpdate(r *CheckResult) (io.ReadCloser, error) {
	isUrlAbsolute, err := isAbsouteUrl(r.Url)
	if err != nil {
		return nil, err
	}

	// check url absolut vs relativ
	download_url := r.Url
	if !isUrlAbsolute {
		download_url, err = c.buildPathBinary(*r)
		if err != nil {
			return nil, err
		}
	}

	resp, err := http.Get(download_url) // #nosec url is given by flag update-uri
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("got empty response body")
	}
	if resp.StatusCode >= 400 {
		if err = resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("got unexpected status code: %s. Can't close response body: %s", resp.Status, err)
		}
		return nil, fmt.Errorf("got unexpected status code: %s", resp.Status)
	}
	return resp.Body, nil
}

// private

func isAbsouteUrl(url string) (bool, error) {
	return regexp.MatchString("^(?:[a-z]+:)?//", url)
}

func (c *Client) buildPathJson(params CheckParams) (string, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, buildPath(params.AppId, params.OS, params.Arch), "latest.json")
	return u.String(), nil
}

func (c *Client) buildPathBinary(params CheckResult) (string, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, buildPath(params.AppId, params.OS, params.Arch), params.Url)
	return u.String(), nil
}

func buildPath(appId, os, arch string) string {
	return path.Join(appId, os, arch)
}

func restCall(urlPath, method string, params url.Values, body *bytes.Buffer) (*http.Response, error) {
	// set up the rest url
	u, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
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
