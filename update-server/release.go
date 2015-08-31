package main

import (
	"io/ioutil"
	"path"
	"os"
	"io"
	"bytes"

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
	// configuration file does not exist
	if _, err := os.Stat(releasesConfigPath()); os.IsNotExist(err) {
		return nil		
	}
	
	// read config file
	data, err := ioutil.ReadFile(releasesConfigPath())
	if err != nil {
		return err
	}

	// transform to releases strucs
	*releases = make(Releases, 0)
	err = yaml.Unmarshal([]byte(data), &releases)
	if err != nil {
		return err
	}

	return nil
}

func (releases *Releases) Update(key string, release Release) error {
	err := releases.Read()
	if err != nil {
		return err
	}
	
	(*releases)[key] = release
	
	// transform to string
	text, err := yaml.Marshal(&releases)
	if err != nil {
		return err
	}
	
	// create the file
	out, err := os.Create(releasesConfigPath())
	if err != nil {
		return err
	}
	defer out.Close()
	
	// copy yml to file 
	_, err = io.Copy(out, bytes.NewBuffer(text))
	if err != nil {
		return err
	}
	
	return nil
}

// private

func releasesConfigPath() string {
	return path.Join(buildsRootPath, "releases.yml")
}