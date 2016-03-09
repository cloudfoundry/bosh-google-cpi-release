#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

deployment_dir="${PWD}/deployment"
manifest_filename="director-manifest.yml"

pushd ${deployment_dir}
  cp -r ./.bosh_init $HOME/

  chmod +x ../bosh-init/bosh-init*

  echo "Using bosh-init version..."
  ../bosh-init/bosh-init* version

  echo "Deleting BOSH Director..."
  ../bosh-init/bosh-init* delete ${manifest_filename}
popd
