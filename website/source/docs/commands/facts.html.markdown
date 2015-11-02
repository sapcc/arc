---
layout: "docs"
page_title: "Commands: Facts"
sidebar_current: "docs-commands-facts"
description: Discover and list facts on this system.
---

# Arc Facts

Command: `arc facts`

## Description

Discover and list facts on this system. The `facts` command collects information from the system to provide a simple
and easy to understand view of the machine where Arc is running.

## Output

The output of the Arc facts command is shown as a JSON. Attributes with `null` value are Facts not implemented on the running system.
The following is some example output from Arc running on a Linux machine:

    {
       "arc_version": "0.1.0-dev(6ae2a88)",
       "default_gateway": "10.97.16.1",
       "default_interface": "eth0",
       "domain": "mo.sap.corp",
       "fqdn": "mo-instance.***REMOVED***",
       "hostname": "mo-56125878f",
       "identity": "linux",
       "init_package": "upstart",
       "ipaddress": "10.97.27.63",
       "macaddress": "00:50:56:b3:7e:8d",
       "memory_available": 901268000,
       "memory_total": 1012568000,
       "memory_used": 720168000,
       "memory_used_percent": 11,
       "organization": "test-org",
       "os": "linux",
       "platform": "redhat",
       "platform_family": "rhel",
       "platform_version": "6.6",
       "project": "test-project"
    }
