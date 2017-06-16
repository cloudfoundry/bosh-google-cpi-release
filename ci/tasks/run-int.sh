#!/usr/bin/env bash

set -e

# inputs
release_dir="$( cd $(dirname $0) && cd ../.. && pwd )"
source ${release_dir}/ci/tasks/utils.sh
workspace_dir="$( cd ${release_dir} && cd .. && pwd )"
environment_dir="${workspace_dir}/environment"
stemcell_dir="${workspace_dir}/stemcell"

: ${STEMCELL:=${stemcell_dir}/stemcell.tgz}
: ${METADATA_FILE:=${environment_dir}/metadata}

metadata="$( cat ${METADATA_FILE} )"

# configuration
: ${google_json_key_data:?}
echo "Creating google json key..."
echo "${google_json_key_data}" > /tmp/google_key.json

export GOOGLE_APPLICATION_CREDENTIALS=/tmp/google_key.json
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
export SERVICE_ACCOUNT=$(echo ${metadata} | jq --raw-output ".ServiceAccount" )
export EXTERNAL_STATIC_IP=$(echo ${metadata} | jq --raw-output ".IntegrationExternalIP" )
export STEMCELL_URL="$( tar -xzf ${STEMCELL} -- stemcell.MF && cat stemcell.MF | grep source_url | awk '{print $2}' )"
export CPI_ASYNC_DELETE=true

export GOPATH=${release_dir}
export PATH=${GOPATH}/bin:$PATH

check_go_version $GOPATH

pushd ${release_dir}/src/bosh-google-cpi
  make testintci
popd

rm /tmp/google_key.json
if [ -d "${environment_dir}" ]; then
  echo "Sleeping... waiting for VM's to be deleted"
  sleep $(( 60 * 5 ))
fi
