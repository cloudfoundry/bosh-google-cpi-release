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
check_param google_address_int
check_param google_address_static_int
check_param stemcell_name


# Initialize deployment artifacts
deployment_dir="${PWD}/deployment"
google_json_key=${deployment_dir}/google_key.json

export INT_STEMCELL="${deployment_dir}/stemcell.tgz"
export NETWORK_NAME=${google_network}
export CUSTOM_NETWORK_NAME=$NETWORK_NAME
export CUSTOM_SUBNETWORK_NAME=${google_subnetwork}
export PRIVATE_IP=${google_address_static_int}
export STEMCELL_URL=${stemcell_url}
# Hardcoded until a standard release cycle is made, or to be stored programatically from
# You can hardcode this for your ENV using `gcloud compute images list`
# This step in concourse pipelines/Google-BOSH-CPI-Release/resources/google-ubuntu-stemcell
export EXISTING_STEMCELL=stemcell-c9b5025e-ceb1-4a59-5553-5a1bca74866f
export TARGET_POOL=${google_target_pool}
export BACKEND_SERVICE=${google_backend_service}
export ZONE=${google_zone}
export REGION=${google_region}
export GOOGLE_PROJECT=${google_project}

echo "Setting up artifacts..."
cp ./bosh-cpi-release/*.tgz ${deployment_dir}/${cpi_release_name}.tgz
cp ./bosh-release/*.tgz ${deployment_dir}/bosh-release.tgz
cp ./stemcell/*.tgz ${deployment_dir}/stemcell.tgz

# Find external IP
echo "Looking for external IP..."
external_ip=$(gcloud compute addresses describe ${google_address_int} --format json | jq -r '.address')
export EXTERNAL_STATIC_IP=${external_ip}

echo "Creating google json key..."
echo "${google_json_key_data}" > ${google_json_key}
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

# Setup Go and run tests
export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
make testintci














