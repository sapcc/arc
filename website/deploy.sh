#!/bin/sh
set -o errexit
if [ "$SSH_KEY" != "" ]; then
  mkdir -p ~/.ssh
  chmod 700 ~/.ssh
  echo "$SSH_KEY" > ~/.ssh/id_rsa #the quotes prevent word splitting
  chmod 600 ~/.ssh/id_rsa
else
  curl -fL http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt | tr -d '\r' > /usr/local/share/ca-certificates/SAP_Global_Root_CA.crt && update-ca-certificates
fi
if [ "$1" = "" ]; then
  echo "usage: $0 [WEBSITE DIRECTORY]"
  exit 1
fi
cd $1
/usr/local/bundle/bin/middleman deploy -b
