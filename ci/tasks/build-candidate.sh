#!/usr/bin/env bash

set -e

source bosh-google-cpi-boshrelease/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

cpi_release_name="bosh-google-cpi"
semver=`cat version-semver/number`

pushd bosh-google-cpi-boshrelease
  echo "Using BOSH CLI version..."
  bosh version

  echo "Creating CPI BOSH Release..."
  bosh create release --name $cpi_release_name --version $semver --with-tarball
popd

mv bosh-google-cpi-boshrelease/dev_releases/$cpi_release_name/$cpi_release_name-$semver.tgz candidate/
