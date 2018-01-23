#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_json_key_data
check_param google_subnetwork_range
check_param google_subnetwork_gw
check_param google_address_static_bats
check_param google_address_static_pair_bats
check_param google_address_static_available_range_bats
check_param base_os
check_param stemcell_name
check_param private_key_data

# Initialize deployment artifacts
deployment_dir="${PWD}"
creds_dir="${PWD}/director-creds"
creds_file="${creds_dir}/creds.yml"
cpi_release_name=bosh-google-cpi
google_json_key=${deployment_dir}/google_key.json
private_key=${deployment_dir}/private_key.pem
bat_config_filename="${deployment_dir}/bat.yml"
infrastructure_metadata="${PWD}/infrastructure/metadata"

read_infrastructure

echo "Setting up artifacts..."
echo "${private_key_data}" > ${private_key}
cp ./stemcell/*.tgz stemcell.tgz

echo "Setting up artifacts..."
cp ./stemcell/*.tgz ${deployment_dir}/stemcell.tgz

echo "${private_key_data}" > ${private_key}

export BAT_STEMCELL="${deployment_dir}/stemcell.tgz"
export BAT_DEPLOYMENT_SPEC="${bat_config_filename}"
export BAT_BOSH_CLI=/usr/bin/bosh2
export BAT_DNS_HOST=${google_address_director_ip}
export BAT_INFRASTRUCTURE=google
export BAT_NETWORKING=dynamic
export BAT_PRIVATE_KEY=${private_key}

echo "Creating google json key..."
echo "${google_json_key_data}" > ${google_json_key}
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

export BAT_DIRECTOR=${google_address_director_ip}
export BAT_DNS_HOST=${google_address_static_director}

echo "Creating private key..."
echo "${private_key_data}" > ${private_key}
chmod go-r ${private_key}
eval $(ssh-agent)
ssh-add ${private_key}

echo "Using BOSH CLI version..."
${BAT_BOSH_CLI} --version

echo "Setting up BOSH v2..."
export BOSH_ENVIRONMENT="${google_address_director_ip}"
export BOSH_CLIENT="admin"
export BOSH_CLIENT_SECRET="$(${BAT_BOSH_CLI} interpolate ${creds_file} --path /admin_password)"
export BOSH_CA_CERT="$(${BAT_BOSH_CLI} interpolate ${creds_file} --path /director_ssl/ca)"

echo "Testing connection to director"
${BAT_BOSH_CLI} env
${BAT_BOSH_CLI} login

echo "Creating ${bat_config_filename}..."
cat > ${bat_config_filename} <<EOF
---
cpi: google
properties:
  stemcell:
    name: ${stemcell_name}
    version: latest
  instances: 1
  vip: ${google_address_bats_ip}
  zone: ${google_zone}
  static_ips: [${google_address_static_pair_bats}]
  networks:
    - name: default
      static_ip: ${google_address_static_bats}
      type: manual
      subnets:
      - range: ${google_subnetwork_range}
        gateway: ${google_subnetwork_gw}
        static: ${google_address_static_available_range_bats}
        cloud_properties:
          network_name: ${google_network}
          subnetwork_name: ${google_subnetwork}
          ephemeral_external_ip: true
          tags:
            - ${google_firewall_internal}
            - ${google_firewall_external}
EOF

pushd bats
  echo "Installing gems..."
  bundle install

  echo "Running BOSH Acceptance Tests..."
  # Disable Unsupported by google cpi (multiple_manual_networks) and deprecated specs
  bundle exec rspec --tag ~multiple_manual_networks --tag ~raw_ephemeral_storage --tag ~changing_static_ip --tag ~network_reconfiguration spec
popd
