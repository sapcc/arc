#!/bin/bash
# vim: set ft=sh

set -e

SOURCE=${1:-source.git}
function require_env
{
  if [[ -z ${!1} ]]; then
    echo "${1} environment variable missing"
    exit 1
  fi
}

require_env ARC_VERSION
require_env CONTAINER
require_env OS_AUTH_URL
require_env OS_USERNAME
require_env OS_USER_DOMAIN_NAME
require_env OS_PASSWORD
require_env OS_PROJECT_ID

export OS_AUTH_VERSION=3

#Linux 
mkdir -p arc/linux/amd64
cp $SOURCE/arc_${ARC_VERSION}_linux_amd64 arc/linux/amd64/
cp $SOURCE/arc_${ARC_VERSION}_linux_amd64 arc/linux/amd64/latest
checksum=$(sha256sum arc/linux/amd64/latest | cut -f1 -d' ')
cat > arc/linux/amd64/latest.json <<EOF
{
  "app_id": "arc",
  "os": "linux",
  "arch": "amd64",
  "checksum": "${checksum}",
  "version": "$ARC_VERSION",
  "url":"arc_${ARC_VERSION}_linux_amd64"
}
EOF
cp arc/linux/amd64/latest.json arc/linux/amd64/${ARC_VERSION}.json

#Windows
mkdir -p arc/windows/amd64
cp $SOURCE/arc_${ARC_VERSION}_windows_amd64.exe arc/windows/amd64/
cp $SOURCE/arc_${ARC_VERSION}_windows_amd64.exe arc/windows/amd64/latest
checksum=$(sha256sum arc/windows/amd64/latest | cut -f1 -d' ')
cat > arc/windows/amd64/latest.json <<EOF
{
  "app_id": "arc",
  "os": "windows",
  "arch": "amd64",
  "checksum": "${checksum}",
  "version": "${ARC_VERSION}",
  "url":"arc_${ARC_VERSION}_windows_amd64.exe"
}
EOF
cp arc/windows/amd64/latest.json arc/windows/amd64/$ARC_VERSION.json

swift upload $CONTAINER arc/

