#!/usr/bin/env bash

set -e

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
make test
