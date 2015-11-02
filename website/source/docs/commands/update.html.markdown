---
layout: "docs"
page_title: "Commands: Update"
sidebar_current: "docs-commands-update"
description: The `update` command check for new updates and update to the last version.
---

# Arc Update

Command: `arc update`

## Description

The `update` command checks for the last version available, asks for user confirmation and triggers an update. When
the update is being triggered the existing Arc binary is replaced with the new one.

## Usage

Usage: `arc update [options]`

The following command-line options are available for this command.
Every option is optional:

* `--force, -f` - Forces an update without any user confirmation.

* `--update-uri` - Update server uri.

* `--no-update, -n` - Just return the last version available, no update is being triggered.

