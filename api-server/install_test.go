package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCloudConfigInstallerlNoApiURL(t *testing.T) {
	info := struct {
		Token       string
		SignURL     string
		EndpointURL string
		UpdateURL   string
		ApiURL      string
	}{
		Token:       "some_token",
		SignURL:     "some_url",
		EndpointURL: "some_EndPoint",
		UpdateURL:   "some_Url",
		ApiURL:      "",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := cloudConfigInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if strings.Contains(w.String(), "--api-uri") {
		t.Error("Cloud config install script should not contain flag --api-uri")
	}
}

func TestCloudConfigInstallerlApiURL(t *testing.T) {
	info := struct {
		Token       string
		SignURL     string
		EndpointURL string
		UpdateURL   string
		ApiURL      string
	}{
		Token:       "some_token",
		SignURL:     "some_url",
		EndpointURL: "some_EndPoint",
		UpdateURL:   "some_Url",
		ApiURL:      "some_Cert_url",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := cloudConfigInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if !strings.Contains(w.String(), "--api-uri") {
		t.Error("Cloud config install script should contain flag --api-uri")
	}
	if !strings.Contains(w.String(), "some_Cert_url") {
		t.Error("Cloud config install script should contain flag value some_Cert_url")
	}
}

func TestShellScriptInstallerNoApiURL(t *testing.T) {
	info := struct {
		Token       string
		SignURL     string
		EndpointURL string
		UpdateURL   string
		ApiURL      string
	}{
		Token:       "some_token",
		SignURL:     "some_url",
		EndpointURL: "some_EndPoint",
		UpdateURL:   "some_Url",
		ApiURL:      "",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := shellScriptInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if strings.Contains(w.String(), "--api-uri") {
		t.Error("Shell config install script should not contain flag --api-uri")
	}
}

func TestShellScriptInstallerApiURL(t *testing.T) {
	info := struct {
		Token       string
		SignURL     string
		EndpointURL string
		UpdateURL   string
		ApiURL      string
	}{
		Token:       "some_token",
		SignURL:     "some_url",
		EndpointURL: "some_EndPoint",
		UpdateURL:   "some_Url",
		ApiURL:      "some_Cert_url",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := shellScriptInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if !strings.Contains(w.String(), "--api-uri") {
		t.Error("Shell config install script should contain flag --api-uri")
	}
	if !strings.Contains(w.String(), "some_Cert_url") {
		t.Error("Shell config install script should contain flag value some_Cert_url")
	}
}

func TestPowerShellScriptInstallerNoApiURL(t *testing.T) {
	info := struct {
		Token       string
		SignURL     string
		EndpointURL string
		UpdateURL   string
		ApiURL      string
	}{
		Token:       "some_token",
		SignURL:     "some_url",
		EndpointURL: "some_EndPoint",
		UpdateURL:   "some_Url",
		ApiURL:      "",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := powershellScriptInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if strings.Contains(w.String(), "--api-uri") {
		t.Error("PowerShell config install script should not contain flag --api-uri")
	}
}

func TestPowerShellScriptInstallerApiURL(t *testing.T) {
	info := struct {
		Token       string
		SignURL     string
		EndpointURL string
		UpdateURL   string
		ApiURL      string
	}{
		Token:       "some_token",
		SignURL:     "some_url",
		EndpointURL: "some_EndPoint",
		UpdateURL:   "some_Url",
		ApiURL:      "some_Cert_url",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := powershellScriptInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if !strings.Contains(w.String(), "--api-uri") {
		t.Error("PowerShell config install script should contain flag --api-uri")
	}
	if !strings.Contains(w.String(), "some_Cert_url") {
		t.Error("PowerShell config install script should contain flag value some_Cert_url")
	}
}
