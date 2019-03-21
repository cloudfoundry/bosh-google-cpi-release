#!/usr/bin/env bash

set -e

source ci/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.6.1.sh

check_param google_json_key_data

# Initialize deployment artifacts
google_json_key=google_key.json
infrastructure_metadata="${PWD}/infrastructure/metadata"

read_infrastructure

# Stemcell stuff
export STEMCELL_VERSION=`cat stemcell/version`
export STEMCELL_FILE=`pwd`/stemcell/image.tgz
pushd stemcell
  tar -zxvf stemcell.tgz
  mv image image.tgz
popd

export NETWORK_NAME=${google_auto_network}
export CUSTOM_NETWORK_NAME=${google_network}
export CUSTOM_SUBNETWORK_NAME=${google_subnetwork}
export PRIVATE_IP=${google_address_int_internal_ip}
export TARGET_POOL=${google_target_pool}
export BACKEND_SERVICE=${google_backend_service}
export REGION_BACKEND_SERVICE=${google_region_backend_service}
export ILB_INSTANCE_GROUP=${google_region_backend_service}
export ZONE=${google_zone}
export REGION=${google_region}
export GOOGLE_PROJECT=${google_project}
export SERVICE_ACCOUNT=${google_service_account}
export EXTERNAL_STATIC_IP=${google_address_int_ip}
export CPI_ASYNC_DELETE=true

echo "Creating google json key..."
echo "${google_json_key_data}" > ${google_json_key}
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

check_go_version ${PWD}/bosh-cpi-src

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
make testintci
