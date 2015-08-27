---
layout: "docs"
page_title: "Arc Api Service - Index"
sidebar_current: "docs-api"
description: The API server, the main interface to Arc, offers a RESTful API service, a job scheduler with log collections and supervision based on heartbeat.
---

# Arc API Service

The API server, the main interface to Arc, offers a RESTful API service, a job scheduler with log collections and supervision based on heartbeat.

## HTTP API

Due to the variety of options, the API is documented in its own section.
See the [HTTP API](/docs/api/api.html) section for more information.

## Job scheduler

Each job executed through the API Server will be persisted in the database. Due to the fact that the server subscribes to all
`message replies` allows the system to track all jobs lifecycle and save the corresponding log chunks.

A job lifecycle comprehends following states:

| State             | Description                                                                 |
|:------------------|-----------------------------------------------------------------------------|
| Queued            | Default status when executing a job                                         |
| Executing         | This status is being set when the target system receives a job              |
| Failed            | Target system replies with a failed status when the job exits with an error |
| Complete          | Target system replies with a complete status when the job has been finished |

Please visit following sections to know more about how to execute jobs, retrieve jobs or logs data.

* [Execute a job](/docs/api/api.html#execute_job)
* [Retrieve job data](/docs/api/api.html#get_job)
* [Retrieve job logs data](/docs/api/api.html#get_job_log)

## Supervision

The Arc server subscribes also to all `registry replies` allowing to persist all facts from all remote Arc agents.

Each remote Arc agent once is connected sends all available facts. Here is an example of all available facts:

```text
{
	os: "darwin",
	online: true,
	project: "test-project",
	hostname: "BERM32186999A",
	identity: "darwin",
	platform: "darwin",
	arc_version: "0.1.0-dev(HEAD)",
	memory_used: 13662908416,
	memory_total: 17179869184,
	organization: "test-org",
	platform_family: "",
	memory_available: 3516960768,
	platform_version: "14.3.0",
	memory_used_percent: 80
}
```

Each time a fact changes the remote agent sends a new `registry reply` just with the facts that have updated. Here is an example
of an update:

```text
{
	memory_used: 17163505664,
	memory_total: 17179869184,
	memory_available: 16363520,
	memory_used_percent: 100
}
```

Please visit following sections to know more about how to retrieve agents and agents facts.

* [Retrieve agents](/docs/api/api.html#list_all_agents)
* [Retrieve agents filtered](/docs/api/api.html#filter_agents)
* [Retrieve agent facts](/docs/api/api.html#list_agent_facts)