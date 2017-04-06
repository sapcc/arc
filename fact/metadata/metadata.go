package metadata

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/common/log"
)

type Source struct{}

var metadataURL = "http://169.254.169.254/openstack/latest/meta_data.json"
var ipv4Url = "http://169.254.169.254/latest/meta-data/public-ipv4"

func New() Source {
	return Source{}
}

func (h Source) Name() string {
	return "metadata"
}

func (h Source) Facts() (map[string]interface{}, error) {
	facts := make(map[string]interface{})

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	uuid := instanceID(client)
	if uuid != "" {
		facts["metadata_uuid"] = uuid
	}

	ips := floatingIP(client)
	if len(ips) > 0 {
		facts["metadata_ipv4"] = strings.Join(ips, ",")
	}

	return facts, nil
}

func floatingIP(client http.Client) []string {
	ips := []string{}
	r, err := client.Get(ipv4Url)
	if err != nil {
		log.Warnf(fmt.Sprint("Error requesting metadata. ", err.Error()))
		return ips
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ips
	}
	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		ips = append(ips, scanner.Text())
	}

	return ips
}

type metaDataID struct {
	UUID string `json:"uuid"`
}

// InstanceID returns the instance id from the metadata
func instanceID(client http.Client) string {
	r, err := client.Get(metadataURL)
	if err != nil {
		log.Warnf(fmt.Sprint("Error requesting metadata. ", err.Error()))
		return ""
	}
	defer r.Body.Close()

	var metadata = new(metaDataID)
	err = json.NewDecoder(r.Body).Decode(metadata)
	if err != nil {
		log.Warnf(fmt.Sprint("Error parsing metadata. ", err.Error()))
		return ""
	}

	return metadata.UUID
}
