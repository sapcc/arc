#!/bin/sh
set -e

if [ -f "$ENV_FILE" ]; then
  echo Sourcing environment file "$ENV_FILE" ...
  . "$ENV_FILE"
fi

# if we want to use certs but the given cert does not exist try to create one from the pki service
# TODO: Think about authentication at some point
if [ ! -z "$ARC_TLS_CLIENT_CERT" ] && [ ! -f "$ARC_TLS_CLIENT_CERT" ] && [ ! -z "$PKI_SERVICE_HOST" ] && [ ! -z "$COMMON_NAME" ]; then
	DIR=`dirname $ARC_TLS_CLIENT_CERT`
	echo gencert -pki-endpoint=http://$PKI_SERVICE_HOST:$PKI_SERVICE_PORT -output-dir=$DIR -cn=$COMMON_NAME transport
	gencert --pki-endpoint=http://$PKI_SERVICE_HOST:$PKI_SERVICE_PORT --output-dir=$DIR --common-name=$COMMON_NAME transport
fi

exec "$@"
