#!/bin/bash

set -e

src_dir="$( cd "$( dirname "$0" )" && cd ../../../.. && pwd )"
workspace_dir="$( cd "${src_dir}/.." && pwd )"

# env
: ${SSH_PRIVATE_KEY:?}
: ${GCE_CREDENTIALS_JSON:?}

# inputs
bosh_cli=$( echo ${workspace_dir}/bosh-cli/bosh-cli-* )
chmod +x "${bosh_cli}"
cpi_dir="$( cd "${workspace_dir}/bosh-cpi-release" && pwd )"
stemcell_dir="$( cd "${workspace_dir}/light-stemcell" && pwd )"
terraform_config="$( cd "${workspace_dir}/terraform" && pwd )"

# outputs
output_dir="$( cd "${workspace_dir}/deployment-state" && pwd )"

mkdir -p "${output_dir}/assets/"
cp ${cpi_dir}/*.tgz "${output_dir}/assets/cpi.tgz"
cp ${stemcell_dir}/*.tgz "${output_dir}/assets/stemcell.tgz"
cp ${bosh_cli} "${output_dir}/assets/bosh"

# make ruby available for BOSH CLI template rendering
source /etc/profile.d/chruby.sh
chruby ruby-2.1.2

mbus_password="$(openssl rand -base64 24 | tr -d '[/+]')"
gce_cloud_provider_mbus="https://mbus:${mbus_password}@$(jq -r -e .skeletal_external_ip ${terraform_config}/metadata):6868"
gce_cloud_provider_agent_mbus="https://mbus:${mbus_password}@0.0.0.0:6868"

pushd "${output_dir}" > /dev/null
  echo "Deploying skeletal instance..."

  echo "${SSH_PRIVATE_KEY}" > bosh.pem # CLI has trouble with newlines in variable

  ${bosh_cli} -n interpolate \
    -v gce_cloud_provider_mbus="${gce_cloud_provider_mbus}" \
    -v gce_cloud_provider_agent_mbus="${gce_cloud_provider_agent_mbus}" \
    -v gce_credentials_json="'${GCE_CREDENTIALS_JSON}'" \
    -v ssh_private_key="bosh.pem" \
    -l "${terraform_config}/metadata" \
    "${src_dir}/ci/stemcell/light/skeletal-deployment.yml" > ./skeletal-deployment.yml

  ${bosh_cli} -n create-env ./skeletal-deployment.yml
popd > /dev/null

echo "Successfully deployed!"
