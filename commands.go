package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/sapcc/arc/arc"
	"github.com/sapcc/arc/commands"
	"github.com/sapcc/arc/fact"
	"github.com/sapcc/arc/fact/agents"
	arc_facts "github.com/sapcc/arc/fact/arc"
	"github.com/sapcc/arc/fact/host"
	"github.com/sapcc/arc/fact/memory"
	"github.com/sapcc/arc/fact/metadata"
	"github.com/sapcc/arc/fact/network"
	"github.com/sapcc/arc/transport"
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
			optApiUri,
			optCertUpdateInterval,
			optCertUpdateThreshold,
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
			optApiUri,
			optUpdateInterval,
			optCertUpdateInterval,
			optCertUpdateThreshold,
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
	{
		Name:        "renewcert",
		Usage:       cmdUsage["docs-commands-renewcert"],
		Description: cmdDescription["docs-commands-renewcert"],
		Flags: []cli.Flag{
			optTlsClientCert,
			optTlsClientKey,
			optTlsCaCert,
			optApiUri,
		},
		Before: func(c *cli.Context) error {
			return config.Load(c)
		},
		Action: cmdRenewCert,
	},
}

func cmdServer(c *cli.Context) {
	code, err := commands.CmdServer(c, config, appName)
	if err != nil {
		log.Error(err)
	}
	os.Exit(code)
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
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Error reading from stdin: %s", err)
		}
		payload = string(bytes)
	}

	if err := tp.Connect(); err != nil {
		log.Fatal("Error connecting to broker ", err)
	}
	defer tp.Disconnect()

	request, err := arc.CreateRequest(agent, action, config.Identity, c.String("identity"), c.Int("timeout"), payload)
	if err != nil {
		log.Fatal(err.Error())
	}

	msgChan, cancelSubscription := tp.SubscribeJob(request.RequestID)
	defer cancelSubscription()
	log.Infof("Sending request %s", request.RequestID)
	if err = tp.Request(request); err != nil {
		log.Fatal(err.Error())
	}
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
	store.AddSource(host.New(&config), 0)
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
	code, err := commands.Status(c)
	if err != nil {
		log.Error(err)
	}
	os.Exit(code)
}

func cmdStart(c *cli.Context) {
	code, err := commands.Start(c)
	if err != nil {
		log.Error(err)
	}
	os.Exit(code)
}

func cmdStop(c *cli.Context) {
	code, err := commands.Stop(c)
	if err != nil {
		log.Error(err)
	}
	os.Exit(code)
}

func cmdRestart(c *cli.Context) {
	code, err := commands.Restart(c)
	if err != nil {
		log.Error(err)
	}
	os.Exit(code)
}

func cmdRenewCert(c *cli.Context) {
	code, err := commands.RenewCert(c, &config)
	if err != nil {
		log.Error(err)
	}
	os.Exit(code)
}
