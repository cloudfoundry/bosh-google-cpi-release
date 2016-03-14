#!/usr/bin/env bash

set -e

check_param build_number
check_param os_name
check_param os_version
check_param os_image

pushd bosh-src
  echo "Installing gems..."
  bundle install --local

  echo "Creating stemcell..."
  CANDIDATE_BUILD_NUMBER=${build_number} bundle exec rake stemcell:build[google,kvm,${os_name},${os_version},go,bosh-os-images,${os_image}]
popd
