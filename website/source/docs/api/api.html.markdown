---
layout: "docs"
page_title: "Arc API Service - HTTP API"
sidebar_current: "docs-api-api"
description: The main interface to Arc is a RESTful HTTP API. The API can be used to perform operations or collect information from one or different Arc servers.
---

<a name="back_to_top"></a>
# HTTP API

The main interface to Arc is a RESTful HTTP API. The API can be used to perform operations or collect
information from one or different Arc servers.

-> **Note:** The term `agent` will be changed to `node` in the near future. Allthough it is not yet changed in the Arc API and this documentation, it is already being used in the <b>[Elektra dashboard](https://dashboard.***REMOVED***)</b>, <b>[Lyra-CLI](https://documentation.***REMOVED***/docs/automation/cli/)</b> and in the <b>[global documentation](https://documentation.***REMOVED***/docs/automation/details.html)</b>.


* [Definition](#definition)
* [Paginating lists](#paginating_lists)
* [List all agents](#list_all_agents)
  * [Filtering agents](#filter_agents)
  * [Showing specific agent facts](#show_facts_agents)
* [Bootstrap agent](#init_agent)
* [Get an agent](#get_agent)
* [Delete an agent](#delete_agent)
* [Show agent facts](#show_agent_facts)
* [Show agent tags](#show_agent_tags)
* [Add an agent tags](#add_agent_tag)
* [Delete an agent tag](#delete_agent_tag)
* [List all jobs](#list_all_jobs)
  * [Filtering jobs](#filter_jobs)
* [Get a job](#get_job)
* [Get a job log](#get_job_log)
* [Execute a job](#execute_job)

<a name="definition"></a>
## Definition

| URL                               | GET                    | PUT                        | POST          | DELETE                    |
|:----------------------------------|:-----------------------|:---------------------------|:--------------|:--------------------------|
| /agents                           | List all agents        | N/A                        | N/A           | N/A                       |
| /agents/init                      | N/A                    | N/A                        | Bootstrap agent | N/A               |
| /agents/{agent-id}                | Get an agent           | N/A                        | N/A           | Delete an agent           |
| /agents/{agent-id}/facts          | Show agent facts       | N/A                        | N/A           | N/A                       |
| /agents/{agent-id}/tags           | Show agent tags        | N/A                        | Add a tag     | Delete a tag              |
| /jobs                             | List all jobs          | N/A                        | Execute a job | N/A                       |
| /jobs/{job-id}                    | Get a job              | N/A                        | N/A           | N/A                       |
| /jobs/{job-id}/log                | Get a job log          | N/A                        | N/A           | N/A                       |
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="paginating_lists"></a>
## Paginating lists
All lists served by the Arc API are by default paginated. The pagination can be used adding following parameters to the request url:

| Parameter                 |Default value    | Description                                           |
|:--------------------------|:----------------|:------------------------------------------------------|
| page                      |1                | Subset of the whole set available                     |
| per_page                  |25               | Limit the number of resources with a maximum of `100` |


In the response header extra parameters will be added to handle the pagination.

| Headers                   | Description                                                 |
|:--------------------------|:------------------------------------------------------------|
| Pagination-Elements       | Returns total elements available from the set               |
| Pagination-Pages          | Returns total number of pages using same per_page parameter |
| Link                      | Returns links to paginate the data                          |

- Example URL:

`https://arc-staging.***REMOVED***/api/v1/jobs?page=2&per_page=50`

- Example response headers:

```text
Content-Type: application/json; charset=UTF-8
Link: </api/v1/jobs?page=2&per_page=50>;rel="self",</api/v1/jobs?page=1&per_page=50>;rel="first",</api/v1/jobs?page=1&per_page=50>;rel="prev",</api/v1/jobs?page=3&per_page=50>;rel="next",</api/v1/jobs?page=3&per_page=50>;rel="last"
Pagination-Elements: 125
Pagination-Pages: 3
X-Served-By: api-c7287855605592c8309e5d58174e6e30-g6vj5
Date: Tue, 26 Jan 2016 13:04:28 GMT
Transfer-Encoding: chunked
```
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="list_all_agents"></a>
## List all agents
- Method: `GET`
- URL: `/agents`
- Example response:

```json
[
	{
		agent_id: "darwin",
		display_name: "test_server",
		created_at: "2015-06-18T15:13:10.164307Z",
		updated_at: "2015-06-18T15:13:10.164307Z"
	},
	...
]
```

Agents are `sorted` by the attribute `display_name`. The `display_name` attribute is a virtual attribute created on time of retrieving the agent information as following:

- check tag with key `name`, if not given
- check fact with key `hostname`, if not given
- set the `agent_id` attribute

<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="filter_agents"></a>
### Filtering agents

We use a self written parser that transforms the filter syntax exposed by the API to a filter expression that can by used by the underlying fact storage system.
Following operators are available:

<div class="filter-operators">

| Comparison Operators      | Description                                       |
|:--------------------------|:--------------------------------------------------|
| =                         | Performs a equal-to comparison                    |
| !=                        | Performs a not-equal-to comparison                |
| ^                         | Performs a comparison with wildcards, where `*` matches zero or more characters and `+` matches exactly one character. (e.g `"hallo" ^ "*ll+"` would be a match) |
| !^                        | The negation of the ^ comparison (does not match) |

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

```json
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
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="show_facts_agents"></a>
### Showing specific agents facts

Adding the `facts` parameter in the [List all agents](#list_all_agents) call facts can be selected dinamically to be shown in the JSON response.
A list of all available facts can be found [here](/docs/server/facts.html).

- Method: `GET`
- URL: `/agents?facts={coma separated facts}`
- Example URL: `/agents?facts=os,online`
- Example response:

```json
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
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="init_agent"></a>
## Bootstrap agent

Create a one-time token for joining an arc agent.
This endpoint supports different output formats that can be selected by providing an appriopiate `Accept` header in the request

- Method: `POST`
- URL: `/agents/init`
- Supported Content-Types: `application/json` (default), `text/cloud-config`, `text/x-shellscript`, `text/x-powershellscript`

Example response for `Accept: application/json`

```json
{
  "token": "4d523051-089f-41ce-aaf7-727fee19c28a",
  "url": "https://arc.example.com/api/v1/agents/init/4d523051-089f-41ce-aaf7-727fee19c28a",
  "endpoint_url": "tls://arc-broker.example.com:8883",
  "update_url": "https://stable.arc.example.com"
}
```

Example  response fpr `Accept: text/cloud-config`

```yaml
#cloud-config
runcmd:
  - - sh
    - -ec
    - |
      curl -f --create-dirs -o /opt/arc/arc https://stable.arc.example.com/arc/linux/amd64/latest
      chmod +x /opt/arc/arc
      /opt/arc/arc init --endpoint tls://arc-broker.example.com:8883 --update-uri https://stable.arc.example.com --registration-url https://arc.example.com/api/v1/agents/init/506a4692-84be-41cc-b5e5-9e4d4184f6cd
```

<a name="get_agent"></a>
## Get an agent
- Method: `GET`
- URL: `/agents/{agent-id}`
- Example URL: `/agents/darwin`
- Example response:

```json
{
	agent_id: "darwin",
	created_at: "2015-06-18T15:13:10.164307Z",
	updated_at: "2015-06-18T15:13:10.164307Z"
}
```
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="delete_agent"></a>
## Delete an agent
- Method: `DELETE`
- URL: `/agents/{agent-id}`
- No response body is provided if the request succeeds
- Example URL: `/agents/darwin`

<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="show_agent_facts"></a>
## Show agent facts
- Method: `GET`
- URL: `/agents/{agent-id}/facts`
- Example URL: `/agents/darwin/facts`
- Example response:

```json
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
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="show_agent_tags"></a>
## Show agent tags
- Method: `GET`
- URL: `/agents/{agent-id}/tags`
- Example URL: `/agents/d84ca366-c963-454f-9bd7-854121a0117e/tags`
- Example response:

```json
{
	pool: "green",
	landscape: "staging",
}
```
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="add_agent_tag"></a>
## Add an agent tags
- Method: `POST`
- URL: `/agents/{agent-id}/tags`
- Content-Types: `application/json`
- No response body is provided if the request succeeds
- Example request:
	- URL: `/agents/d84ca366-c963-454f-9bd7-854121a0117e/tags`  
	- Body:

  ```text
  {"pool":"green","landscape":"staging"}
  ```

All tag keys musst be `alphanumeric [a-z0-9A-Z]` and have non empty values. In case of an error
body will contain the error messages as JSON.

In case of adding an existing tag the value will be replaced with the new submitted.

<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="delete_agent_tag"></a>
## Delete an agent tag
- Method: `DELETE`
- URL: `/agents/{agent-id}/tags/{tag-key}`
- No response body is provided if the request succeeds
- Example URL: `/agents/d84ca366-c963-454f-9bd7-854121a0117e/tags/pool`

<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="list_all_jobs"></a>
## List all jobs
- Method: `GET`
- URL: `/jobs`
- Example response:

```json
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

Jobs are `sorted` by creation date newest first.

<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="filter_jobs"></a>
### Filtering jobs
Adding the `agent_id` parameter in the [List all jobs](#list_all_jobs) call jobs can be filtered by agent.

- Example URL: `/jobs?agent_id=d84ca366-c963-454f-9bd7-854121a0117e`
- Example response:

```json
[
	{
		request: {
			version: 1,
			sender: "darwin",
			request_id: "24c744df-1773-429d-b9cf-a0f280e514ca",
			to: "d84ca366-c963-454f-9bd7-854121a0117e",
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
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="get_job"></a>
## Get a job
- Method: `GET`
- URL: `/jobs/{job-id}`
- Example URL: `/jobs/24c744df-1773-429d-b9cf-a0f280e514ca`
- Example response:

```json
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
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

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
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>

<a name="execute_job"></a>
## Execute a job
- Method: `POST`
- URL: `/jobs`
- Content-Types: `application/json`
- Example request body:

```json
{
	"to": "darwin",
	"timeout": 60,
	"agent": "execute",
	"action": "script",
	"payload": "echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""
}
```

- Example response:

```json
{
	request_id: "692f9dd9-f0f4-4332-89e9-ac556a343bd4"
}
```
<a href="#back_to_top" class="back_to_top">Top &uarr;</a>
