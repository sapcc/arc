#!/bin/sh
set -o errexit 
if [ "$SSH_KEY" != "" ]; then
  mkdir -p ~/.ssh
  chmod 700 ~/.ssh
  echo "$SSH_KEY" > ~/.ssh/id_rsa #the quotes prevent word splitting
  chmod 600 ~/.ssh/id_rsa
fi
echo "~/.ssh/id_rsa:"
cat ~/.ssh/id_rsa
if [ "$1" = "" ]; then
  echo "usage: $0 [WEBSITE DIRECTORY]"
  exit 1
fi
cd $1
middleman deploy -b
