#!/usr/bin/env bash

set -e

export GOPATH=${PWD}/bosh-google-cpi-boshrelease
export PATH=${GOPATH}/bin:$PATH

cd ${PWD}/bosh-google-cpi-boshrelease/src/github.com/frodenas/bosh-google-cpi
bin/test
