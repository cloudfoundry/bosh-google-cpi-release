#!/usr/bin/env bash

set -e

source ci/ci/tasks/utils.sh

check_go_version ${PWD}/bosh-cpi-src

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
make test
