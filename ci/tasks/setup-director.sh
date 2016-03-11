#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_project
check_param google_default_zone
check_param google_json_key_data
check_param google_static_ip
check_param google_network
check_param private_key_user
check_param private_key_data
check_param director_password
check_param director_username

deployment_dir="${PWD}/deployment"
cpi_release_name=bosh-google-cpi
google_json_key=${deployment_dir}/google_key.json
private_key=${deployment_dir}/private_key.pem
manifest_filename="director-manifest.yml"

echo "Setting up artifacts..."
cp ./bosh-cpi-release/*.tgz ${deployment_dir}/${cpi_release_name}.tgz
cp ./bosh-release/*.tgz ${deployment_dir}/bosh-release.tgz
cp ./stemcell/*.tgz ${deployment_dir}/stemcell.tgz

echo "Creating google json key..."
echo "${google_json_key_data}" > ${google_json_key}
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

echo "Creating private key..."
echo "${private_key_data}" > ${private_key}
chmod go-r ${private_key}
eval $(ssh-agent)
ssh-add ${private_key}

echo "Creating ${manifest_filename}..."
cat > "${deployment_dir}/${manifest_filename}"<<EOF
---
name: bosh
releases:
  - name: bosh
    url: file://bosh-release.tgz
  - name: ${cpi_release_name}
    url: file://${cpi_release_name}.tgz

resource_pools:
  - name: vms
    network: private
    stemcell:
      url: file://stemcell.tgz
    cloud_properties:
      machine_type: n1-standard-2
      root_disk_size_gb: 40
      root_disk_type: pd-standard
      service_scopes:
        - compute
        - devstorage.full_control

disk_pools:
  - name: disks
    disk_size: 32_768
    cloud_properties:
      type: pd-standard

networks:
  - name: private
    type: dynamic
    cloud_properties:
      network_name: ${google_network}
      tags:
        - bosh-ci
  - name: public
    type: vip

jobs:
  - name: bosh
    instances: 1

    templates:
      - name: nats
        release: bosh
      - name: redis
        release: bosh
      - name: postgres
        release: bosh
      - name: powerdns
        release: bosh
      - name: blobstore
        release: bosh
      - name: director
        release: bosh
      - name: health_monitor
        release: bosh
      - name: registry
        release: bosh
      - name: google_cpi
        release: bosh-google-cpi

    resource_pool: vms
    persistent_disk_pool: disks

    networks:
      - name: private
        default:
          - dns
          - gateway
      - name: public
        static_ips:
          - ${google_static_ip}

    properties:
      nats:
        address: 127.0.0.1
        user: nats
        password: nats-password

      redis:
        listen_address: 127.0.0.1
        address: 127.0.0.1
        password: redis-password

      postgres: &db
        listen_address: 127.0.0.1
        host: 127.0.0.1
        user: postgres
        password: postgres-password
        database: bosh
        adapter: postgres

      dns:
        address: ${google_static_ip}
        domain_name: microbosh
        db: *db
        recursor: 8.8.8.8

      registry:
        address: ${google_static_ip}
        host: ${google_static_ip}
        db: *db
        http:
          user: registry
          password: registry-password
          port: 25777
        username: registry
        password: registry-password
        port: 25777

      blobstore:
        address: ${google_static_ip}
        port: 25250
        provider: dav
        director:
          user: director
          password: director-password
        agent:
          user: agent
          password: agent-password

      director:
        address: 127.0.0.1
        name: micro-google
        db: *db
        cpi_job: google_cpi
        user_management:
          provider: local
          local:
            users:
              - name: ${director_username}
                password: ${director_password}
              - name: hm
                password: hm-password
      hm:
        director_account:
          user: hm
          password: hm-password
        resurrector_enabled: true

      google: &google_properties
        project: ${google_project}
        default_zone: ${google_default_zone}

      agent:
        mbus: nats://nats:nats-password@${google_static_ip}:4222
        ntp: *ntp
        blobstore:
           options:
             endpoint: http://${google_static_ip}:25250
             user: agent
             password: agent-password

      ntp: &ntp
        - 169.254.169.254

cloud_provider:
  template:
    name: google_cpi
    release: bosh-google-cpi

  ssh_tunnel:
    host: ${google_static_ip}
    port: 22
    user: ${private_key_user}
    private_key: ${private_key}

  mbus: https://mbus:mbus-password@${google_static_ip}:6868

  properties:
    google: *google_properties
    agent:
      mbus: https://mbus:mbus-password@0.0.0.0:6868
      blobstore:
        provider: local
        options:
          blobstore_path: /var/vcap/micro_bosh/data/cache
      ntp: *ntp
EOF

pushd ${deployment_dir}
  function finish {
    echo "Final state of director deployment:"
    echo "=========================================="
    cat director-manifest-state.json
    echo "=========================================="

    cp -r $HOME/.bosh_init ./
  }
  trap finish ERR

  chmod +x ../bosh-init/bosh-init*

  echo "Using bosh-init version..."
  ../bosh-init/bosh-init* version

  echo "Deploying BOSH Director..."
  ../bosh-init/bosh-init* deploy ${manifest_filename}

  trap - ERR
  finish
popd
