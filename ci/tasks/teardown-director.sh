#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

creds_file="${PWD}/director-creds/creds.yml"
deployment_dir="${PWD}/deployment"
google_json_key=${deployment_dir}/google_key.json
manifest_filename="director-manifest.yml"
manifest_state_filename="manifest-state.json"
infrastructure_metadata="${PWD}/infrastructure/metadata"

read_infrastructure

echo "Creating google json key..."
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

# Export prefixed variables so they are accessible
echo "Populating environment with BOSH_ prefixed vars"
export BOSH_CONFIG=${deployment_dir}/.boshconfig

pushd ${deployment_dir}
  echo "Destroying BOSH Director..."
  ./bosh delete-env bosh-deployment/bosh.yml \
      --state=${manifest_state_filename} \
      --vars-store=${creds_file} \
      -o bosh-deployment/gcp/cpi.yml \
      -o bosh-deployment/gcp/gcs-blobstore.yml \
      -o bosh-deployment/external-ip-not-recommended.yml \
      -o ops_local_cpi.yml \
      -o ops_local_stemcell.yml \
      -o ops_add_vcap.yml \
      -v director_name=micro-google \
      -v internal_cidr=${google_subnetwork_range} \
      -v internal_gw=${google_subnetwork_gw} \
      -v internal_ip=${google_address_static_director} \
      -v external_ip=${google_address_director_ip} \
      --var-file gcp_credentials_json=${google_json_key} \
      -v project_id=${google_project} \
      -v zone=${google_zone} \
      -v "tags=[${google_firewall_internal}, ${google_firewall_external}]" \
      -v network=${google_network} \
      -v subnetwork=${google_subnetwork} \
      -v bucket_name=${google_test_bucket_name} \
     --var-file director_gcs_credentials_json=${google_json_key} \
     --var-file agent_gcs_credentials_json=${google_json_key}
popd
