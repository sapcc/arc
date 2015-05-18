package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"gitHub.***REMOVED***/monsoon/arc/server"
	"gitHub.***REMOVED***/monsoon/arc/transport"
)

var Commands = []cli.Command{
	{
		Name:   "server",
		Usage:  "Run the agent",
		Action: cmdServer,
	},
	{
		Name:   "execute",
		Usage:  "Remote execute action",
		Action: cmdExecute,
	},
}

func cmdServer(c *cli.Context) {
	doneChan := make(chan bool)

	tp, err := transport.New(config)
	if err != nil {
		log.Fatal(err)
	}
	server := server.New(tp)

	go server.Run(doneChan)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
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

func cmdExecute(c *cli.Context) {

	tp, err := transport.New(config)
	if err != nil {
		log.Fatal(err)
	}

	tp.Connect()
	defer tp.Disconnect()

	for {
		select {
		case <-time.After(1 * time.Second):
			return
		}
	}

}
