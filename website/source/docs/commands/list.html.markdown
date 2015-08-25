---
layout: "docs"
page_title: "Commands: List"
sidebar_current: "docs-commands-list"
description: The `list` command shows all available agents and their actions.
---

# Arc List

Command: `arc list`

The `list` command shows all available agents and their actions. The payload requirements is action dependent
and it should be looked up in the corresponding implementation. Some payload examples can be found in the
[run command](/docs/commands/run.html) documentation. A complete explanation over the available agents can be found
in the [arc server](/docs/server/agents.html) documentation topic.

## Example

Example: `arc list`

```text
  Agent               Actions
----------------------------------------
  rpc                 enable,disable,ping,sleep,version
  execute             enable,disable,command,script
  chef                enable,disable,zero
```