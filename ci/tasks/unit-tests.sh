#!/usr/bin/env bash

set -e

source ci/ci/tasks/utils.sh

export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

check_go_version $GOPATH

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
make test
