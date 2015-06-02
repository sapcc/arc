API
========

| URI                              | GET																																				| PUT   | POST   										| DELETE										|
| --------------------------------- |:-------------------------------------------------------------------------:|:-----:|:-------------------------:|--------------------------:|
| /agent														| list all agents (with attr: id)																						| N/A   | N/A												| N/A												|
| /agent/facts											| list all facts ids															 													| N/A   | N/A												| N/A												|
| /agent/{agent-id}/facts						| list all facts from an agent with agent-id			 													| N/A   | N/A												| N/A												|
| /agent/{agent-id}/facts/{fact-id}	| Gets the fact with the unique agent-id and fact-id												| N/A   | N/A												| N/A												|
| /job															| list all jobs 																														| N/A   | N/A												| N/A												|
| /job/{job-id}											| list all agents (attr: status, payload) where the job will be triggered		| N/A   | N/A												| killing/canceling the job	|
| /job/{job-id}/{agent-id}/log			| Gets the log from the job with the job-id	and agent-id										| N/A   | N/A												| N/A												|


Facts
--------
Retrieve status and settings from one agent:

| Name								| Key	 					| Description														|
| --------------------|:-------------:|--------------------------------------:|
| Status							| status				|																				|
| Memory							|	memory				|																				|
| Plattform						|	os						|																				|
| Architecture				|	arch					|																				|