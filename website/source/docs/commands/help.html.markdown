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

Example: `arc update -help`

It shows the help for the [`update` command](/docs/commands/update.html).

```text
NAME:
   update - Update current binary to the latest version

USAGE:
   command update [command options] [arguments...]

OPTIONS:
   --force, -f                                  No confirmation is needed
   --update-uri "http://localhost:3000/updates" Update server uri [$ARC_UPDATE_URI]
   --no-update, -n                              No update is triggered
```