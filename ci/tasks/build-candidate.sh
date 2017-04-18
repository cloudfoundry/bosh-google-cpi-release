#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

cpi_release_name="bosh-google-cpi"
semver=`cat version-semver/number`

pushd bosh-cpi-src
  echo "Using BOSH CLI version..."
  bosh version

  echo "Exposing release semver to bosh-google-cpi"
  echo ${semver} > "src/bosh-google-cpi/release"

  # We have to use the --force flag because we just added the `src/bosh-google-cpi/release` file
  echo "Creating CPI BOSH Release..."
  bosh create release --name ${cpi_release_name} --version ${semver} --with-tarball --force
popd

image_path=bosh-cpi-src/dev_releases/${cpi_release_name}/${cpi_release_name}-${semver}.tgz
echo -n $(sha1sum $image_path | awk '{print $1}') > $image_path.sha1

mv ${image_path} candidate/
mv ${image_path}.sha1 candidate/
