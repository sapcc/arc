package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCloudConfigInstallerlNoRenewCertURL(t *testing.T) {
	info := struct {
		Token        string
		SignURL      string
		EndpointURL  string
		UpdateURL    string
		RenewCertURL string
	}{
		Token:        "some_token",
		SignURL:      "some_url",
		EndpointURL:  "some_EndPoint",
		UpdateURL:    "some_Url",
		RenewCertURL: "",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := cloudConfigInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if strings.Contains(w.String(), "--renew-cert-uri") {
		t.Error("Cloud config install script should not contain flag --renew-cert-uri")
	}
}

func TestCloudConfigInstallerlRenewCertURL(t *testing.T) {
	info := struct {
		Token        string
		SignURL      string
		EndpointURL  string
		UpdateURL    string
		RenewCertURL string
	}{
		Token:        "some_token",
		SignURL:      "some_url",
		EndpointURL:  "some_EndPoint",
		UpdateURL:    "some_Url",
		RenewCertURL: "some_Cert_url",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := cloudConfigInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if !strings.Contains(w.String(), "--renew-cert-uri") {
		t.Error("Cloud config install script should contain flag --renew-cert-uri")
	}
	if !strings.Contains(w.String(), "some_Cert_url") {
		t.Error("Cloud config install script should contain flag value some_Cert_url")
	}
}

func TestShellScriptInstallerNoRenewCertURL(t *testing.T) {
	info := struct {
		Token        string
		SignURL      string
		EndpointURL  string
		UpdateURL    string
		RenewCertURL string
	}{
		Token:        "some_token",
		SignURL:      "some_url",
		EndpointURL:  "some_EndPoint",
		UpdateURL:    "some_Url",
		RenewCertURL: "",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := shellScriptInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if strings.Contains(w.String(), "--renew-cert-uri") {
		t.Error("Shell config install script should not contain flag --renew-cert-uri")
	}
}

func TestShellScriptInstallerRenewCertURL(t *testing.T) {
	info := struct {
		Token        string
		SignURL      string
		EndpointURL  string
		UpdateURL    string
		RenewCertURL string
	}{
		Token:        "some_token",
		SignURL:      "some_url",
		EndpointURL:  "some_EndPoint",
		UpdateURL:    "some_Url",
		RenewCertURL: "some_Cert_url",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := shellScriptInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if !strings.Contains(w.String(), "--renew-cert-uri") {
		t.Error("Shell config install script should contain flag --renew-cert-uri")
	}
	if !strings.Contains(w.String(), "some_Cert_url") {
		t.Error("Shell config install script should contain flag value some_Cert_url")
	}
}

func TestPowerShellScriptInstallerNoRenewCertURL(t *testing.T) {
	info := struct {
		Token        string
		SignURL      string
		EndpointURL  string
		UpdateURL    string
		RenewCertURL string
	}{
		Token:        "some_token",
		SignURL:      "some_url",
		EndpointURL:  "some_EndPoint",
		UpdateURL:    "some_Url",
		RenewCertURL: "",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := powershellScriptInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if strings.Contains(w.String(), "--renew-cert-uri") {
		t.Error("PowerShell config install script should not contain flag --renew-cert-uri")
	}
}

func TestPowerShellScriptInstallerRenewCertURL(t *testing.T) {
	info := struct {
		Token        string
		SignURL      string
		EndpointURL  string
		UpdateURL    string
		RenewCertURL string
	}{
		Token:        "some_token",
		SignURL:      "some_url",
		EndpointURL:  "some_EndPoint",
		UpdateURL:    "some_Url",
		RenewCertURL: "some_Cert_url",
	}
	//var w io.Writer
	var w bytes.Buffer
	err := powershellScriptInstaller.Execute(&w, info)
	if err != nil {
		t.Error("Expected to not have an error")
	}
	if !strings.Contains(w.String(), "--renew-cert-uri") {
		t.Error("PowerShell config install script should contain flag --renew-cert-uri")
	}
	if !strings.Contains(w.String(), "some_Cert_url") {
		t.Error("PowerShell config install script should contain flag value some_Cert_url")
	}
}
