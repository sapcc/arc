# API
### Defenition:

| URI                               | GET                                        | PUT                        | POST          | DELETE                    |
|:----------------------------------|:-------------------------------------------|:---------------------------|:--------------|:--------------------------|
| /agent                            | List all agents (with attr: id)            | N/A                        | N/A           | N/A                       |
| /agent/facts                      | List all facts ids                         | N/A                        | N/A           | N/A                       |
| /agent/{agent-id}/facts           | List all facts from an agent with agent-id | N/A                        | N/A           | N/A                       |
| /agent/{agent-id}/facts/{fact-id} | Gets the fact with agent-id and fact-id    | N/A                        | N/A           | N/A                       |
| /jobs                             | List all jobs                              | N/A                        | Execute a job | N/A                       |
| /jobs/{job-id}                    | Gets the job with the job-id               | Update the job with job-id | N/A           | killing/canceling the job |
| /jobs/{job-id}/log                | Gets the log from the job with the job-id  | N/A                        | N/A           | N/A                       |


### Examples:
##### Execute a job:
- Method: `POST`
- URL: `http://localhost:3000/jobs`
- Body: `{"sender":"me","to":"you","timeout":1,"agent":"007","action":"hhmm","payload":"payload"}`

# Facts
Retrieve status and settings from one agent:

| Name                | Key          | Description          |
|:--------------------|:-------------|:---------------------|
| Status              | status       | ...                  |
| Memory              |	 memory       | ...                  |
| Plattform           | os           | ...                  |
| Architecture        | arch         | ...                  |