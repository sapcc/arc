---
layout: "docs"
page_title: "Commands: Server"
sidebar_current: "docs-commands-server"
description: Run the Arc daemon.
---

# Arc Server

Command: `arc server`

## Description

Run the Arc daemon.

Due to the power and complexity of this command, the Arc server is documented in its own section.
See the [Arc Server](/docs/server/basics.html) section for more information on how to use this command and the options it has.

## Usage

Usage: `arc server command [command options] [arguments...]`

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
* `--no-auto-update` - Specifies if the server should NO trigger auto updates. You can also have the default value
set from the environment via the variable $ARC_NO_AUTO_UPDATE.
* `--update-uri` - Update server uri. If this isn't set, the default transport will be set to http<nolink>://localhost:3000/updates.
You can also have the default value set from the environment via the variable $ARC_UPDATE_URI.
* `--update-interval` - Time update interval in seconds. If this isn't set, the default transport will be set to 21600 seconds.
You can also have the default value set from the environment via the variable $ARC_UPDATE_INTERVAL.