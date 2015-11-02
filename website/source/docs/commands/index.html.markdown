---
layout: "docs"
page_title: "Commands"
sidebar_current: "docs-commands"
description: |-
  Arc is controlled via a very easy to use command-line interface (CLI).
---

# Arc Commands (CLI)

## Description

Arc is controlled via a very easy to use command-line interface (CLI).
Arc is only a single command-line application: `arc`. This application
then takes a subcommand such as "server" or "run". The complete list of
subcommands is in the navigation to the left.

The Arc CLI is a well-behaved command line application. In erroneous
cases, a non-zero exit status will be returned. It also responds to `-h` and `--help`
as you'd most likely expect.

## Usage

To view a list of the available commands at any time, just run `arc` with
no arguments:

    $ arc
    NAME:
       arc - Remote job execution galore

    USAGE:
       arc [global options] command [command options] [arguments...]

    VERSION:
       0.1.0-dev(72a904e)

    AUTHOR(S):
       Fabian Ruff <fabian.ruff@sap.com> Arturo Reuschenbach Puncernau <a.reuschenbach.puncernau@sap.com>

    COMMANDS:
       server	Run the Arc daemon
       run		Execute an agent action on a remote Arc server
       list		List available agents and actions
       facts	Discover and list facts on this system
       update	Update current binary to the latest version
       init		Initialize server configuration
       status	Service status
       start	Start agent service
       stop		Stop agent service
       restart	Restart agent service
       help, h	Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --config-file, -c "/opt/arc/arc.cfg" Load config file [$ARC_CONFIGFILE]
       --log-level, -l "info"               Log level [$ARC_LOG_LEVEL]
       --help, -h                           show help
       --version, -v                        print the version

To get help for any specific command, pass the `-h` flag to the relevant
subcommand. For example, to see help about the `run` subcommand:

    $ arc run -h
    NAME:
       run - Execute an agent action on a remote Arc server

    USAGE:
       command run [command options] [arguments...]

    OPTIONS:
       --transport, -T "mqtt" Transport backend driver [$ARC_TRANSPORT]
       --endpoint, -e [--endpoint option --endpoint option]	Endpoint url(s) for selected transport [$ARC_ENDPOINT]
       --tls-client-cert      Client cert to use for TLS [$ARC_TLS_CLIENT_CERT]
       --tls-client-key       Private key used in client TLS auth [$ARC_TLS_CLIENT_KEY]
       --tls-ca-cert          CA to verify transport endpoints [$ARC_TLS_CA_CERT]
       --timeout, -t "60"     Timeout for executing the action
       --identity, -i         Target system
       --payload, -p          Payload for action
       --stdin, -s            Read payload from stdin
