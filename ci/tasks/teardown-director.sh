#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

deployment_dir="${PWD}/deployment"
google_json_key=${deployment_dir}/google_key.json
manifest_filename="director-manifest.yml"

echo "Creating google json key..."
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

pushd ${deployment_dir}
  cp -r ./.bosh_init $HOME/

  chmod +x ../bosh-init/bosh-init*

  echo "Using bosh-init version..."
  ../bosh-init/bosh-init* version

  echo "Deleting BOSH Director..."
  ../bosh-init/bosh-init* delete ${manifest_filename}
popd
