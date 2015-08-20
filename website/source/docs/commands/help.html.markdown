---
layout: "docs"
page_title: "Commands: Help"
sidebar_current: "docs-commands-help"
description: The `help` shows a list of commands or help for one command.
---

# Arc Help

Command: `arc help`

The `help` command shows a list of commands or help for one specific command.

## Example

Example: `arc execute -help`

It shows the help for the command execute.

```text
NAME:
    execute - Execute remote agent action

USAGE:
    command execute [command options] [arguments...]

OPTIONS:
    --timeout, -t "60"  Timeout for executing the action
    --identity, -i      Target system
    --payload, -p       Payload for action
    --stdin, -s         Read payload from stdin
