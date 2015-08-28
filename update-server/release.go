package main

import (
	"io/ioutil"
	"path"

	//log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Release struct {
	Uid      string `json:"uid"`
	Filename string `json:"filename"`
	App      string `json:"app"`
	Os       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Date     string `json:"date"`
}

type Releases map[string]Release

func (releases *Releases) Read() error {
	releasesConfigPath := path.Join(buildsRootPath, "releases.yml")
	
	data, err := ioutil.ReadFile(releasesConfigPath)
	if err != nil {
		return err
	}

	*releases = make(Releases, 0)
	err = yaml.Unmarshal([]byte(data), &releases)
	if err != nil {
		return err
	}

	return nil
}