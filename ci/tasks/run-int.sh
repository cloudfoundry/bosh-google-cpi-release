#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

# inputs
release_dir="$( cd $(dirname $0) && cd ../.. && pwd )"
workspace_dir="$( cd ${release_dir} && cd .. && pwd )"
ci_environment_dir="${workspace_dir}/environment"

: ${METADATA_FILE:=${ci_environment_dir}/metadata}

metadata="$( cat ${METADATA_FILE} )"

# configuration
: ${google_json_key_data:?}
export GOOGLE_PROJECT=$(echo ${metadata} | jq --raw-output ".ProjectID" )
export REGION=$(echo ${metadata} | jq --raw-output ".Region" )
export ZONE=$(echo ${metadata} | jq --raw-output ".Zone" )
export NETWORK_NAME=$(echo ${metadata} | jq --raw-output ".AutoNetwork" )
export CUSTOM_NETWORK_NAME=$(echo ${metadata} | jq --raw-output ".CustomNetwork" )
export CUSTOM_SUBNETWORK_NAME=$(echo ${metadata} | jq --raw-output ".Subnetwork" )
export PRIVATE_IP=$(echo ${metadata} | jq --raw-output ".IntegrationStaticIPs" )
export TARGET_POOL=$(echo ${metadata} | jq --raw-output ".TargetPool" )
export BACKEND_SERVICE=$(echo ${metadata} | jq --raw-output ".BackendService" )
export ILB_INSTANCE_GROUP=$(echo ${metadata} | jq --raw-output ".ILBInstanceGroup" )
export REGION_BACKEND_SERVICE=$(echo ${metadata} | jq --raw-output ".RegionBackendService" )
export SERVICE_ACCOUNT=$(echo ${metadata} | jq --raw-output ".TargetPool" )
export EXTERNAL_STATIC_IP=$(echo ${metadata} | jq --raw-output ".IntegrationExternalIP" )
export INT_STEMCELL="stemcell.tgz"
export STEMCELL_URL="$( tar -xzf stemcell/stemcell.tgz -- stemcell.MF && cat stemcell.MF | grep source_url | awk '{print $2}' )"
export CPI_ASYNC_DELETE=true


# Initialize deployment artifacts
google_json_key=google_key.json

echo "Setting up artifacts..."
cp ./stemcell/*.tgz stemcell.tgz

echo "Creating google json key..."
echo "${google_json_key_data}" > ${google_json_key}
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${GOOGLE_PROJECT}
gcloud config set compute/region ${REGION}
gcloud config set compute/zone ${ZONE}

# Setup Go and run tests
export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

check_go_version $GOPATH

cd ${PWD}/bosh-cpi-src/src/bosh-google-cpi
env
make testintci
