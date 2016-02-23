#!/bin/bash

DOMAIN=${DOMAIN:-monsoon2}
PROJECT=${PROJECT:-Arc_Test}
USERNAME=${USERNAME:-$USER}
KEYSTONE_ENDPOINT=${KEYSTONE_ENDPOINT:-https://identity.***REMOVED***:5000/v3}

function HELP {
  echo usage $0 [-d DOMAIN] [-p PROJECT] [-u USER] [-e KEYSTONE_ENDPOINT]
}

while getopts u:d:p:e:vh FLAG; do
  case $FLAG in
    u)
      USERNAME=$OPTARG
      ;;
    d)
      DOMAIN=$OPTARG
      ;;
    p)
      PROJECT=$OPTARG
      ;;
    e)
      KEYSTONE_ENDPOINT=$OPTARG
      ;;
    v)
      VERBOSE=1
      ;;
    h)
      HELP
      exit 0
      ;;
    \?)
      HELP
      exit 1
      ;;
  esac
done

if [ -z "$PASSWORD" ]; then
  read -p "Enter password for user $USERNAME: " -s PASSWORD
  echo
fi

json=$(curl --silent -D /dev/stderr \
      -H "Content-Type: application/json" \
      -d '
{ "auth": {
    "identity": {
      "methods": ["password"],
      "password": {
        "user": {
          "name": "'$USERNAME'",
          "domain": { "name": "'$DOMAIN'"   },
          "password": "'$PASSWORD'"
        }
      }
    },
    "scope": {
      "project": {
        "name": "'$PROJECT'",
        "domain": {"name":"'$DOMAIN'"} 
      }
    }
  }
}
' \
  $KEYSTONE_ENDPOINT/auth/tokens?nocatalog \
  2> >(grep -E HTTP\|X-Subject-Token >&2))


if echo $json |grep -q error; then
  echo $json | jq -r .error.message
  exit 1
fi
if [ -n "$VERBOSE" ]; then
  echo "$json" | jq .
else
  echo -n "User: "
  echo $json | jq -r ".token.user.id"
  echo -n "Project: "
  echo $json | jq -r ".token.project.id"
  echo -n "Domain: "
  echo $json | jq -r ".token.project.domain_id"
  echo -n "Roles: "
  echo $json | jq -r '.token.roles |map(.name)|join(", ")'
fi
