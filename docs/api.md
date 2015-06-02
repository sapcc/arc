| URI                              | GET																																				| PUT   | POST   										| DELETE										|
| --------------------------------- |:-------------------------------------------------------------------------:|:-----:|:------------------------:|---------------------------:|
| /agent														| list all agents (with attr: status/availability)													| N/A   | N/A												| N/A												|
| /agent/facts											| list all facts ids															 													| N/A   | N/A												| N/A												|
| /agent/{agent-id}/facts						| list all facts																	 													| N/A   | N/A												| N/A												|
| /agent/{agent-id}/facts/{fact-id}	| Gets the fact with the unique agent-id and fact-id												| N/A   | N/A												| N/A												|

| /agent/{agent-id}/job							| list all jobs from the given agent-id																			| N/A   | N/A												| N/A												|
| /agent/{agent-id}/job/{job-id}		| Gets the job with the unique agent-id and job-id													| N/A   | N/A												| killing/canceling the job	|

| /job															| list all jobs (with attr: status)																					| N/A   | N/A												| N/A												|
| /job/{job-id}											| list all agents (attr: status, payload) where the job will be triggered		| N/A   | N/A												| N/A												|
| /job/{job-id}/{agent-id}/log			| Gets the log from the job with the job-id	and agent-id										| N/A   | N/A												| N/A												|