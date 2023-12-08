#!/usr/bin/env bash

set -eu -o pipefail

source ci/ci/tasks/utils.sh

check_param google_subnetwork_range
check_param stemcell_name
check_param public_key_data
check_param private_key_data

creds_dir="${PWD}/director-creds"
creds_file="${creds_dir}/${cpi_source_branch}-creds.yml"
infrastructure_metadata="${PWD}/infrastructure/metadata"

read_infrastructure

echo "Creating bats env config..."
cat > bats-config/bats.env <<EOF
export BOSH_ENVIRONMENT="${google_address_director_ip}"
export BOSH_CLIENT="admin"
export BOSH_CLIENT_SECRET="$(bosh interpolate ${creds_file} --path /admin_password)"
export BOSH_CA_CERT="$(bosh interpolate ${creds_file} --path /director_ssl/ca)"

export BAT_INFRASTRUCTURE=google
export BAT_RSPEC_FLAGS="--tag ~multiple_manual_networks --tag ~raw_ephemeral_storage --tag ~changing_static_ip --tag ~network_reconfiguration --tag ~dns"
EOF

echo "Creating bats-config..."
cat > bats-config/bats-config.yml <<EOF
---
cpi: google
properties:
  stemcell:
    name: ${stemcell_name}
    version: latest
  instances: 1
  vip: ${google_address_bats_ip}
  zone: ${google_zone}
  ssh_key_pair:
    public_key: "${public_key_data}"
    private_key: "$(echo "${private_key_data}" | sed 's/$/\\n/' | tr -d '\n')"
  static_ips: [${google_address_bats_internal_ip_pair}]
  networks:
    - name: default
      static_ip: ${google_address_bats_internal_ip}
      type: manual
      subnets:
      - range: ${google_subnetwork_range}
        gateway: ${google_subnetwork_gateway}
        static: ${google_address_bats_internal_ip_static_range}
        reserved: ${google_address_director_internal_ip},${google_address_int_internal_ip}
        cloud_properties:
          network_name: ${google_network}
          subnetwork_name: ${google_subnetwork}
          ephemeral_external_ip: true
          tags:
            - ${google_firewall_internal}
            - ${google_firewall_external}
EOF
