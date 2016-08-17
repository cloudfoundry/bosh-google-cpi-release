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
check_param google_subnetwork_gw
check_param google_firewall_internal
check_param google_address_director
check_param base_os
check_param stemcell_name
check_param director_username
check_param director_password

deployment_dir="${PWD}/deployment"
dummy_manifest_filename="${deployment_dir}/${base_os}-dummy-manifest.yml"
google_json_key=${deployment_dir}/google_key.json

echo "Creating google json key..."
echo "${google_json_key_data}" > ${google_json_key}
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

echo "Looking for director IP..."
director_ip=$(gcloud compute addresses describe ${google_address_director} --format json | jq -r '.address')

echo "Using BOSH CLI version..."
bosh version

echo "Targeting BOSH director..."
bosh -n target ${director_ip}
bosh login ${director_username} ${director_password}

echo "Creating ${dummy_manifest_filename}..."
cat > ${dummy_manifest_filename} <<EOF
---
name: dummy
director_uuid: $(bosh status --uuid)

releases:
  - name: dummy
    version: latest

compilation:
  workers: 1
  network: private
  reuse_compilation_vms: true
  cloud_properties:
    machine_type: n1-standard-2
    zone: ${google_zone}
    root_disk_size_gb: 20
    root_disk_type: pd-standard

update:
  canaries: 1
  canary_watch_time: 3000-90000
  update_watch_time: 3000-90000
  max_in_flight: 1

resource_pools:
  - name: default
    stemcell:
      name: ${stemcell_name}
      version: latest
    network: private
    cloud_properties:
      machine_type: n1-standard-2
      zone: ${google_zone}
      root_disk_size_gb: 20
      root_disk_type: pd-standard

networks:
  - name: private
    type: manual
    subnets:
    - range: ${google_subnetwork_range}
      gateway: ${google_subnetwork_gw}
      cloud_properties:
        network_name: ${google_network}
        subnetwork_name: ${google_subnetwork}
        tags:
          - ${google_firewall_internal}

jobs:
  - name: dummy
    template: dummy
    instances: 1
    resource_pool: default
    networks:
      - name: private
        default: [dns, gateway]
EOF

echo "Uploading BOSH Dummy Release..."

pushd dummy-boshrelease
  bosh -n create release --force
  bosh -n upload release --skip-if-exists
popd

echo "Uploading Stemcell..."
bosh -n upload stemcell ${deployment_dir}/stemcell.tgz --skip-if-exists

echo "Deploying Dummy Release..."
bosh -d ${dummy_manifest_filename} -n deploy

echo "Deleting Dummy Release..."
bosh -n delete deployment dummy

echo "Cleaning up artifacts..."
bosh -n cleanup --all
