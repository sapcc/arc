package server

import (
	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/onos/onos"
	"gitHub.***REMOVED***/monsoon/onos/transport"
)

type Server interface {
	Run()
	Stop()
}

type server struct {
	stopChan  chan bool
	doneChan  chan<- bool
	transport transport.Transport
}

func New(doneChan chan<- bool, transport transport.Transport) Server {
	stopChan := make(chan bool)
	return &server{stopChan, doneChan, transport}
}

func (s *server) Run() {
	defer close(s.doneChan)

	s.transport.Connect()
	msgChan := s.transport.Subscribe()

	for {
		select {
		case <-s.stopChan:
			log.Debug("Server received stop signal")
			s.transport.Disconnect()
			return
		case msg := <-msgChan:
			dispatchMessage(msg)
		}
	}

}

func (s *server) Stop() {
	log.Info("Stopping Server")
	close(s.stopChan)
}

func dispatchMessage(m *onos.Message) {
	log.Infof("Received message with requestID %s for agent %s\n", m.RequestID, m.Agent)

}
