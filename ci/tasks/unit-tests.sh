#!/usr/bin/env bash

set -e

: ${GO_VERSION:?}

export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

current=$(go version)
if [[ "$current" != *"$GO_VERSION"* ]]; then
  echo "Go version is incorrect"
  exit 1
fi

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
make test
