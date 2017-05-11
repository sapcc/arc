package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/commands"
	"gitHub.***REMOVED***/monsoon/arc/fact"
	"gitHub.***REMOVED***/monsoon/arc/fact/agents"
	arc_facts "gitHub.***REMOVED***/monsoon/arc/fact/arc"
	"gitHub.***REMOVED***/monsoon/arc/fact/host"
	"gitHub.***REMOVED***/monsoon/arc/fact/memory"
	"gitHub.***REMOVED***/monsoon/arc/fact/metadata"
	"gitHub.***REMOVED***/monsoon/arc/fact/network"
	"gitHub.***REMOVED***/monsoon/arc/server"
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/updater"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

var cliCommands = []cli.Command{
	{
		Name:        "server",
		Usage:       cmdUsage["docs-commands-server"],
		Description: cmdDescription["docs-commands-server"],
		Flags: []cli.Flag{
			optTransport,
			optEndpoint,
			optTlsClientCert,
			optTlsClientKey,
			optTlsCaCert,
			optUpdateUri,
			optUpdateInterval,
		},
		Before: func(c *cli.Context) error {
			return config.Load(c)
		},
		Action: cmdServer,
	},
	{
		Name:        "run",
		Usage:       cmdUsage["docs-commands-run"],
		Description: cmdDescription["docs-commands-run"],
		Flags: []cli.Flag{
			optTransport,
			optEndpoint,
			optTlsClientCert,
			optTlsClientKey,
			optTlsCaCert,
			optTimeout,
			optIdentity,
			optPayload,
			optStdin,
		},
		Before: func(c *cli.Context) error {
			return config.Load(c)
		},
		Action: cmdExecute,
	},
	{
		Name:        "list",
		Usage:       cmdUsage["docs-commands-list"],
		Description: cmdDescription["docs-commands-list"],
		Action:      cmdList,
	},
	{
		Name:        "facts",
		Usage:       cmdUsage["docs-commands-facts"],
		Description: cmdDescription["docs-commands-facts"],
		Flags: []cli.Flag{
			optTlsClientCert,
			optTlsClientKey,
			optTlsCaCert,
		},
		Before: func(c *cli.Context) error {
			return config.Load(c)
		},
		Action: cmdFacts,
	},
	{
		Name:        "update",
		Usage:       cmdUsage["docs-commands-update"],
		Description: cmdDescription["docs-commands-update"],
		Flags: []cli.Flag{
			optForce,
			optUpdateUri,
			optNoUpdate,
		},
		Action: cmdUpdate,
	},
	{
		Name:        "init",
		Usage:       cmdUsage["docs-commands-init"],
		Description: cmdDescription["docs-commands-init"],
		Action:      cmdInit,
		Flags: []cli.Flag{
			optTransport,
			optEndpoint,
			optTlsClientCert,
			optTlsClientKey,
			optTlsCaCert,
			optUpdateUri,
			optUpdateInterval,
			optRegistrationUrl,
			optInstallDir,
			optCommonName,
		},
	},
	{
		Name:        "status",
		Usage:       cmdUsage["docs-commands-status"],
		Description: cmdDescription["docs-commands-status"],
		Action:      cmdStatus,
		Flags: []cli.Flag{
			optInstallDir,
		},
	},
	{
		Name:        "start",
		Usage:       cmdUsage["docs-commands-start"],
		Description: cmdDescription["docs-commands-start"],
		Action:      cmdStart,
		Flags: []cli.Flag{
			optInstallDir,
		},
	},
	{
		Name:        "stop",
		Usage:       cmdUsage["docs-commands-stop"],
		Description: cmdDescription["docs-commands-stop"],
		Action:      cmdStop,
		Flags: []cli.Flag{
			optInstallDir,
		},
	},
	{
		Name:        "restart",
		Usage:       cmdUsage["docs-commands-restart"],
		Description: cmdDescription["docs-commands-restart"],
		Action:      cmdRestart,
		Flags: []cli.Flag{
			optInstallDir,
		},
	},
}

