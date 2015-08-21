---
layout: "docs"
page_title: "Arc"
sidebar_current: "docs-internals-architecture"
description: |-
  Arc
---


Overview
========
The [INSERT COOL NAME HERE] allows clients to schedule jobs on agents (servers).
The agents process the jobs and report back on the progress and result of the job.

In general the client starts the communication by sending a request to one or more agents.
During the lifecycle of a job the agent sends one more replys to the agent.
At minimum the agent sends one reply when the execution of a job has come to a permanent end.

```
+------------+                     +------------+
|            | +-----Request-----> |            |
|   Client   |                     |   Agent    |
|            |                     |  (Server)  |
|            | <-----Reply-------+ |            |
|            |         .           |            |
|            |         .           |            |
|            |         .           |            |
|            |                     |            |
|            | <-----Reply-------+ |            |
+------------+                     +------------+
```

Protocol
--------
The protocal is based on JSON messages that are passed between the client and agent(s).

The are two types of messages `Request` and `Reply` which share a common structure with some additional fields specific to each type.

### Request
```
{
  "version": int,       // e.g. 1
  "requestid": string,  // e.g. "133B0939-76F4-4C9B-99AB-7A6A873E8C9E",
  "type": "request",    // fixed
  "sender": string      // sender identifier
  "agent": string,      // e.g. "provision"
  "action": string,     // e.g. "execute"
  "payload": string,    // action specfic request payload
  "timeout": int        // timeout in seconds for processing this job
  "to": string          // recipient(s) of the request, e.g. indentity/mo-123456, project/someproject
}
```

### Reply

```
{
  "version": int,       // e.g. 1
  "requestid": string,  // e.g. "133B0939-76F4-4C9B-99AB-7A6A873E8C9E",
  "type": "reply",      // fixed
  "sender": string      // sender identifier
  "agent": string,      // e.g. "provision"
  "action": string,     // e.g. "execute"
  "payload": string,    // action specfic reply payload
  "state": string,      // "queued", "executing", "completed" or "failed"
  "final": bool         // indicates a final message for the job
}
```
