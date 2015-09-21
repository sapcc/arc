#!/bin/sh
set -e

if [ -f "$ENV_FILE" ]; then
  echo Sourcing environment file "$ENV_FILE" ...
  . "$ENV_FILE"
fi

if [ ! -z $PKI_ENDPOINT ]; then
	DIR=`dirname $ARC_TLS_CA_CERT`
	gencert --pki-endpoint=$PKI_ENDPOINT --output-dir=$DIR --common-name=$COMMON_NAME
fi

exec "$@"
