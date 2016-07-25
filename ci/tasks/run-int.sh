#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_project
check_param google_region
check_param google_zone
check_param google_json_key_data
check_param google_network
check_param google_subnetwork
check_param google_target_pool
check_param google_backend_service
check_param google_address_static_int
check_param google_address_int

# Initialize deployment artifacts
google_json_key=google_key.json

export INT_STEMCELL="stemcell.tgz"
export NETWORK_NAME=${google_auto_network}
export CUSTOM_NETWORK_NAME=${google_network}
export CUSTOM_SUBNETWORK_NAME=${google_subnetwork}
export PRIVATE_IP=${google_address_static_int}
export TARGET_POOL=${google_target_pool}
export BACKEND_SERVICE=${google_backend_service}
export ZONE=${google_zone}
export REGION=${google_region}
export GOOGLE_PROJECT=${google_project}
export STEMCELL_URL=`cat stemcell/url | sed "s|gs://|https://storage.googleapis.com/|"`

# Divine the raw stemcell URL
stemcell_url_base=`cat stemcell/url | sed "s|gs://|https://storage.googleapis.com/|"`
stemcell_url_base=${stemcell_url_base/light-/}
stemcell_url_base=${stemcell_url_base/\.tgz/-raw\.tar\.gz}
export STEMCELL_URL=$stemcell_url_base

echo "Setting up artifacts..."
cp ./stemcell/*.tgz stemcell.tgz


echo "Creating google json key..."
echo "${google_json_key_data}" > ${google_json_key}
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

# Find external IP
echo "Looking for external IP..."
external_ip=$(gcloud compute addresses describe ${google_address_int} --format json | jq -r '.address')
export EXTERNAL_STATIC_IP=${external_ip}

# Setup Go and run tests
export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
env
make testintci
