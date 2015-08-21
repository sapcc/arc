---
layout: "docs"
page_title: "Agent"
sidebar_current: "docs-server-basics"
description: |-
 TBD
---

# Arc Server 

## Running a Server

## Stopping a Server 

An agent can be stopped in two ways: gracefully or forcefully. To gracefully
halt an agent, send the process an interrupt signal (usually
`Ctrl-C` from a terminal or running `kill -INT consul_pid` ). When gracefully exiting, the agent first notifies
the cluster it intends to leave the cluster. This way, other cluster members
notify the cluster that the node has _left_.

Alternatively, you can force kill the agent by sending it a kill signal.
When force killed, the agent ends immediately. The rest of the cluster will
eventually (usually within seconds) detect that the node has died and
notify the cluster that the node has _failed_.

It is especially important that a server node be allowed to leave gracefully
so that there will be a minimal impact on availability as the server leaves
the consensus quorum.

For client agents, the difference between a node _failing_ and a node _leaving_
may not be important for your use case. For example, for a web server and load
balancer setup, both result in the same outcome: the web node is removed
from the load balancer pool.
