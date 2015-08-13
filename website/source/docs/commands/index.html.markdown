---
layout: "docs"
page_title: "Commands"
sidebar_current: "docs-commands"
description: |-
  ARC is controlled via a very easy to use command-line interface (CLI).
---

# ARC Commands (CLI)

ARC is controlled via a very easy to use command-line interface (CLI).
ARC is only a single command-line application: `arc`. This application
then takes a subcommand such as "server" or "execute". The complete list of
subcommands is in the navigation to the left.

The `ARC` CLI is a well-behaved command line application. In erroneous
cases, a non-zero exit status will be returned. It also responds to `-h` and `--help`
as you'd most likely expect.

To view a list of the available commands at any time, just run `arc` with
no arguments:

```text
$ arc
NAME:
    arc - Remote job execution galore

USAGE:
    arc [global options] command [command options] [arguments...]

VERSION:
    0.1.0-dev(90de7b9)

AUTHOR(S):
    Fabian Ruff <fabian.ruff@sap.com> Arturo Reuschenbach Puncernau <a.reuschenbach.puncernau@sap.com>

COMMANDS:
    server     Run the agent
    execute    Execute remote agent action
    list       List available agents and actions
    facts      Discover and list facts on this system
    update     Check for new updates and update to the last version
    init       Initialize server configuration
    status     Service status
    help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
    --config-file, -c "/opt/arc/arc.cfg"
        Load config file [$ARC_CONFIGFILE]
    --transport, -t "mqtt"
        Transport backend driver [$ARC_TRANSPORT]
    --endpoint, -e [--endpoint option --endpoint option]
        Endpoint url(s) for selected transport [$ARC_ENDPOINT]
    --tls-ca-cert
        CA to verify transport endpoints [$ARC_TLS_CA_CERT]
    --tls-client-cert
        Client cert to use for TLS [$ARC_TLS_CLIENT_CERT]
    --tls-client-key
        Private key used in client TLS auth [$ARC_TLS_CLIENT_KEY]
    --log-level, -l "info"
        Log level [$ARC_LOG_LEVEL]
    --no-auto-update
        Should NO trigger auto updates [$ARC_NO_AUTO_UPDATE]
    --update-interval "21600"
        Time update interval in seconds [$ARC_UPDATE_INTERVAL]
    --update-uri "http://localhost:3000/updates"
        Update server uri [$ARC_UPDATE_URI]
    --help, -h
        Show help
    --version, -v
        Print the version

```

To get help for any specific command, pass the `-h` flag to the relevant
subcommand. For example, to see help about the `execute` subcommand:

```text
$ arc execute -h
NAME:
    execute - Execute remote agent action

USAGE:
    command execute [command options] [arguments...]

OPTIONS:
    --timeout, -t "60"    Timeout for executing the action
    --identity, -i        Target system
    --payload, -p         Payload for action
    --stdin, -s           Read payload from stdin
```
