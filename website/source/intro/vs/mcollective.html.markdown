---
layout: "intro"
page_title: "Arc vs. MCollective"
sidebar_current: "vs-other-mco"
description: |-
 TBD
---

# Arc vs. MCollective 

Arc is heavily influenced by [MCollective](http://docs.puppetlabs.com/mcollective/), also known as The Marionette Collective, and shares a similar architecture and principles in many ways:

* Both use a publish/subscribe messaging middleware
* Arc has a similar plugin concept for agents/actions that can be remotely executed
* Arc implements a fact collection/registration mechanism comparable to MCollective 

In fact the initial development of Arc was largely driven by the experience of running a large MCollective installation and tries to address some specific problems that surfaced while operating MCollective:

### Support for long running jobs
TBD

### Self contained installer
TBD

### Multi-tenant capable message middleware
TBD

### Auto updates
TBD

### Miscellaneous
* Async facts, only transmit when changed
* Track server availability (MQTT Last will, etc)
