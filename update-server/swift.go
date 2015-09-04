package main

import (
	"fmt"
	"errors"
	"github.com/ncw/swift"
)

func swiftTest() {
	// Create a connection	
	c, err := NewSwiftConnection("concourse", "secret", "http://automation.***REMOVED***:5000/v3", "monsooncc", "arc_releases_development")
	if err != nil {
		panic(err)
	}

	// Authenticate
	err = c.Connection.Authenticate()
	if err != nil {
		panic(err)
	}
	
	// check and create container
	err = c.CheckAndCreateContainer()
	if err != nil {
		panic(err)
	}	
	
	// save example file
	obj := SwiftObject{
		Name: "test_file",
		Content: []byte("Here is a string...."),
	}	
	err = obj.Save(c)
	if err != nil {
		panic(err)
	}
	
	// Get all files in the container
	filenames, err := c.GetAllContainerNames()
	if err != nil {
		panic(err)
	}
	fmt.Println(filenames)
	
	// Get file
	objBack := SwiftObject{}	
	data, err := objBack.Get(c, "test_file")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

type SwiftConnection struct {
	Connection  swift.Connection  `json:"connection"`
	Container   string 		`json:"container"`
}

func NewSwiftConnection(username, password, authUrl, domain string, container string) (*SwiftConnection, error) {
	if username == "" || password == "" || authUrl == "" || domain == "" || container == "" {
		return nil, errors.New("Not enough arguments in call NewSwiftConnection")
	}
	
	return &SwiftConnection{
		swift.Connection{
			UserName: username,
			ApiKey:   password,
			AuthUrl:  authUrl,
			Domain:   domain,
		},
		container,
	}, nil
}


type SwiftObject struct {
	Content     []byte 			 `json:"content"`
	Name        string       `json:"name"`
}

func (obj *SwiftObject) Get(c *SwiftConnection, name string) ([]byte, error) {
	data, err := c.Connection.ObjectGetBytes(c.Container, name)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (obj *SwiftObject) Save(c *SwiftConnection) error {	
	err := c.Connection.ObjectPutBytes(c.Container, obj.Name, obj.Content, "application/octet-stream")
	if err != nil {
		return err
	}
	return nil
}

func (c *SwiftConnection) GetAllContainerNames() ([]string, error) {
	names, err := c.Connection.ObjectNames(c.Container, nil)
	if err != nil {
		return nil, err
	}	
	return names, nil
}

func (c *SwiftConnection) CheckAndCreateContainer() error {
	_, _, err := c.Connection.Container(c.Container)
	if err == swift.ContainerNotFound {		
		err = c.Connection.ContainerCreate(c.Container, nil)
		if err != nil {
			return err
		}		
	} else if err != nil {
		return err
	}	
	return nil
}

func (c *SwiftConnection) DeleteContainer() error {
	err := c.Connection.ContainerDelete(c.Container)
	if err != nil {
		return err
	}
	return nil
}
 