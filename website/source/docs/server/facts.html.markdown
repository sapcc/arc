---
layout: "docs"
page_title: "Facts"
sidebar_current: "docs-server-facts"
description: |-
  TBD
---

# Facts

Following facts are available:

- **agents:** list all available agents of the Arc node. See following section to now more about the available [agents](/docs/server/agents.html).

- **arc_version:** shows the version of the installed Arc node.

- **cert_expiration:** shows the number of hours to the expiration date of the certificate used by the Arc node.

- **default_gateway:** shows the default gateway from the instance where the Arc node is installed.

- **default_interface:** shows the default interface from the instance where the Arc node is installed.

- **domain:** DNS domain from the instance where the Arc node is installed.

- **fqdn:** shows the fully qualified domain name from the instance where the Arc node is installed.

- **hostname:** shows the hostname from the instance where the Arc node is installed.

- **identity:** shows the Arc node id.

- **init_package:** shows the available Linux Service Management (Systemd, SysV or Upstart). This fact not available for Windows images.

- **ipaddress:** shows the fixed IP address from the instance where the Arc node is installed.

- **macaddress:** shows the mac address from the instance where the Arc node is installed.

- **memory_available:** shows the RAM available for programs to allocate on the instance where the Arc node is installed.

- **memory_total:** shows the total amount of RAM on the instance where the Arc node is installed.

- **memory_used:** shows RAM used by programs on the instance where the Arc node is installed.

- **memory_used_percent:** shows the percentage of RAM used by programs on the instance where the Arc node is installed.

- **metadata_availability_zone:** shows the availability zone where the instance is deployed. This fact is provided by the metadata service.

- **metadata_name:** shows the instance name provided by the metadata service.

- **metadata_public_ipv4:** shows the floating IP address provided by the metadata service.

- **metadata_uuid:** shows the instance uuid provided by the metadata service.

- **online:** boolean value that defines if the Arc node installed on the machine is able to communicate with the broker.

- **organization:** shows the organization id.

- **os:** shows the installed operating system on the instance where the Arc node is installed.

- **platform:** shows the platform from the instance where the Arc node is installed.

- **platform_family:** shows the platform family from the instance where the Arc node is installed.

- **platform_version:** shows the platform version from the instance where the Arc node is installed.

- **project:** shows the project id.

Facts values example:

```text
{
	agents: {"chef"=>"enabled", "execute"=>"enabled", "rpc"=>"enabled"},
	arc_version: "20160118.2 (341fb82), go1.5.3",
  cert_expiration: 17495,
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
