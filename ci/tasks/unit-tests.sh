#!/usr/bin/env bash

set -e

source ci/ci/tasks/utils.sh

check_go_version ${PWD}/bosh-google-cpi-release

cd ${PWD}/bosh-google-cpi-release/src/bosh-google-cpi
make test