func cmdServer(c *cli.Context) {
	log.Infof("Starting server version %s. identity: %s, project: %s, organization: %s", version.Version, config.Identity, config.Project, config.Organization)

	// update object and ticker
	var up *updater.Updater
	tickChan := time.NewTicker(time.Second * time.Duration(c.Int("update-interval")))
	if c.String("update-uri") != "" {
		// create update object
		up = updater.New(map[string]string{
			"version":   version.Version,
			"appName":   appName,
			"updateUri": c.String("update-uri"),
		})
		log.Infof("Updater setup with interval %v, version %q, app name %q and update uri %q", c.Int("update-interval"), version.Version, appName, c.String("update-uri"))
	} else {
		// ticker will be stoped if no update uri is given
		tickChan.Stop()
	}

	tp, err := transport.New(config, true)
	if err != nil {
		log.Fatal(err)
	}

	if err := tp.Connect(); err != nil {
		log.Fatalf("Failed to connect to broker: %s", err)
	}
	server := server.New(config, tp)

	go server.Run()

	gracefulChan := make(chan os.Signal, 1)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(gracefulChan, syscall.SIGTERM)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case s := <-shutdownChan:
			log.Infof("Captured %v", s)
			server.Stop()
		case s := <-gracefulChan:
			log.Infof("Captured %v", s)
			server.GracefulShutdown()
		case <-server.Done():
			os.Exit(0)
		case <-tickChan.C:
			go func() {
				success, err := up.CheckAndUpdate()
				if err != nil {
					log.Error(err)
				}
				if success {
					server.GracefulShutdown()
					tickChan.Stop()
				}
			}()
		}
	}

}

func cmdExecute(c *cli.Context) {
	tp, err := transport.New(config, false)
	if err != nil {
		log.Fatal(err)
	}

	registry := arc.AgentRegistry()

	agent := c.Args().Get(0)
	action := c.Args().Get(1)
	if !registry.HasAction(agent, action) {
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
		bytes, _ := ioutil.ReadAll(os.Stdin)
		payload = string(bytes)
	}

	if err := tp.Connect(); err != nil {
		log.Fatal("Error connecting to broker ", err)
	}
	defer tp.Disconnect()

	request, err := arc.CreateRequest(agent, action, config.Identity, c.String("identity"), c.Int("timeout"), payload)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

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
				fmt.Print(reply.Payload)
			}
			if state == arc.Complete {
				log.Infof("Job %s completed successfully", reply.RequestID)
				return
			}
			if state == arc.Failed {
				log.Warnf("Job %s failed", reply.RequestID)
				exitCode = 1
				return
			}

		case <-time.After(time.Duration(c.Int("timeout")+2) * time.Second):
			log.Warnf("Timeout waiting for job %s\n", request.RequestID)
			exitCode = 1
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

func cmdFacts(c *cli.Context) {
	store := fact.NewStore()
	store.AddSource(host.New(), 0)
	store.AddSource(memory.New(), 0)
	store.AddSource(network.New(), 0)
	store.AddSource(arc_facts.New(config), 0)
	store.AddSource(agents.New(), 0)
	store.AddSource(metadata.New(true), 0)
	j, err := json.MarshalIndent(store.Facts(), " ", "  ")
	if err != nil {
		log.Warnf("Failed to generate json: %s", err)
		exitCode = 1
		return
	}
	fmt.Println(string(j))
}

func cmdUpdate(c *cli.Context) {
	code, err := commands.Update(c, map[string]interface{}{"appName": appName})
	if err != nil {
		log.Error(err)
	}
	os.Exit(code)
}

func cmdInit(c *cli.Context) {
	code, err := commands.Init(c, appName)
	if err != nil {
		log.Error(err)
	}
	os.Exit(code)
}

func cmdStatus(c *cli.Context) {
	code, _ := commands.Status(c)
	os.Exit(code)
}

func cmdStart(c *cli.Context) {
	code, _ := commands.Start(c)
	os.Exit(code)
}

func cmdStop(c *cli.Context) {
	code, _ := commands.Stop(c)
	os.Exit(code)
}

func cmdRestart(c *cli.Context) {
	code, _ := commands.Restart(c)
	os.Exit(code)
}
