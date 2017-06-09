#!/usr/bin/env bash

set -e

# inputs
release_dir="$( cd $(dirname $0) && cd ../.. && pwd )"
workspace_dir="$( cd ${release_dir} && cd .. && pwd )"
ci_environment_dir="${workspace_dir}/environment"
bosh_deployment="${workspace_dir}/bosh-deployment"
certification="${workspace_dir}/certification"
bosh_cli="${workspace_dir}/bosh-cli/*bosh-cli-*"
chmod +x $bosh_cli

# outputs
ci_output_dir="${workspace_dir}/director-config"

# environment
: ${google_json_key_data:?}
: ${director_password:?}
: ${METADATA_FILE:=${ci_environment_dir}/metadata}
: ${OUTPUT_DIR:=${ci_output_dir}}

if [ ! -d ${OUTPUT_DIR} ]; then
  echo -e "OUTPUT_DIR '${OUTPUT_DIR}' does not exist"
  exit 1
fi
if [ ! -f ${METADATA_FILE} ]; then
  echo -e "METADATA_FILE '${METADATA_FILE}' does not exist"
  exit 1
fi

metadata="$( cat ${METADATA_FILE} )"
tmpdir="$(mktemp -d /tmp/bosh-director-artifacts.XXXXXXXXXX)"

BOSH_RELEASE_URI="file://$( echo bosh-release/*.tgz )"
CPI_RELEASE_URI="file://$( echo cpi-release/*.tgz )"
STEMCELL_URI="file://$( echo stemcell/*.tgz )"

# configuration
: ${project_id:=$(echo ${metadata} | jq --raw-output ".ProjectID" )}
: ${zone:=$(echo ${metadata} | jq --raw-output ".Zone" )}
: ${network:=$(echo ${metadata} | jq --raw-output ".CustomNetwork" )}
: ${subnetwork:=$(echo ${metadata} | jq --raw-output ".Subnetwork" )}
: ${subnetwork_cidr:=$(echo ${metadata} | jq --raw-output ".SubnetworkCIDR" )}
: ${google_subnetwork_gw:=$(echo ${metadata} | jq --raw-output ".SubnetworkGateway" )}
: ${internal_tag:=$(echo ${metadata} | jq --raw-output ".InternalTag" )}
: ${external_tag:=$(echo ${metadata} | jq --raw-output ".ExternalTag" )}
: ${director_external_ip:=$(echo ${metadata} | jq --raw-output ".DirectorExternalIP" )}
: ${director_internal_ip:=$(echo ${metadata} | jq --raw-output ".DirectorInternalIP" )}

# keys
shared_key="shared.pem"
echo "${PRIVATE_KEY_DATA}" > "${OUTPUT_DIR}/${shared_key}"

# env file generation
cat > "${OUTPUT_DIR}/director.env" <<EOF
#!/usr/bin/env bash

export BOSH_ENVIRONMENT=${director_external_ip}
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET="${director_password}"
EOF

cat > /tmp/gcp_creds.yml <<EOF
---
project_id: ${project_id}
zone: ${zone}
network: ${network}
subnetwork: ${subnetwork}
tags: [${internal_tag},${external_tag}]
external_ip: ${director_external_ip}
internal_cidr: ${subnetwork_cidr}
internal_gw: ${google_subnetwork_gw}
internal_ip: ${director_internal_ip}
director_name: bosh
dns_recursor_ip: 169.254.169.254
EOF

cat > "${OUTPUT_DIR}/cloud-config.yml" <<EOF
EOF

${bosh_cli} interpolate \
  --ops-file ${bosh_deployment}/gcp/cpi.yml \
  --ops-file ${bosh_deployment}/powerdns.yml \
  --ops-file ${bosh_deployment}/external-ip-not-recommended.yml \
  --ops-file ${certification}/shared/assets/ops/custom-releases.yml \
  --ops-file ${certification}/gcp/assets/ops/custom-releases.yml \
  -v bosh_release_uri="${BOSH_RELEASE_URI}" \
  -v cpi_release_uri="${CPI_RELEASE_URI}" \
  -v stemcell_uri="${STEMCELL_URI}" \
  -l /tmp/gcp_creds.yml \
  -v "gcp_credentials_json='${google_json_key_data}'" \
  ${bosh_deployment}/bosh.yml > "${OUTPUT_DIR}/director.yml"

echo -e "Successfully generated manifest!"
echo -e "Manifest:    ${OUTPUT_DIR}/director.yml"
echo -e "Env:         ${OUTPUT_DIR}/director.env"
echo -e "CloudConfig: ${OUTPUT_DIR}/cloud-config.yml"
echo -e "Artifacts:   ${tmpdir}/"
