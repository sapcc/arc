package server

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/onos/onos"
	"gitHub.***REMOVED***/monsoon/onos/transport"
	"golang.org/x/net/context"
)

type Server interface {
	Run(done chan<- bool)
	Stop()
}

type server struct {
	stopChan    chan bool
	transport   transport.Transport
	activeJobs  map[string]func()
	rootContext context.Context
	cancel      func()
}

func New(transport transport.Transport) Server {
	stopChan := make(chan bool)
	activeJobs := make(map[string]func())
	return &server{stopChan, transport, activeJobs, nil, nil}
}

func (s *server) Run(doneChan chan<- bool) {
	defer close(doneChan)

	s.transport.Connect()
	defer s.transport.Disconnect()
	incomingChan := s.transport.Subscribe()

	s.rootContext, s.cancel = context.WithCancel(context.Background())
	done := s.rootContext.Done()

	for {
		select {
		case <-done:
			log.Debug("Server was stopped")
			return
		case msg := <-incomingChan:
			go s.handleJob(msg)
		}
	}

}

func (s *server) Stop() {
	log.Info("Stopping Server")
	s.cancel()
}

func (s *server) handleJob(msg *onos.Request) {
	log.Infof("Dispatching message with requestID %s to agent %s\n", msg.RequestID, msg.Agent)
	jobContext, _ := context.WithTimeout(s.rootContext, time.Duration(msg.Timeout)*time.Second)

	outChan := make(chan *onos.Reply)
	go onos.ExecuteAction(jobContext, msg, outChan)

	for m := range outChan {
		s.transport.Reply(m)
	}
	log.Infof("Job %s completed\n", msg.RequestID)
}
