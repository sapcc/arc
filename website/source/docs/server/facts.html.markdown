---
layout: "docs"
page_title: "Facts"
sidebar_current: "docs-server-facts"
description: |-
  TBD
---

# Facts

Following facts are available:

- agents
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
- metadata_availability_zone
- metadata_name
- metadata_public_ipv4
- metadata_uuid
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
	agents: {"chef"=>"enabled", "execute"=>"enabled", "rpc"=>"enabled"},
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
	metadata_availability_zone: "eu-de-1b",
	metadata_name: "rel7_test",
	metadata_public_ipv4: "10.47.1.38",
	metadata_uuid: "8c2a3086-425e-47a2-a12b-d98ad31dc471",
	online: true,
	organization: "o-monsoon2",
	os: "linux",
	platform: "ubuntu",
	platform_family: "debian",
	platform_version: "14.04",
	project: "p-6167b588e"
}
```
