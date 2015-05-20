package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/server"
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/updater"
)

var Commands = []cli.Command{
	{
		Name:   "server",
		Usage:  "Run the agent",
		Action: cmdServer,
	},
	{
		Name:  "execute",
		Usage: "Execute remote agent action",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "timeout, t",
				Usage: "timeout for executing the action",
				Value: 60,
			},
			cli.StringFlag{
				Name:  "identity, i",
				Usage: "target system",
				Value: "",
			},
			cli.StringFlag{
				Name:  "payload,p",
				Usage: "payload for action",
				Value: "",
			},
			cli.BoolFlag{
				Name:  "stdin,s",
				Usage: "read payload from stdin",
			},
		},
		Action: cmdExecute,
	},
	{
		Name:   "list",
		Usage:  "list available agents and actions",
		Action: cmdList,
	},
}

func cmdServer(c *cli.Context) {
	doneChan := make(chan bool)

	// Ticker containing a channel that will send the time with a period
	tickChan := time.NewTicker(time.Second * time.Duration(c.GlobalInt("update-interval")))
	// updater object
	up := updater.New(map[string]string{
		"version":   Version,
		"appName":   appName,
		"updateUri": c.GlobalString("update-uri"),
	})

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
		case <-tickChan.C:
			if !c.GlobalBool("no-auto-update") {
				go up.Update(tickChan)
			}
		}
	}

}

func cmdExecute(c *cli.Context) {

	tp, err := transport.New(config)
	if err != nil {
		log.Fatal(err)
	}

	registry := arc.AgentRegistry()

	agent := c.Args().Get(0)
	action := c.Args().Get(1)
	if registry.HasAction(agent, action) == false {
		log.Fatal("You need to provide a valid agent and action name")
	}

	if c.String("identity") == "" {
		log.Fatal("Target identity not given.")
	}

	if c.Int("timeout") < 1 {
		log.Fatal("timeout needs to be a positive integer")
	}

	payload := c.String("payload")

	if c.Bool("stdin") {
		if c.String("payload") != "" {
			log.Fatal("--stdin and --payload are mutually exclusive")
		}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			payload = scanner.Text()
		}
	}

	if err := tp.Connect(); err != nil {
		log.Fatal("Error connecting to broker ", err)
	}
	defer tp.Disconnect()

	request := arc.CreateRequest(agent, action, c.String("identity"), c.Int("timeout"), payload)

	msgChan, cancelSubscription := tp.SubscribeJob(request.RequestID)
	defer cancelSubscription()
	log.Infof("Sending request %s", request.RequestID)
	tp.Request(request)
	state := arc.Queued

	for {
		select {
		case reply := <-msgChan:
			log.Debug(reply)

			if state == arc.Queued && reply.State == arc.Executing {
				log.Infof("Job %s started executing", reply.RequestID)
			}
			state = reply.State

			if reply.Payload != "" {
				log.Info("Payload: ", reply.Payload)
			}
			if state == arc.Complete {
				log.Infof("Job %s completed successfully", reply.RequestID)
				return
			}
			if state == arc.Failed {
				log.Errorf("Job %s failed", reply.RequestID)
				return
			}

		case <-time.After(time.Duration(c.Int("timeout")) * time.Second):
			log.Warnf("Timeout waiting for job %s\n", request.RequestID)
			return
		}
	}

}

func cmdList(c *cli.Context) {

	registry := arc.AgentRegistry()

	fmt.Printf("  %-20sActions\n", "Agent")
	fmt.Println(strings.Repeat("-", 40))
	for _, agent := range registry.Agents() {
		fmt.Printf("  %-20s%s\n", agent, strings.Join(registry.Actions(agent), ","))
	}

}
