#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

deployment_dir="${PWD}/deployment"
google_json_key=${deployment_dir}/google_key.json
manifest_filename="director-manifest.yml"
manifest_state_filename="manifest-state.json"
certs=certs.yml

echo "Creating google json key..."
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

# Export prefixed variables so they are accessible
echo "Populating environment with BOSH_ prefixed vars"
export BOSH_CONFIG=${deployment_dir}/.boshconfig
export BOSH_director_username=$director_username
export BOSH_director_password=$director_password
export BOSH_cpi_release_name=$cpi_release_name
export BOSH_google_zone=$google_zone
export BOSH_google_project=$google_project
export BOSH_google_address_static_director=$google_address_static_director
export BOSH_director_ip=$director_ip
export BOSH_google_test_bucket_name=$google_test_bucket_name
export BOSH_google_network=$google_network
export BOSH_google_subnetwork_gw=$google_subnetwork_gw
export BOSH_google_subnetwork=$google_subnetwork
export BOSH_google_subnetwork_range=$google_subnetwork_range
export BOSH_google_firewall_internal=$google_firewall_internal
export BOSH_google_firewall_external=$google_firewall_external
export BOSH_google_json_key_data=$google_json_key_data

pushd ${deployment_dir}
  echo "Destroying BOSH Director..."
  ./bosh delete-env ${manifest_filename} --state ${manifest_state_filename} --vars-store ${certs} --vars-env=BOSH
popd
