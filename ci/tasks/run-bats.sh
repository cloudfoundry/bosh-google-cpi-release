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
check_param google_subnetwork_range
check_param google_subnetwork_gw
check_param google_firewall_internal
check_param google_firewall_external
check_param google_address_director
check_param google_address_bats
check_param google_address_static_bats
check_param google_address_static_pair_bats
check_param google_address_static_available_range_bats
check_param base_os
check_param stemcell_name
check_param bat_vcap_password
check_param private_key_data


# Initialize deployment artifacts
deployment_dir="${PWD}"
cpi_release_name=bosh-google-cpi
google_json_key=${deployment_dir}/google_key.json
private_key=${deployment_dir}/private_key.pem
manifest_filename="director-manifest.yml"
bat_manifest_filename="${deployment_dir}/${base_os}-bats-manifest.yml"
bat_config_filename="${deployment_dir}/${base_os}-bats-config.yml"

echo "Setting up artifacts..."
echo "${private_key_data}" > ${private_key}
cp ./stemcell/*.tgz stemcell.tgz

echo "Setting up artifacts..."
cp ./stemcell/*.tgz ${deployment_dir}/stemcell.tgz
echo "${private_key_data}" > ${private_key}

export BAT_STEMCELL="${deployment_dir}/stemcell.tgz"
export BAT_DEPLOYMENT_SPEC="${bat_config_filename}"
export BAT_VCAP_PASSWORD="${bat_vcap_password}"
export BAT_INFRASTRUCTURE=google
export BAT_NETWORKING=dynamic
export BAT_VCAP_PRIVATE_KEY=${private_key}

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
export BAT_DIRECTOR=${director_ip}
export BAT_DNS_HOST=${director_ip}

echo "Looking for bats IP..."
bats_ip=$(gcloud compute addresses describe ${google_address_bats} --format json | jq -r '.address')

echo "Creating private key..."
echo "${private_key_data}" > ${private_key}
chmod go-r ${private_key}
eval $(ssh-agent)
ssh-add ${private_key}

echo "Using BOSH CLI version..."
bosh version

echo "Targeting BOSH director..."
bosh -n target ${BAT_DIRECTOR}

echo "Creating ${bat_manifest_filename}..."
cat > ${bat_manifest_filename} <<EOF

EOF

echo "Creating ${bat_config_filename}..."
cat > ${bat_config_filename} <<EOF
---
cpi: google
manifest_template_path: ${bat_manifest_filename}
properties:
  uuid: $(bosh status --uuid)
  stemcell:
    name: ${stemcell_name}
    version: latest
  instances: 1
  vip: ${bats_ip}
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
  ./write_gemfile
  bundle install

  echo "Running BOSH Acceptance Tests..."
  bundle exec rspec --tag ~multiple_manual_networks --tag ~raw_ephemeral_storage --tag ~changing_static_ip spec
popd
