---
layout: "docs"
page_title: "Arc API Service - Prerequisites"
sidebar_current: "docs-api-prerequisites"
description: Arc
---

# Prerequisites

## PostgreSQL database

Install and configure a [PostgreSQL](http://www.postgresql.org/) database. We recommend to install PostgreSQL via homebrew:

```text
brew install postgres
```

## Database migration with goose

[goose](https://bitbucket.org/liamstask/goose) is a database migration tool that allows you to manage your database's
evolution by creating incremental SQL or Go scripts.

To setup the database for the `Arc API Service` please follow this steps:

* Install goose:

```text
go get bitbucket.org/liamstask/goose/cmd/goose
```

* Create the databases as defined in the `dbconf.yml`. For development purpose you will just need to create the development
and test database.

* Run migration from the root API Server folder for all environments needed:

```text
goose -env environment up
```

Or use the flag `-path` to indicate where to find the migration files like in the following example:

```text
goose -env development -path api-server/db/ up
```

You should get a similar output for each environment:

```text
goose: migrating db environment 'development', current version: 0, target: 20150707152559
OK    20150624111545_CreateJobs.sql
OK    20150624112329_CreateLogs.sql
OK    20150624112740_CreateLogParts.sql
OK    20150624113227_CreateJsonReplaceFunc.sql
OK    20150707152559_CreateAgents.sql
```

* To roll back a single migration from the current version run following goose command specifying the environment:

```text
goose -env development down
```