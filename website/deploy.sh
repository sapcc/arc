#!/bin/sh
set -o errexit
if [ "$SSH_KEY" != "" ]; then
  mkdir -p ~/.ssh
  chmod 700 ~/.ssh
  echo "$SSH_KEY" > ~/.ssh/id_rsa #the quotes prevent word splitting
  chmod 600 ~/.ssh/id_rsa
  # hack since concourse5 can not deal with git@ protocol when getting the sources
  git remote set-url origin git@gitHub.***REMOVED***:monsoon/arc.git
fi
if [ "$1" = "" ]; then
  echo "usage: $0 [WEBSITE DIRECTORY]"
  exit 1
fi
cd $1
/usr/local/bundle/bin/middleman deploy -b
