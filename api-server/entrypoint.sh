#!/bin/sh
set -e

if [ "$1" = "api-server" ]; then
  # if we want to use certs but the given cert does not exist try to create one from the pki service
  # TODO: Think about authentication at some point
  PKI_HOST=${ARC_PKI_SERVICE_HOST:-$PKI_SERVICE_HOST}
  PKI_PORT=${ARC_PKI_SERVICE_PORT:-$PKI_SERVICE_PORT}
  if [ ! -z "$ARC_TLS_CLIENT_CERT" ] && [ ! -f "$ARC_TLS_CLIENT_CERT" ] && [ ! -z "$PKI_HOST" ] && [ ! -z "$COMMON_NAME" ]; then
    DIR=`dirname $ARC_TLS_CLIENT_CERT`
    echo gencert -pki-endpoint=http://$PKI_HOST:$PKI_PORT -output-dir=$DIR -cn=$COMMON_NAME transport
    gencert -pki-endpoint=http://$PKI_HOST:$PKI_PORT -output-dir=$DIR -cn=$COMMON_NAME transport
  fi

  # run migrations
  goose -env=$ARC_ENV status
  goose -env=$ARC_ENV up

fi

exec "$@"
