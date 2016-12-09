#!/bin/bash
# vim: set ft=sh

set -e

VERSION=$(cat arc.version/current)
function require_env
{
  if [[ -z ${!1} ]]; then
    echo "${1} environment variable missing"
    exit 1
  fi
}

require_env CONTAINER
require_env SOURCE_CONTAINER
require_env OS_AUTH_URL
require_env OS_USERNAME
require_env OS_PASSWORD
require_env OS_USER_DOMAIN_NAME
require_env OS_PROJECT_ID

export OS_AUTH_VERSION=3

eval $(swift auth)

set -x

curl -f -X COPY $OS_STORAGE_URL/$SOURCE_CONTAINER/arc/windows/amd64/arc_${VERSION}_windows_amd64.exe \
  -H "Destination: $CONTAINER/arc/windows/amd64/arc_${VERSION}_windows_amd64.exe" \
  -H "X-Auth-Token: $OS_AUTH_TOKEN"
curl -f -X COPY $OS_STORAGE_URL/$SOURCE_CONTAINER/arc/windows/amd64/latest \
  -H "Destination: $CONTAINER/arc/windows/amd64/latest" \
  -H "X-Auth-Token: $OS_AUTH_TOKEN"
curl -f -X COPY $OS_STORAGE_URL/$SOURCE_CONTAINER/arc/windows/amd64/${VERSION}.json \
  -H "Destination: $CONTAINER/arc/windows/amd64/${VERSION}.json" \
  -H "X-Auth-Token: $OS_AUTH_TOKEN"
curl -f -X COPY $OS_STORAGE_URL/$SOURCE_CONTAINER/arc/windows/amd64/latest.json \
  -H "Destination: $CONTAINER/arc/windows/amd64/latest.json" \
  -H "X-Auth-Token: $OS_AUTH_TOKEN"

curl -f -X COPY $OS_STORAGE_URL/$SOURCE_CONTAINER/arc/linux/amd64/arc_${VERSION}_linux_amd64 \
  -H "Destination: $CONTAINER/arc/linux/amd64/arc_${VERSION}_linux_amd64" \
  -H "X-Auth-Token: $OS_AUTH_TOKEN"
curl -f -X COPY $OS_STORAGE_URL/$SOURCE_CONTAINER/arc/linux/amd64/latest \
  -H "Destination: $CONTAINER/arc/linux/amd64/latest" \
  -H "X-Auth-Token: $OS_AUTH_TOKEN"
curl -f -X COPY $OS_STORAGE_URL/$SOURCE_CONTAINER/arc/linux/amd64/${VERSION}.json \
  -H "Destination: $CONTAINER/arc/linux/amd64/${VERSION}.json" \
  -H "X-Auth-Token: $OS_AUTH_TOKEN"
curl -f -X COPY $OS_STORAGE_URL/$SOURCE_CONTAINER/arc/linux/amd64/latest.json \
  -H "Destination: $CONTAINER/arc/linux/amd64/latest.json" \
  -H "X-Auth-Token: $OS_AUTH_TOKEN"
