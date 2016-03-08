#!/usr/bin/env bash

set -e

cd ${GOPATH}/src/github.com/frodenas/bosh-google-cpi
bin/test
