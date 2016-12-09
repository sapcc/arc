#!/bin/sh
# vim: set ft=sh

set -e
require_env () {
  val=$(eval echo \${$1})
  if [ -z $val ]; then
    echo "$1 environment variable missing"
    exit 1
  fi
}


require_env "API_SERVER"
require_env "API_VERSION"
require_env "ARC_VERSION"
require_env "KEYSTONE_ENDPOINT"
require_env "DOMAIN"
require_env "PROJECT"
require_env "USERNAME"
require_env "PASSWORD"

TOKEN=$(get-token)
export TOKEN

echo "Running smoke tests"
smoke

echo "Checking that agents are running the current version"
updated-test -latest-version=$ARC_VERSION

for id in $(curl -s -H X-Auth-Token:$TOKEN $API_SERVER/api/v1/agents | jq -r '.[]|.agent_id')
do
  echo "Running job test for $id:"
  AGENT_IDENTITY=$id job-test
  echo "Running fact test for $id:"
  AGENT_IDENTITY=$id fact-test
done

