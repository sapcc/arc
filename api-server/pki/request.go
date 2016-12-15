package pki

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Scheme returns the request schema
func Scheme(req *http.Request) string {

	if len(req.Header["X-Forwarded-Ssl"]) > 0 && req.Header["X-Forwarded-Ssl"][0] == "on" {
		return "https"
	}
	if len(req.Header["X-Forwarded-Scheme"]) > 0 {
		return req.Header["X-Forwarded-Scheme"][0]
	}
	if len(req.Header["X-Forwarded-Proto"]) > 0 {
		return strings.Split(req.Header["X-Forwarded-Proto"][0], ",")[0]
	}
	if req.TLS != nil {
		return "https"
	}
	return "http"

}

// HostWithPort returns the request host and port
func HostWithPort(req *http.Request) string {
	if len(req.Header["X-Forwarded-Host"]) > 0 {
		forwarded_hosts := regexp.MustCompile(`,\s?`).Split(req.Header["X-Forwarded-Host"][0], -1)
		return forwarded_hosts[len(forwarded_hosts)-1]
	}
	return req.Host
}

// Host returns the request host
func Host(req *http.Request) string {
	return strings.Split(HostWithPort(req), ":")[0]
}

// Port returns the request port
func Port(req *http.Request) int {
	if parts := strings.Split(HostWithPort(req), ":"); len(parts) > 1 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	if len(req.Header["X-Forwarded-Port"]) > 0 {
		port, _ := strconv.Atoi(req.Header["X-Forwarded-Port"][0])
		return port
	}
	return defaultPorts(Scheme(req))
}

func defaultPorts(scheme string) int {
	if scheme == "https" {
		return 443
	}
	return 80
}
