#!/usr/bin/env bash

set -e

source /etc/profile.d/chruby-with-ruby-2.1.2.sh

semver=`cat version-semver/number`

cd bosh-google-cpi-boshrelease

echo "Using BOSH CLI version..."
bosh version

cpi_release_name="bosh-google-cpi"

echo "Building CPI BOSH Release..."
bosh create release --name $cpi_release_name --version $semver --with-tarball

mv dev_releases/$cpi_release_name/$cpi_release_name-$semver.tgz ../candidate/
