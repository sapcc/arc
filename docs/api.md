# API
### Defenition:

| URI                               | GET                                        | PUT                        | POST          | DELETE                    |
|:----------------------------------|:-------------------------------------------|:---------------------------|:--------------|:--------------------------|
| /agent                            | List all agents (with attr: id)            | N/A                        | N/A           | N/A                       |
| /agent/facts                      | List all facts ids                         | N/A                        | N/A           | N/A                       |
| /agent/{agent-id}/facts           | List all facts from an agent with agent-id | N/A                        | N/A           | N/A                       |
| /jobs                             | List all jobs                              | N/A                        | Execute a job | N/A                       |
| /jobs/{job-id}                    | Gets the job with the job-id               | N/A                        | N/A           | N/A                      |
| /jobs/{job-id}/log               | Gets the log from the job with the job-id  | N/A                        | N/A           | N/A                       |


### Agents and facts:
##### Get all agents
- Method: `GET`
- URL: `/agents`
- Response:

```javascript
[
	{
		agent_id: "darwin",
		created_at: "2015-06-18T15:13:10.164307Z",
		updated_at: "2015-06-18T15:13:10.164307Z"
	},
	...
]
```

##### Get an agent
- Method: `GET`
- URL: `/agents/darwin`
- Response:

```javascript
{
	agent_id: "darwin",
	created_at: "2015-06-18T15:13:10.164307Z",
	updated_at: "2015-06-18T15:13:10.164307Z"
}
```

##### Get agent facts:
- Method: `GET`
- URL: `/agents/darwin/facts`
- Response:

```javascript
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

### Jobs

##### Get all job:
- Method: `GET`
- URL: `/jobs`
- Response:

```javascript
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
			payload: "echo "Scritp start" for i in {1..10} do echo $i sleep 1s done echo "Scritp done""
		},
		status: "failed",
		created_at: "2015-06-18T15:23:23.595169Z",
		updated_at: "2015-06-18T15:25:32.09501Z"
	},
	...
]
```

##### Get a job:
- Method: `GET`
- URL: `/jobs/24c744df-1773-429d-b9cf-a0f280e514ca`
- Response:

```javascript
{
	request: {
		version: 1,
		sender: "darwin",
		request_id: "24c744df-1773-429d-b9cf-a0f280e514ca",
		to: "darwin",
		timeout: 60,
		agent: "execute",
		action: "script",
		payload: "echo "Scritp start" for i in {1..10} do echo $i sleep 1s done echo "Scritp done""
	},
	status: "failed",
	created_at: "2015-06-18T15:23:23.595169Z",
	updated_at: "2015-06-18T15:25:32.09501Z"
}
```

##### Get a job log:
- Method: `GET`
- URL: `/jobs/24c744df-1773-429d-b9cf-a0f280e514ca/log`
- Response:

```javascript
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

##### Execute a job:
- Method: `POST`
- URL: `/jobs`
- Body: 

```javascript
{
	"to": "darwin",
	"timeout": 60,
	"agent": "execute",
	"action": "script",
	"payload": "echo \"Scritp start\"\n\nfor i in {1..10}\ndo\n\techo $i\n  sleep 1s\ndone\n\necho \"Scritp done\""
}
```

- Response:

```javascript
{
	request_id: "692f9dd9-f0f4-4332-89e9-ac556a343bd4"
}
```