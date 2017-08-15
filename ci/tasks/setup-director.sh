#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_project
check_param google_region
check_param google_zone
check_param google_json_key_data
check_param google_test_bucket_name
check_param google_network
check_param google_subnetwork
check_param google_subnetwork_range
check_param google_subnetwork_gw
check_param google_firewall_internal
check_param google_firewall_external
check_param google_address_director
check_param google_address_static_director
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

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

echo "Looking for director IP..."
director_ip=$(gcloud compute addresses describe ${google_address_director} --format json | jq -r '.address')

echo "Creating private key..."
echo "${private_key_data}" > ${private_key}
chmod go-r ${private_key}
eval $(ssh-agent)
ssh-add ${private_key}

echo "Generating public key from vcap private"
public_key="public.key"
openssl rsa -in ${private_key} -pubout > ${public_key}

# Export prefixed variables so they are accessible
echo "Populating environment with BOSH_ prefixed vars"
export BOSH_director_username=$director_username
export BOSH_director_password=$director_password
export BOSH_cpi_release_name=$cpi_release_name
export BOSH_google_zone=$google_zone
export BOSH_google_project=$google_project
export BOSH_google_address_static_director=$google_address_static_director
export BOSH_director_ip=$director_ip
export BOSH_google_test_bucket_name=$google_test_bucket_name
export BOSH_google_network=$google_network
export BOSH_google_subnetwork_gw=$google_subnetwork_gw
export BOSH_google_subnetwork=$google_subnetwork
export BOSH_google_subnetwork_range=$google_subnetwork_range
export BOSH_google_firewall_internal=$google_firewall_internal
export BOSH_google_firewall_external=$google_firewall_external

export BOSH_google_json_key_data=$google_json_key_data

echo "Creating ${manifest_filename}..."
cat > "${deployment_dir}/${manifest_filename}"<<EOF
---
name: bosh
releases:
  - name: bosh
    url: file://bosh-release.tgz
  - name: ((cpi_release_name))
    url: file://((cpi_release_name)).tgz
  - name: os-conf
    url: https://bosh.io/d/github.com/cloudfoundry/os-conf-release?v=12
    sha1: af5a2c9f228b9d7ec4bd051d71fef0e712fa1549

resource_pools:
  - name: vms
    network: private
    stemcell:
      url: file://stemcell.tgz
    cloud_properties:
      zone: ((google_zone))
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
    type: manual
    subnets:
    - range: ((google_subnetwork_range))
      gateway: ((google_subnetwork_gw))
      cloud_properties:
        network_name: ((google_network))
        subnetwork_name: ((google_subnetwork))
        tags:
          - ((google_firewall_internal))
          - ((google_firewall_external))
  - name: public
    type: vip

instance_groups:
- name: bosh
  instances: 1

  jobs:
  - name: nats
    release: bosh
  - name: postgres-9.4
    release: bosh
  - name: blobstore
    release: bosh
  - name: director
    release: bosh
  - name: health_monitor
    release: bosh
  - name: powerdns
    release: bosh
  - name: google_cpi
    release: bosh-google-cpi
  - name: user_add
    release: os-conf

  resource_pool: vms
  persistent_disk_pool: disks

  networks:
    - name: private
      static_ips: [((google_address_static_director))]
      default:
        - dns
        - gateway
    - name: public
      static_ips:
        - ((director_ip))

  properties:
    nats:
      address: 127.0.0.1
      user: nats
      password: nats-password

    postgres: &db
      listen_address: 127.0.0.1
      host: 127.0.0.1
      user: postgres
      password: postgres-password
      database: bosh
      adapter: postgres

    dns:
      address: ((google_address_static_director))
      domain_name: microbosh
      db: *db
      recursor: 169.254.169.254

    registry:
      address: ((google_address_static_director))
      host: ((google_address_static_director))
      db: *db
      http:
        user: registry

    blobstore:
      provider: gcs
      bucket_name: ((google_test_bucket_name))
      credentials_source: static
      json_key: |
        $(echo $google_json_key_data | tr -d '\n')
      address: ((google_address_static_director))
      director:
        user: director
        password: director-password
      agent:
        user: agent
        password: agent-password
      port: 25250

    director:
      address: 127.0.0.1
      name: micro-google
      db: *db
      cpi_job: google_cpi
      ssl:
        key: ((director_ssl.private_key))
        cert: ((director_ssl.certificate))
      user_management:
        provider: local
        local:
          users:
            - name: ((director_username))
              password: ((director_password))
            - name: hm
              password: hm-password
    hm:
      director_account:
        user: hm
        password: hm-password
      resurrector_enabled: true

    google: &google_properties
      project: ((google_project))

    users:
      - name: vcap
        public_key: $(ssh-keygen -i -m PKCS8 -f ${public_key})

    agent:
      mbus: nats://nats:nats-password@((google_address_static_director)):4222
      ntp: *ntp
      blobstore:
          options:
            endpoint: http://((google_address_static_director)):25250
            user: agent
            password: agent-password

    ntp: &ntp
      - 169.254.169.254

cloud_provider:
  template:
    name: google_cpi
    release: ((cpi_release_name))

  mbus: https://mbus:mbus-password@((director_ip)):6868

  properties:
    google: *google_properties
    agent: {mbus: "https://mbus:mbus-password@0.0.0.0:6868"}
    blobstore: {provider: local, path: /var/vcap/micro_bosh/data/cache}
    ntp: *ntp

misc:
  ca_cert: ((director_ssl.ca))

EOF

cert_template=certs.yml.tpl
echo "Creating ${cert_template}..."
cat > "${deployment_dir}/${cert_template}"<<EOF
variables:
- name: default_ca
  type: certificate
  options:
    is_ca: true
    common_name: bosh_ca
- name: director_ssl
  type: certificate
  options:
    ca: default_ca
    common_name: ((internal_ip))
    alternative_names: [((internal_ip))]
EOF

pushd ${deployment_dir}
  function finish {
    cp director-manifest-state.json manifest-state.json
    echo "Final state of director deployment:"
    echo "=========================================="
    cat manifest-state.json
    echo "=========================================="

    cp -r $HOME/.bosh ./
  }
  trap finish ERR

  echo "Fetching bosh-cli V2"
  curl -Ls \
    https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.28-linux-amd64 \
    -o bosh2
  chmod +x bosh2

  echo "Using bosh2 version..."
  ./bosh2 --version

  echo "Generating certificates"
  certs=certs.yml
  ./bosh2 interpolate ${cert_template} -v internal_ip=${director_ip} --vars-store ${certs}

  echo "Deploying BOSH Director..."
  ./bosh2 create-env ${manifest_filename} --vars-store ${certs} --vars-env=BOSH

  echo "Logging into BOSH Director"
  # We need to fetch and specify the CA certificate as bosh-cli V2
  # strictly validates certificate with no insecure option.
  ./bosh2 interpolate certs.yml --path /director_ssl/ca > ca_cert.pem
  ./bosh2 alias-env micro-google --environment ${director_ip} --ca-cert ca_cert.pem

  # We have to export these to get non-interactive login
  export BOSH_CLIENT=$BOSH_director_username
  export BOSH_CLIENT_SECRET=$BOSH_director_password
  ./bosh2 login -e micro-google

  trap - ERR
  finish
popd
