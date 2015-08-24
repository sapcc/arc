---
layout: "intro"
page_title: "Introduction"
sidebar_current: "what"
description: |-
 Arc is a remote job execution framework. Its main focus is orchestrating and administering large clusters of servers in a cloud environment.
---

# Introduction to Arc

Welcome to the introduction to Arc! This guide is the best place to start
with Arc. We cover what Arc is, what problems it can solve, how it compares
to existing software, and how you can get started using it. If you are familiar
with the basics of Arc, the [documentation](/docs/index.html) provides a more
detailed reference of available features.

## What is Arc?
Arc is a remote job execution framework. Its main focus is orchestrating and administering large clusters of servers in a cloud environment.
It uses an *agent*-based approach, as such a lightweight process needs to be installed on every instance of the cluster.

Some of the key features of Arc:

 * Self-contained server binary with no dependencies (windows/unix)
 * Built-in auto-update functionality
 * Uses a messaging middleware instead of direct connections
 * Multi-tenant capable secure authorisation and encryption 

## Basic Architecture of Arc

<div class="center">
![Arc Architecture](architecture.png)
</div>

Each system that should be managed via Arc runs the **Server** component of Arc. Each server maintains a persistent connection to a central **Messaging Middleware** waiting for requests it should act on. Jobs are scheduled primarily via an asynchronous HTTP **API Service**.
The Api service takes care of end-user authentication and tracks the lifecycle of scheduled jobs, store job logs and and maintains an inventory of registered servers. 

## Next Steps

* See [how Arc compares to other software](/intro/vs/index.html) to assess how it fits into your
existing infrastructure.
* Continue onwards with the [getting started guide](/intro/getting-started/install.html)
to get Arc up and running.
