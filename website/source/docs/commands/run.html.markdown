---
layout: "docs"
page_title: "Commands: Run"
sidebar_current: "docs-commands-run"
description: The `arc run` execute an agent action [REFERENCE ACTIONS DOC PAGE FROM HERE] on a remote Arc server.
---

# Arc Run

Command: `arc run`

The `arc run` execute an agent action [REFERENCE ACTIONS DOC PAGE FROM HERE] on a remote Arc server.

## Usage

Usage: `arc run [command options] [arguments...]`

The following command-line options are available for this command:

* `--transport, -T` - Transport backend driver. If this isn't set, the default transport will be set to MQTT. You can
also have the default value set from the environment via the variable $ARC_TRANSPORT.
* `--endpoint, -e [--endpoint option --endpoint option]` -	Endpoint url(s) for selected transport. You can also have
the default value set from the environment via the variable $ARC_ENDPOINT.
* `--tls-client-cert`- Client cert to use for TLS. You can also have the default value set from the environment via
the variable $ARC_TLS_CLIENT_CERT.
* `--tls-client-key` - Private key used in client TLS authentication. You can also have the default value set from
the environment via the variable $ARC_TLS_CLIENT_KEY.
* `--tls-ca-cert` - CA to verify transport endpoints. You can also have the default value set from the environment via
the variable $ARC_TLS_CA_CERT.
* `--timeout, -t` -  Timeout for executing the action. If this isn't set, the default timeout will be set to 60 seconds.
* `--identity, -i` - Target system. To get the identity of the target system you can run the [`facts` command](/docs/commands/facts.html) on the
remote system.
* `--payload, -p` - Payload for action. [REFERENCE ACTIONS DOC PAGE FROM HERE]
* `--stdin, -s` - Read action payload from stdin. [REFERENCE ACTIONS DOC PAGE FROM HERE]

## Examples

Example: `arc run -endpoint tcp://localhost:1883 -identity darwin rpc version`

It prints the current version of the Arc running on the remote server with the identity `darwin`.

```text
0.1.0-dev(ae07667)
```

Example: `arc run -endpoint tcp://localhost:1883 -identity darwin -payload "echo Script start; for i in {1..5}; do echo \$i; sleep 1s; done; echo Script done" execute script`

```text
Script start
1
2
3
4
5
Script done
```

It runs the script given as a payload on the remote server with identity `darwin`. If the script gets to complicated or too long to give it as a payload, there is possibility to
write a bash script file and give it to the `run` command as the following example:

Example: `arc run -endpoint tcp://localhost:1883 -identity darwin -stdin execute script < script.sh`

```text
#!/bin/bash
# script.sh

echo Scritp start
for i in {1..5}; do
	echo $i
	sleep 1s
done
echo Scritp done
```
