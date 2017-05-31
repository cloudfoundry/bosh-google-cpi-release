#!/usr/bin/env bash

set -e

export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

release_go_version="$(cat "$GOPATH/packages/golang/spec" | \
  grep linux | awk '{print $2}' | sed "s/golang\/go\(.*\)\.linux-amd65.tar.gz/\1/")"

current=$(go version)
if [[ "$current" != *"$release_go_version"* ]]; then
  echo "Go version is incorrect"
  exit 1
fi

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
make test
