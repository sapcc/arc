package commands

import (
	"fmt"
	"os"
	"strings"

	"net/url"

	"github.com/codegangsta/cli"
	"gitHub.***REMOVED***/monsoon/arc/api-server/pki"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
)

// RenewCert download a new cert
func RenewCert(c *cli.Context, cfg *arc_config.Config) (int, error) {
	// check api renew cert uri
	renewCertURI, err := renewCertURI(c)
	if err != nil {
		return 1, err
	}

	// check hours left
	_, hoursLeft, err := pki.RenewCert(cfg, renewCertURI, int64(0))
	if err != nil {
		return 1, err
	}
	fmt.Printf("Cert expires in %d hours. \n", hoursLeft)

	// get the new cert
	success, _, err := pki.RenewCert(cfg, c.String("cert-renew-uri"), int64(87600)) // 10 years in hours
	if err != nil {
		return 1, err
	}

	if success {
		fmt.Println("Cert successfully downloaded")
	} else {
		fmt.Printf("Cert expires in %d hours. \n", hoursLeft)
	}

	return 0, nil
}

func renewCertURI(c *cli.Context) (string, error) {
	// check api renew cert uri
	uri := c.String("renew-cert-uri")
	uriType := 0
	if uri == "" {
		uri = os.Getenv("ARC_UPDATE_URI")
		uriType = 1
	}

	if uri == "" {
		return "", fmt.Errorf("No renew cert URI found")
	}

	// Parse the URL and ensure there are no errors.
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("Failed to parse uri: %v", err)
	}
	host := u.Host
	// in case of update uri split the url to get the last part
	if uriType == 1 {
		strSlice := strings.SplitN(host, ".", 2)
		host = strSlice[1]
	}

	// add path to the renew cert URI
	renewCertURL := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/v1/agents/renew",
	}

	return renewCertURL.String(), nil
}
