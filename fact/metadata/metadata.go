package metadata

import (
	"bufio"
	"encoding/json"
	"fmt"
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

	data := metaDataInfo(client)
	if data != nil {
		if data.UUID != "" {
			facts["metadata_uuid"] = data.UUID
		}
		if data.AvailabilityZone != "" {
			facts["metadata_availability_zone"] = data.AvailabilityZone
		}
		if data.Name != "" {
			facts["metadata_name"] = data.Name
		}
	}

	ips := floatingIP(client)
	if len(ips) > 0 {
		facts["metadata_public_ipv4"] = strings.Join(ips, ",")
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

	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		if scanner.Text() != "" {
			ips = append(ips, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Warnf(fmt.Sprint("Error scanning ipv4s. ", err.Error()))
		return ips
	}

	return ips
}

type metaData struct {
	UUID             string `json:"uuid"`
	AvailabilityZone string `json:"availability_zone"`
	Name             string `json:"name"`
}

// InstanceID returns the instance id from the metadata
func metaDataInfo(client http.Client) *metaData {
	r, err := client.Get(metadataURL)
	if err != nil {
		log.Warnf(fmt.Sprint("Error requesting metadata. ", err.Error()))
		return nil
	}
	defer r.Body.Close()

	var metadata = new(metaData)
	err = json.NewDecoder(r.Body).Decode(metadata)
	if err != nil {
		log.Warnf(fmt.Sprint("Error parsing metadata. ", err.Error()))
		return nil
	}

	return metadata
}
