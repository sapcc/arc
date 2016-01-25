---
layout: "docs"
page_title: "Arc API Service - HTTP API"
sidebar_current: "docs-api-api"
description: The main interface to Arc is a RESTful HTTP API. The API can be used to perform operations or collect information from one or different Arc servers.
---

# HTTP API

The main interface to Arc is a RESTful HTTP API. The API can be used to perform operations or collect
information from one or different Arc servers.

* [Definition](#definition)
* [List all agents](#list_all_agents)
  * [Filtering agents](#filter_agents)
  * [Showing specific agent facts](#show_facts_agents)	
* [Get an agent](#get_agent)
* [Delete an agent](#delete_agent)
* [List agent facts](#list_agent_facts)
* [List agent tags](#list_agent_tags)
* [Add an agent tag](#add_agent_tag)
* [Delete an agent tag](#delete_agent_tag)
* [List all jobs](#list_all_jobs)
* [Get a job](#get_job)
* [Get a job log](#get_job_log)
* [Execute a job](#execute_job)

<a name="definition"></a>
## Definition

| URL                               | GET                    | PUT                        | POST          | DELETE                    |
|:----------------------------------|:-----------------------|:---------------------------|:--------------|:--------------------------|
| /agents                           | List all agents        | N/A                        | N/A           | N/A                       |
| /agents/{agent-id}                | Get an agent           | N/A                        | N/A           | Delete an agent           |
| /agents/{agent-id}/facts          | List agent facts       | N/A                        | N/A           | N/A                       |
| /agents/{agent-id}/tags           | List agent tags        | N/A                        | Add a tag     | Delete a tag              |
| /jobs                             | List all jobs          | N/A                        | Execute a job | N/A                       |
| /jobs/{job-id}                    | Get a job              | N/A                        | N/A           | N/A                       |
| /jobs/{job-id}/log                | Get a job log          | N/A                        | N/A           | N/A                       |

<a name="list_all_agents"></a>
## List all agents
- Method: `GET`
- URL: `/agents`
- Example response:

```text
[
	{
		agent_id: "darwin",
		created_at: "2015-06-18T15:13:10.164307Z",
		updated_at: "2015-06-18T15:13:10.164307Z"
	},
	...
]
```

<a name="filter_agents"></a>
### Filtering agents
We use a self written parser that transforms the filter syntax exposed by the API to a filter expression that can by used by the underlying fact storage system.
Following operators are available:

<div class="filter-operators">

| Comparison Operators      | Description                                |
|:--------------------------|:-------------------------------------------|
| =                         | Performs a equal-to comparison             |
| !=                        | Performs a not-equal-to comparison         |

</div>

<div class="filter-operators">

| Logical Operators         | Description                                |
|:--------------------------|:-------------------------------------------|
| OR                        | Performs a logical OR on the condition     |
| AND                       | Performs a logical AND on the condition    |
| NOT                       | Performs a logical NOT on the condition    |
| ()                        | When parenthesized conditions are nested, the innermost condition is evaluated first |

</div>

<div class="filter-operators">

| Attribute types           | Description                                |
|:--------------------------|:-------------------------------------------|
| @attribute_name           | Will match a Fact defined by the attribute name |
| attribute_name            | Will match a Tag defined by the attribute name    |

</div>

- Method: `GET`
- URL: `/agents?q={filter}`
- Example URL: `/agents?q=@os+%3D+%22darwin%22+OR+%28landscape+%3D+%22staging%22+AND+pool+%3D+%22green%22%29`
- Example filter not encoded: `@os = "darwin" OR (landscape = "staging" AND pool = "green")`
- Example response:

```text
[
	{
		agent_id: "darwin",
		project: "test-project",
		organization: "test-org",
		created_at: "2015-07-08T14:17:47.184796Z",
		updated_at: "2015-08-26T10:07:06.310511Z"
	},
	...
]
```

<a name="show_facts_agents"></a>
### Showing specific agents facts

Adding the `facts` parameter in the [List all agents](#list_all_agents) call facts can be selected dinamically to be shown in the JSON response.
A list of all available facts can be found [here](/docs/server/facts.html).

- Method: `GET`
- URL: `/agents?facts={coma separated facts}`
- Example URL: `/agents?facts=os,online`
- Example response:

```text
[
	{
		agent_id: "mo-d90b4b6fe",
		project: "p-ea1868652",
		organization: "o-monsoon2",
		facts: {
			online: false,
			os: "linux"
		},
		created_at: "2015-10-22T11:02:56.359063Z",
		updated_at: "2015-11-11T13:13:37.955197Z",
		updated_with: "6f79af48-00ff-4dbb-b57e-1061e5d3c635",
		updated_by: "api-6994d829bae7811cca179eb72c3cf634-kiebm"
		}
	...
]
```

Use the special key `all` as a value to the `facts` parameter to show all available facts. If the key is being added in combination with other facts this 
will be ignored.

<a name="get_agent"></a>
## Get an agent
- Method: `GET`
- URL: `/agents/{agent-id}`
- Example URL: `/agents/darwin`
- Example response:

```text
{
	agent_id: "darwin",
	created_at: "2015-06-18T15:13:10.164307Z",
	updated_at: "2015-06-18T15:13:10.164307Z"
}
```

<a name="delete_agent"></a>
## Delete an agent
- Method: `DELETE`
- URL: `/agents/{agent-id}`
- Example URL: `/agents/darwin`
- Example response:

```text
Agent with id "darwin" deleted. 
```

<a name="list_agent_facts"></a>
## List agent facts
- Method: `GET`
- URL: `/agents/{agent-id}/facts`
- Example: `/agents/darwin/facts`
- Example response:

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

<a name="list_agent_tags"></a>
## List agent tags
- Method: `GET`
- URL: `/agents/{agent-id}/tacts`
- Example: `/agents/d84ca366-c963-454f-9bd7-854121a0117e/tags`
- Example response:

```text
{
	pool: "green",
	landscape: "staging",
}
```

<a name="add_agent_tag"></a>
## Add an agent tag
- Method: `POST`
- URL: `/agents/{agent-id}/tags`
- Example: 
	- URL: `/agents/d84ca366-c963-454f-9bd7-854121a0117e/tags`
	- Body: `pool=green&lanscape=staging`
	- Headers: `Content-Type: application/x-www-form-urlencoded`
- Example response:

```text
All tags saved!
```

<a name="delete_agent_tag"></a>
## Delete an agent tag
- Method: `DELETE`
- URL: `/agents/{agent-id}/tags/{tag-key}`
- Example: `/agents/d84ca366-c963-454f-9bd7-854121a0117e/tags/pool`
- Example response:

```text
Tag from agent with id "d84ca366-c963-454f-9bd7-854121a0117e" and value "pool" is removed!
```

<a name="list_all_jobs"></a>
## List all jobs
- Method: `GET`
- URL: `/jobs`
- Example response:

```text
[
	{
		request: {
			version: 1,
			sender: "darwin",
			request_id: "24c744df-1773-429d-b9cf-a0f280e514ca",
			to: "darwin",
			timeout: 60,
			agent: "execute",
			action: "script",
			payload: "echo "Script start" for i in {1..10} do echo $i sleep 1s done echo "Script done""
		},
		status: "failed",
		created_at: "2015-06-18T15:23:23.595169Z",
		updated_at: "2015-06-18T15:25:32.09501Z"
	},
	...
]
```

<a name="get_job"></a>
## Get a job
- Method: `GET`
- URL: `/jobs/{job-id}`
- Example URL: `/jobs/24c744df-1773-429d-b9cf-a0f280e514ca`
- Example response:

```text
{
	request: {
		version: 1,
		sender: "darwin",
		request_id: "24c744df-1773-429d-b9cf-a0f280e514ca",
		to: "darwin",
		timeout: 60,
		agent: "execute",
		action: "script",
		payload: "echo "Script start" for i in {1..10} do echo $i sleep 1s done echo "Script done""
	},
	status: "completed",
	created_at: "2015-06-18T15:23:23.595169Z",
	updated_at: "2015-06-18T15:25:32.09501Z"
}
```

<a name="get_job_log"></a>
## Get a job log
- Method: `GET`
- URL: `/jobs/{job-id}/log`
- Example URL: `/jobs/24c744df-1773-429d-b9cf-a0f280e514ca/log`
- Example response:

```text
Reading package lists... Done
Building dependency tree
Reading state information... Done
The following extra packages will be installed:
  unzip
The following NEW packages will be installed:
  unzip zip
0 upgraded, 2 newly installed, 0 to remove and 0 not upgraded.
Need to get 455 kB of archives.
After this operation, 993 kB of additional disk space will be used.
Do you want to continue? [Y/n] Y
Get:1 http://archive.ubuntu.com/ubuntu/ trusty/main unzip amd64
6.0-9ubuntu1 [193 kB]
Get:2 http://archive.ubuntu.com/ubuntu/ trusty/main zip amd64 3.0-8 [262
kB]
Fetched 455 kB in 0s (456 kB/s)
Selecting previously unselected package unzip.
(Reading database ... 11813 files and directories currently installed.)
Preparing to unpack .../unzip_6.0-9ubuntu1_amd64.deb ...
Unpacking unzip (6.0-9ubuntu1) ...
Selecting previously unselected package zip.
Preparing to unpack .../archives/zip_3.0-8_amd64.deb ...
Unpacking zip (3.0-8) ...
Processing triggers for mime-support (3.54ubuntu1.1) ...
Setting up unzip (6.0-9ubuntu1) ...
Setting up zip (3.0-8) ...
```

<a name="execute_job"></a>
## Execute a job
- Method: `POST`
- URL: `/jobs`
- Example body:

```text
{
	"to": "darwin",
	"timeout": 60,
	"agent": "execute",
	"action": "script",
	"payload": "echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""
}
```

- Example response:

```text
{
	request_id: "692f9dd9-f0f4-4332-89e9-ac556a343bd4"
}
```
