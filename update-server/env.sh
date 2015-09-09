#!/bin/sh
set -e

if [ -f "$ENV_FILE" ]; then
  echo Sourcing environment file "$ENV_FILE" ...
  . "$ENV_FILE"
fi

exec "$@"
