package server

import (
	"gitHub.***REMOVED***/monsoon/onos/onos"
	"gitHub.***REMOVED***/monsoon/onos/transport"
)

type Server interface {
	Run()
	Stop()
}

type server struct {
	stopChan  <-chan bool
	doneChan  chan<- bool
	transport transport.Transport
}

func New(stopChan <-chan bool, doneChan chan<- bool, transport transport.Transport) Server {
	return &server{stopChan, doneChan, transport}
}

func (s *server) Run() {
	defer close(s.doneChan)

	messageHandler := func(msg onos.Message) {
	}

	s.transport.Connect()
	s.transport.Subscribe(messageHandler)

	for {
		select {
		case <-s.stopChan:
			break
		}
	}

	s.transport.Disconnect()

}

func (s *server) Stop() {
	close(s.stopChan)
}
