---
layout: "docs"
page_title: "Facts"
sidebar_current: "docs-server-facts"
description: |-
  TBD
---

# Facts 

Following facts are available:

- arc_version
- default_gateway
- default_interface
- domain
- fqdn
- hostname
- identity
- init_package
- ipaddress
- macaddress
- memory_available
- memory_total
- memory_used
- memory_used_percent
- online
- organization
- os
- platform
- platform_family
- platform_version
- project

Facts values example:

```text
{
	arc_version: "20160118.2 (341fb82), go1.5.3",
	default_gateway: "10.44.57.1",
	default_interface: "eth0",
	domain: "4.lab.***REMOVED***",
	fqdn: "mo-f39dcd562.4.lab.***REMOVED***",
	hostname: "mo-f39dcd562",
	identity: "1164dd36-e595-4f58-9685-c48935be7261",
	init_package: "upstart",
	ipaddress: "10.44.57.94",
	macaddress: "00:50:56:8c:d9:ca",
	memory_available: 869340000,
	memory_total: 1010428000,
	memory_used: 734800000,
	memory_used_percent: 14,
	online: true,
	organization: "o-monsoon2",
	os: "linux",
	platform: "ubuntu",
	platform_family: "debian",
	platform_version: "14.04",
	project: "p-6167b588e"
}
```