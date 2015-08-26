---
layout: "docs"
page_title: "Arc Api Service"
sidebar_current: "docs-api-basics"
description: |-
  Arc
---

# API Service

The API server offers a RESTful service to Arc, a job scheduler with log collections and supervision based on heartbeat.

## Setup

* Postgresql

The API Server uses a Postgres database.
Install from Brew

* Goose

The Goose
Install Goose with
Run migration

## Routine Scheduler

## Running

Usage: `api-server [global options] command [command options] [arguments...]`

The following global command-line options are available:

* `--transport, -T` - Transport backend driver. If this isn't set, the default transport will be set to MQTT. You can
also have the default value set from the environment via the variable $ARC_TRANSPORT.
* `log-level, -l` - Log level. If this isn't set, the default log level will be set to info.
* `--endpoint, -e [--endpoint option --endpoint option]` -	Endpoint url(s) for selected transport. You can also have
the default value set from the environment via the variable $ARC_ENDPOINT.
* `--tls-client-cert`- Client cert to use for TLS. You can also have the default value set from the environment via
the variable $ARC_TLS_CLIENT_CERT.
* `--tls-client-key` - Private key used in client TLS authentication. You can also have the default value set from
the environment via the variable $ARC_TLS_CLIENT_KEY.
* `--tls-ca-cert` - CA to verify transport endpoints. You can also have the default value set from the environment via
the variable $ARC_TLS_CA_CERT.
* `--bind-address, -b` - Update server URL. If this isn't set, the default bind address will be set to 0.0.0.0:3000.
* `--env` - Environment to use (development, test, production). If this isn't set, the default transport will be set to
development. You can also have the default value set from the environment via the variable $ARC_ENV.
* `--db-config, -c` - Database configuration file.  If this isn't set, the default db config will be set to db/dbconf.yml.
You can also have the default value set from the environment via the variable $ARC_DB_CONFIG.
* `--help, -h` - Show help
* `--version, -v` - Print the version

## Stoping
