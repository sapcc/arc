---
layout: "docs"
page_title: "Arc Api Service"
sidebar_current: "docs-api-basics"
description: Arc
---

# API Service

## Prerequisites

* Postgres database

    Install and configure a Postgres database. We recommend to install Postgres via homebrew:

    ```text
    brew install postgres
    ```

* Database migration with goose

    goose is a database migration tool that allows you to manage your database's evolution by creating incremental SQL or Go scripts.
Please visit [goose](https://bitbucket.org/liamstask/goose) own website for more detail.

  * Install goose

      ```text
      go get bitbucket.org/liamstask/goose/cmd/goose
      ```
  * Create an `arc_dev` and `arc_test` database.

  * Run migration from the api-server folder for all environments needed.

      ```text
      goose -env development up
      ```

      ```text
      goose -env test up
      ```

      You should get similar output like this for each environment:

      ```text
      goose: migrating db environment 'development', current version: 0, target: 20150707152559
      OK    20150624111545_CreateJobs.sql
      OK    20150624112329_CreateLogs.sql
      OK    20150624112740_CreateLogParts.sql
      OK    20150624113227_CreateJsonReplaceFunc.sql
      OK    20150707152559_CreateAgents.sql
      ```


## Running

It is important to set the mandatory parameters (endpoint and db configuration file) when running the API Server.
Here is an example of how to run the api server:

```text
api-server -endpoint tcp://localhost:1883 -db-config api-server/db/dbconf.yml
```

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

The `API server` can be stopped in two ways: gracefully or forcefully. To gracefully halt the server, send the process
an interrupt signal (usually `Ctrl-C` from a terminal or running `kill -INT arc_pid` ).

Alternatively, you can force kill the server by sending it a kill signal. When force killed, the server ends immediately.