#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh

cpi_release_name="bosh-google-cpi"

pushd bosh-cpi-src
  echo "Using BOSH CLI version..."
  bosh version

  echo "Creating CPI BOSH Release..."
  bosh create release --name ${cpi_release_name}
popd
