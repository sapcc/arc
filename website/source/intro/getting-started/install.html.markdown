---
layout: "intro"
page_title: "Installing Arc"
sidebar_current: "gettingstarted-install"
description: |-
 TBD
---

# Install Arc

The central component of Arc is the **Arc Server** which needs to be
 installed on every system that should be managed by Arc.
To make installation easy, Arc is distributed as a
[binary package](/downloads.html) for all supported platforms and
architectures. This page will not cover how to compile Arc from
source.

## Installing Arc
Just download the appropriate binary for your platform and place the
 binary in a directory of your choosing. 
It is recommended to put it in a dedicated directory (e.g. `/opt/arc` or `C:\arc`).

Though not strictly necessary you might want to add the installation directory to
 your `PATH` environment variable for easy access to the `arc` command line interface.

## Verifying the Installation
After installing Arc, verify the installation worked by opening a new terminal session
 and run the `arc` (or `arc.exe`) command without any argument. You should see output similar to this:

```
$ /opt/arc/arc
NAME:
   arc - Remote job execution galore

USAGE:
   arc [global options] command [command options] [arguments...]

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
   --config-file, -c "/opt/arc/arc.cfg"	Load config file [$ARC_CONFIGFILE]
   --log-level, -l "info"		Log level [$ARC_LOG_LEVEL]
   --help, -h				show help
   --version, -v			print the version
```


## Next Steps

Arc is installed but before its ready for operation. We need to 
[run the Messaging Middleware](middleware.html)!
