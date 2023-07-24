#!/usr/bin/env bash

set -e

export HOME=/home/non-root-user
sudo chown -R non-root-user $(pwd)

source ci/ci/tasks/utils.sh

check_param google_test_bucket_name
check_param google_subnetwork_range
check_param private_key_user
check_param private_key_data
check_param google_json_key_data

creds_file="${PWD}/director-creds/${cpi_source_branch}-creds.yml"
state_file="${PWD}/director-state/${cpi_source_branch}-manifest-state.json"
cpi_release_name=bosh-google-cpi
infrastructure_metadata="${PWD}/infrastructure/metadata"
deployment_dir="${PWD}/deployment"
google_json_key=$HOME/.config/gcloud/application_default_credentials.json
private_key=${HOME}/private_key.pem

echo "Creating google json key..."
mkdir -p $HOME/.config/gcloud/
echo "${google_json_key_data}" > ${google_json_key}

read_infrastructure

export BOSH_CLI=${deployment_dir}/bosh

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

echo "Creating private key..."
echo "${private_key_data}" > ${private_key}
chmod go-r ${private_key}
eval $(ssh-agent)
ssh-add ${private_key}

echo "Generating public key from vcap private"
public_key="public.key"
openssl rsa -in ${private_key} -pubout > ${public_key}

export BOSH_CONFIG=${deployment_dir}/.boshconfig

echo "Creating ops files..."
# Use the locally built CPI
cat > "${deployment_dir}/ops_local_cpi.yml" <<EOF
---
- type: replace
  path: /releases/name=${cpi_release_name}?
  value:
    name: ${cpi_release_name}
    url: file://${deployment_dir}/${cpi_release_name}.tgz
EOF

# Use locally sourced stemcell
cat > "${deployment_dir}/ops_local_stemcell.yml" <<EOF
---
- type: replace
  path: /resource_pools/name=vms/stemcell?
  value:
    url: file://${deployment_dir}/stemcell.tgz
EOF

# Allow user vcap to SSH into director
cat > "${deployment_dir}/ops_add_vcap.yml" <<EOF
---
- type: replace
  path: /releases/name=os-conf?
  value:
    name: os-conf
    version: 18
    url: https://bosh.io/d/github.com/cloudfoundry/os-conf-release?v=18
    sha1: 78d79f08ff5001cc2a24f572837c7a9c59a0e796

- type: replace
  path: /instance_groups/name=bosh/jobs/-
  value:
    name: user_add
    release: os-conf
    properties:
      users:
      - name: vcap
        public_key: $(ssh-keygen -i -m PKCS8 -f ${public_key})
EOF

pushd ${deployment_dir}
  function finish {
    echo "Final state of director deployment:"
    echo "=========================================="
    cat ${state_file}
    echo "=========================================="

    cp -r $HOME/.bosh ./
  }
  trap finish ERR

  echo "Using bosh version..."
  ${BOSH_CLI} --version

  echo "Deploying BOSH Director..."
  ${BOSH_CLI} create-env bosh-deployment/bosh.yml \
      --state=${state_file} \
      --vars-store=${creds_file} \
      -o bosh-deployment/gcp/cpi.yml \
      -o bosh-deployment/gcp/gcs-blobstore.yml \
      -o bosh-deployment/external-ip-not-recommended.yml \
      -o ops_local_cpi.yml \
      -o ops_local_stemcell.yml \
      -o ops_add_vcap.yml \
      -v director_name=micro-google \
      -v internal_cidr=${google_subnetwork_range} \
      -v internal_gw=${google_subnetwork_gateway} \
      -v internal_ip=${google_address_director_internal_ip} \
      -v external_ip=${google_address_director_ip} \
      --var-file gcp_credentials_json=${google_json_key} \
      -v project_id=${google_project} \
      -v zone=${google_zone} \
      -v "tags=[${google_firewall_internal}, ${google_firewall_external}]" \
      -v network=${google_network} \
      -v subnetwork=${google_subnetwork} \
      -v bucket_name=${google_test_bucket_name} \
     --var-file director_gcs_credentials_json=${google_json_key} \
     --var-file agent_gcs_credentials_json=${google_json_key}

  echo "Smoke testing connection to BOSH Director"
  export BOSH_ENVIRONMENT="${google_address_director_ip}"
  export BOSH_CLIENT="admin"
  export BOSH_CLIENT_SECRET="$(${BOSH_CLI} interpolate ${creds_file} --path /admin_password)"
  export BOSH_CA_CERT="$(${BOSH_CLI} interpolate ${creds_file} --path /director_ssl/ca)"
  ${BOSH_CLI} env
  ${BOSH_CLI} login

  trap - ERR
  finish
popd
