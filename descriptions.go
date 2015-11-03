package main

var cmdUsage = map[string]string{
	"docs-commands-facts":   `Discover and list facts on this system.`,
	"docs-commands-help":    `The "help" shows a list of commands or help for one command.`,
	"docs-commands":         `Arc is controlled via a very easy to use command-line interface (CLI).`,
	"docs-commands-init":    `The "init" command initializes server configuration.`,
	"docs-commands-list":    `The "list" command shows all available agents and their actions.`,
	"docs-commands-restart": `The "restart" command restart agent services.`,
	"docs-commands-run":     `The "arc run" execute an agent action [REFERENCE ACTIONS DOC PAGE FROM HERE] on a remote Arc server.`,
	"docs-commands-server":  `Run the Arc daemon.`,
	"docs-commands-start":   `The "start" command start agent service.`,
	"docs-commands-status":  `The "status" command gives the service status.`,
	"docs-commands-stop":    `The "stop" command stop agent services.`,
	"docs-commands-update":  `The "update" command check for new updates and update to the last version.`,
}

var cmdDescription = map[string]string{
	"docs-commands-facts": `Discover and list facts on this system. The "facts" command collects information from the system to provide a simple
and easy to understand view of the machine where Arc is running.`,
	"docs-commands-help": `The "help" command shows a list of commands or help for one specific command.`,
	"docs-commands": `Arc is controlled via a very easy to use command-line interface (CLI).
Arc is only a single command-line application: "arc". This application
then takes a subcommand such as "server" or "run". The complete list of
subcommands is in the navigation to the left.

The Arc CLI is a well-behaved command line application. In erroneous
cases, a non-zero exit status will be returned. It also responds to "-h" and "--help"
as you'd most likely expect.`,
	"docs-commands-init": `Coming soon...`,
	"docs-commands-list": `The "list" command shows all available agents and their actions. The payload requirements is action dependent
and it should be looked up in the corresponding implementation. Some payload examples can be found in the
[run command](/docs/commands/run.html) documentation. A complete explanation over the available agents can be found
in the [arc server](/docs/server/agents.html) documentation topic.`,
	"docs-commands-restart": `Coming soon...`,
	"docs-commands-run": `The "arc run" execute an agent action [REFERENCE ACTIONS DOC PAGE FROM HERE] on a remote Arc server.

Example: "arc run -endpoint tcp://localhost:1883 -identity darwin rpc version"

It prints the current version of the Arc running on the remote server with the identity "darwin".

    0.1.0-dev(ae07667)

Example: "arc run -endpoint tcp://localhost:1883 -identity darwin -payload "echo Script start; for i in {1..5}; do echo \$i; sleep 1s; done; echo Script done" execute script"

    Script start
    1
    2
    3
    4
    5
    Script done

It runs the script given as a payload on the remote server with identity "darwin". If the script gets to complicated or too long to give it as a payload, there is possibility to
write a bash script file and give it to the "run" command as the following example:

Example: "arc run -endpoint tcp://localhost:1883 -identity darwin -stdin execute script < script.sh"

    #!/bin/bash
    # script.sh

    echo Scritp start
    for i in {1..5}; do
    	echo $i
    	sleep 1s
    done
    echo Scritp done`,
	"docs-commands-server": `Run the Arc daemon.

Due to the power and complexity of this command, the Arc server is documented in its own section.
See the [Arc Server](/docs/server/basics.html) section for more information on how to use this command and the options it has.`,
	"docs-commands-start":  `Coming soon...`,
	"docs-commands-status": `Coming soon...`,
	"docs-commands-stop":   `Coming soon...`,
	"docs-commands-update": `The "update" command checks for the last version available, asks for user confirmation and triggers an update. When
the update is being triggered the existing Arc binary is replaced with the new one.`,
}
