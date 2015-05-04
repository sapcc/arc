package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	_ "gitHub.***REMOVED***/monsoon/onos/agents/rpc"
	"gitHub.***REMOVED***/monsoon/onos/onos"
	"gitHub.***REMOVED***/monsoon/onos/transport"
)

func main() {
	flag.Parse()
	if printVersion {
		fmt.Printf("Onos %s\n", Version)
		os.Exit(0)
	}

	err := initConfig()
	if err != nil {
		log.Fatal("Configuration error: ", err.Error())
	}

	transport, err := transport.New(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	//stopChan := make(chan bool)
	doneChan := make(chan bool)
	errChan := make(chan error, 10)

	transport.Connect()
	transport.Subscribe(func(onos.Message) {})

	//setup signal handlers
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	log.Debug("Waiting for something to happen...")
	for {
		select {
		case err := <-errChan:
			log.Error(err.Error())
		case s := <-signalChan:
			log.Info(fmt.Sprintf("Captured %v. Exiting...", s))
			close(doneChan)
		case <-doneChan:
			os.Exit(0)
		}
	}
}
