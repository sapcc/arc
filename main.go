package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	_ "gitHub.***REMOVED***/monsoon/onos/agents/rpc"
	"gitHub.***REMOVED***/monsoon/onos/server"
	"gitHub.***REMOVED***/monsoon/onos/transport"
)

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal("Configuration error: ", err.Error())
	}
	if printVersion {
		fmt.Printf("Onos %s\n", Version)
		os.Exit(0)
	}

	transport, err := transport.New(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	doneChan := make(chan bool)

	server := server.New(doneChan, transport)
	go server.Run()

	//setup signal handlers
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	log.Debug("Waiting for something to happen...")
	for {
		select {
		case s := <-signalChan:
			log.Info(fmt.Sprintf("Captured %v. Exiting...", s))
			server.Stop()
		case <-doneChan:
			os.Exit(0)
		}
	}
}
