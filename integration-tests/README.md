Arc integration tests
=============================

- Environment variable ARC_API_SERVER overrides -api-server parameter
- Environment variable ARC_UPDATE_SERVER -update-server parameter

There are two ways to provide authentication:

- Giving swift parameters:
  - keystone-endpoint
  - username
  - password
  - project
  - domain
- Giving a valid token:
  - token

smokes
-----------------
Environment variable LATEST_VERSION override -latest-version parameter

```text
bin/smoke -api-server https://arc-api-server -update-server http://arc-update-server -token XXX -test.v
```

Job (works on linux and windows)
-----------------
Environment variable AGENT_IDENTITY overrides -agent-identity parameter

```text
bin/job-test -api-server https://arc-api-server -arc-agent mo-2a9b97c0c -test.v
```

updated and online (check all connected agents)
-----------------
Environment variable LATEST_VERSION override -latest-version parameter

```text
bin/updated-online-test -api-server https://arc-api-server -update-server http://arc-update-server -arc-last-deployed-version 20150916.5 -test.v
```

Facts
-----------------
Environment variable AGENT_IDENTITY overrides -agent-identity parameter

```text
bin/facts-test -api-server https://arc-api-server -arc-agent mo-dfcb6e3fd -test.v
````
