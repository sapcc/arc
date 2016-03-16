package server

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/fact"
	"gitHub.***REMOVED***/monsoon/arc/fact/agents"
	arc_facts "gitHub.***REMOVED***/monsoon/arc/fact/arc"
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
	config      arc_config.Config
	transport   transport.Transport
	activeJobs  map[string]func()
	jobsMutex   sync.Mutex
	rootContext context.Context
	cancel      func()
	wg          sync.WaitGroup
	factStore   *fact.Store
}

func New(config arc_config.Config, transport transport.Transport) Server {
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
	defer func() {
		log.Debug("Disconnecting transport")
		s.transport.Disconnect()
	}()
	incomingChan, cancelSubscription := s.transport.Subscribe(s.config.Identity)
	defer func() {
		log.Debug("Cancelling subscriptions")
		cancelSubscription()
	}()

	s.rootContext, s.cancel = context.WithCancel(context.Background())
	done := s.rootContext.Done()

	s.factStore = s.setupFactStore()
	factUpdates := s.factStore.Updates()

	for {
		select {
		case <-done:
			log.Debug("Exiting sever run loop")
			return
		case update := <-factUpdates:
			log.Debug("Processing fact update")
			j, err := json.Marshal(update)
			if err == nil {
				if req, err := arc.CreateRegistration(s.config.Organization, s.config.Project, s.config.Identity, string(j)); err == nil {
					s.transport.Registration(req)
				} else {
					log.Warn("Failed to create registration message", err)
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
	jobContext, cancel := arc.NewJobContext(s.rootContext, time.Duration(msg.Timeout)*time.Second, s.factStore)
	//Accroding to https://www.youtube.com/watch?v=3EW1hZ8DVyw
	//not canceling a context is a memory leak.
	defer cancel()

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
	job := arc.NewJob(s.config.Identity, msg, outChan)
	go arc.ExecuteAction(jobContext, job)

	for m := range outChan {
		s.transport.Reply(m)
	}
	log.Infof("Job %s completed", msg.RequestID)
}

func (s *server) setupFactStore() *fact.Store {
	store := fact.NewStore()
	store.AddSource(host.New(), 1*time.Minute)
	store.AddSource(memory.New(), 1*time.Minute)
	store.AddSource(network.New(), 1*time.Minute)
	store.AddSource(arc_facts.New(s.config), 0)
	store.AddSource(agents.New(), 1*time.Minute)
	return store
}
