#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param base_os
check_param stemcell_name
check_param google_static_ip
check_param google_network
check_param bat_vcap_password
check_param bat_google_static_ip

deployment_dir="${PWD}/deployment"
bat_manifest_filename="${deployment_dir}/${base_os}-bats-manifest.yml"
bat_config_filename="${deployment_dir}/${base_os}-bats-config.yml"
private_key=${deployment_dir}/private_key.pem

export BAT_DIRECTOR=${google_static_ip}
export BAT_STEMCELL="${deployment_dir}/stemcell.tgz"
export BAT_DEPLOYMENT_SPEC="${bat_config_filename}"
export BAT_VCAP_PASSWORD="${bat_vcap_password}"
export BAT_DNS_HOST=${google_static_ip}
export BAT_INFRASTRUCTURE=google
export BAT_NETWORKING=dynamic
export BAT_VCAP_PRIVATE_KEY=${private_key}

echo "Creating private key..."
eval $(ssh-agent)
ssh-add ${private_key}

echo "Using BOSH CLI version..."
bosh version

echo "Targeting BOSH director..."
bosh -n target ${BAT_DIRECTOR}

echo "Creating ${bat_manifest_filename}..."
cat > ${bat_manifest_filename} <<EOF
---
name: <%= properties.name || "bat" %>
director_uuid: <%= properties.uuid %>

releases:
  - name: bat
    version: <%= properties.release || "latest" %>

compilation:
  workers: <%= properties.compilation_workers || 2 %>
  network: default
  reuse_compilation_vms: true
  cloud_properties:
    machine_type: <%= properties.machine_type || "n1-standard-4" %>
    root_disk_size_gb: <%= properties.root_disk_size_gb || 20 %>
    root_disk_type: <%= properties.root_disk_type || "pd-standard" %>
    <% if properties.zone %>
    zone: <%= properties.zone %>
    <% end %>

update:
  canaries: <%= properties.canaries || 1 %>
  canary_watch_time: 3000-90000
  update_watch_time: 3000-90000
  max_in_flight: <%= properties.max_in_flight || 1 %>

networks:
  <% properties.networks.each do |network| %>
  - name: <%= network.name %>
    type: <%= network.type %>
    dns: <%= p('dns').inspect %>
    cloud_properties:
      <% if network.cloud_properties.network_name %>
      network_name: <%= network.cloud_properties.network_name %>
      <% end %>
      ephemeral_external_ip: <%= network.cloud_properties.ephemeral_external_ip || false %>
      tags: <%= network.cloud_properties.tags || [] %>
  <% end %>

  - name: vip
    type: vip

resource_pools:
  - name: common
    network: default
    stemcell:
      name: <%= properties.stemcell.name %>
      version: "<%= properties.stemcell.version %>"
    cloud_properties:
      machine_type: <%= properties.machine_type || "n1-standard-4" %>
      root_disk_size_gb: <%= properties.root_disk_size_gb || 20 %>
      root_disk_type: <%= properties.root_disk_type || "pd-standard" %>
      <% if properties.zone %>
      zone: <%= properties.zone %>
      <% end %>
    <% if properties.password %>
    env:
      bosh:
        password: <%= properties.password %>
    <% end %>

jobs:
  - name: <%= properties.job || "batlight" %>
    templates: <% (properties.templates || ["batlight"]).each do |template| %>
    - name: <%= template %>
    <% end %>
    instances: <%= properties.instances %>
    resource_pool: common
    <% if properties.persistent_disk %>
    persistent_disk: <%= properties.persistent_disk %>
    <% end %>
    networks:
    <% properties.job_networks.each_with_index do |network, i| %>
      - name: <%= network.name %>
        <% if i == 0 %>
        default: [dns, gateway]
        <% end %>
    <% end %>
    <% if properties.use_vip %>
      - name: vip
        static_ips:
          - <%= properties.vip %>
    <% end %>

properties:
  batlight:
    <% if properties.batlight.fail %>
    fail: <%= properties.batlight.fail %>
    <% end %>
    <% if properties.batlight.missing %>
    missing: <%= properties.batlight.missing %>
    <% end %>
    <% if properties.batlight.drain_type %>
    drain_type: <%= properties.batlight.drain_type %>
    <% end %>
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
  vip: ${bat_google_static_ip}
  networks:
    - name: default
      type: dynamic
      cloud_properties:
        network_name: ${google_network}
        ephemeral_external_ip: true
        tags:
          - bosh-ci
EOF

pushd bats
   echo "Installing gems..."
  ./write_gemfile
  bundle install

  echo "Running BOSH Acceptance Tests..."
  bundle exec rspec spec
popd
