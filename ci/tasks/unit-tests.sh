#!/usr/bin/env bash

set -e

export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

cd ${PWD}/bosh-cpi-src/src/github.com/frodenas/bosh-google-cpi
bin/test
