package server

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/fact"
	"gitHub.***REMOVED***/monsoon/arc/fact/host"
	"gitHub.***REMOVED***/monsoon/arc/fact/memory"
	"gitHub.***REMOVED***/monsoon/arc/fact/network"
	"gitHub.***REMOVED***/monsoon/arc/transport"
)

type Server interface {
	Run()
	Stop()
	GracefulShutdown()
	Done() <-chan struct{}
}

type server struct {
	doneChan    chan struct{}
	config      arc.Config
	transport   transport.Transport
	activeJobs  map[string]func()
	jobsMutex   sync.Mutex
	rootContext context.Context
	cancel      func()
	wg          sync.WaitGroup
}

func New(config arc.Config, transport transport.Transport) Server {
	return &server{
		doneChan:   make(chan struct{}),
		transport:  transport,
		activeJobs: make(map[string]func()),
		config:     config,
	}
}

func (s *server) Done() <-chan struct{} {
	return s.doneChan
}

func (s *server) GracefulShutdown() {
	s.jobsMutex.Lock()
	runningJobs := len(s.activeJobs)
	s.jobsMutex.Unlock()
	log.Infof("Graceful shutdown triggered. Waiting for %d jobs to finish processing.", runningJobs)
	go func() {
		s.wg.Wait()
		log.Info("No jobs pending, shutting down")
		s.Stop()
	}()
}

func (s *server) Run() {
	log.Infof("Starting server. Pid %d", os.Getpid())
	defer log.Info("Server stopped")
	defer close(s.doneChan)

	s.transport.Connect()
	defer s.transport.Disconnect()
	incomingChan, cancelSubscription := s.transport.Subscribe(s.config.Identity)
	defer cancelSubscription()

	s.rootContext, s.cancel = context.WithCancel(context.Background())
	done := s.rootContext.Done()

	facts := setupFactStore()

	for {
		select {
		case <-done:
			return
		case update := <-facts.Updates():
			j, err := json.Marshal(update)
			if err == nil {
				if req, err := arc.CreateRegistrationMessage(s.config.Identity, string(j)); err == nil {
					s.transport.Request(req)
				} else {
					log.Warn("Failed to create registratrion request ", err)
				}
			} else {
				log.Warn("Failed to serialize fact update: ", err)
			}
		case msg := <-incomingChan:
			go s.handleJob(msg)
		}
	}

}

func (s *server) Stop() {
	log.Info("Stopping Server")
	s.cancel()
}

func (s *server) handleJob(msg *arc.Request) {
	log.Infof("Dispatching message with requestID %s to agent %s", msg.RequestID, msg.Agent)
	jobContext, cancel := context.WithTimeout(s.rootContext, time.Duration(msg.Timeout)*time.Second)

	s.wg.Add(1)
	defer s.wg.Done()

	//save a reference to the cancel method of the job context
	s.jobsMutex.Lock()
	s.activeJobs[msg.RequestID] = cancel
	s.jobsMutex.Unlock()
	defer func() {
		s.jobsMutex.Lock()
		delete(s.activeJobs, msg.RequestID)
		s.jobsMutex.Unlock()
	}()

	outChan := make(chan *arc.Reply)
	go arc.ExecuteAction(jobContext, s.config.Identity, msg, outChan)

	for m := range outChan {
		s.transport.Reply(m)
	}
	log.Infof("Job %s completed", msg.RequestID)
}

func setupFactStore() *fact.Store {
	s := fact.NewStore()
	s.AddSource(host.New(), 1*time.Minute)
	s.AddSource(memory.New(), 1*time.Minute)
	s.AddSource(network.New(), 1*time.Minute)
	return s
}
