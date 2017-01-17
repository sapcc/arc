#!/bin/sh
set -e

if [ "$1" = "api-server" ]; then
  # run migrations
  goose -env=$ARC_ENV status
  goose -env=$ARC_ENV up

fi

exec "$@"
