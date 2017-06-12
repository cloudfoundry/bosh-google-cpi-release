#!/usr/bin/env bash

set -e

: ${bat_vcap_password:?}
: ${director_password:?}
: ${stemcell_name:?}

source certification/shared/utils.sh
source /etc/profile.d/chruby.sh
chruby 2.1.7

# inputs
release_dir="$( cd $(dirname $0) && cd ../.. && pwd )"
workspace_dir="$( cd ${release_dir} && cd .. && pwd )"
ci_environment_dir="${workspace_dir}/environment"
director_config="${workspace_dir}/director-config"
bats_dir="${workspace_dir}/bats"
director_state_dir="${workspace_dir}/director-state"
bosh_cli="${workspace_dir}/bosh-cli/*bosh-cli-*"
chmod +x $bosh_cli

creds_path() { $bosh_cli int $director_state_dir/creds.yml --path="$1" ; }

metadata="$( cat ${ci_environment_dir}/metadata )"

# configuration
: ${zone:=$(    echo ${metadata} | jq --raw-output ".Zone" )}
: ${network:=$(echo ${metadata} | jq --raw-output ".CustomNetwork" )}
: ${subnetwork:=$(echo ${metadata} | jq --raw-output ".Subnetwork" )}
: ${subnetwork_cidr:=$(echo ${metadata} | jq --raw-output ".SubnetworkCIDR" )}
: ${google_subnetwork_gw:=$(echo ${metadata} | jq --raw-output ".SubnetworkGateway" )}
: ${director_external_ip:=$(echo ${metadata} | jq --raw-output ".DirectorExternalIP" )}
: ${internal_tag:=$(echo ${metadata} | jq --raw-output ".InternalTag" )}
: ${external_tag:=$(echo ${metadata} | jq --raw-output ".ExternalTag" )}
: ${bats_ip:=$( echo ${metadata} | jq --raw-output ".BATsExternalIP" )}
: ${bats_static_ip_pair:=$(echo ${metadata} | jq --raw-output ".BATsStaticIPPair" )}
: ${bats_static_ip:=$(echo ${metadata} | jq --raw-output ".BATsStaticIP" )}
: ${bats_reserved_range:=$(echo ${metadata} | jq --raw-output ".ReservedRange" )}


# outputs
output_dir="${workspace_dir}/bats-config"
bats_spec="${output_dir}/bats-config.yml"
bats_env="${output_dir}/bats.env"
ssh_key="${output_dir}/shared.pem"

# env file generation
cat > "${bats_env}" <<EOF
#!/usr/bin/env bash

export BOSH_ENVIRONMENT=${director_external_ip}
export BOSH_CLIENT="admin"
export BOSH_CLIENT_SECRET="$( creds_path /admin_password )"
export BOSH_CA_CERT="$( creds_path /director_ssl/ca )"
export BAT_DNS_HOST=${director_external_ip}
export BAT_INFRASTRUCTURE=google
export BAT_NETWORKING=dynamic
export BAT_RSPEC_FLAGS="--tag ~multiple_manual_networks --tag ~raw_ephemeral_storage --tag ~changing_static_ip"

# bosh2 ssh info
export BOSH_GW_HOST=${director_external_ip}
export BOSH_GW_USER=jumpbox
export BAT_PRIVATE_KEY="$( creds_path /jumpbox_ssh/private_key )"
EOF

# BATs spec generation
cat > "${bats_spec}" <<EOF
---
cpi: google
properties:
  stemcell:
    name: ${stemcell_name}
    version: latest
  instances: 1
  vip: ${bats_ip}
  zone: ${zone}
  static_ips: [${bats_static_ip_pair}]
  networks:
    - name: default
      static_ip: ${bats_static_ip}
      type: manual
      subnets:
      - range: ${subnetwork_cidr}
        gateway: ${google_subnetwork_gw}
        static: ${bats_reserved_range}
        cloud_properties:
          network_name: ${network}
          subnetwork_name: ${network}
          ephemeral_external_ip: true
          tags:
            - ${internal_tag}
            - ${external_tag}
EOF

cp ${director_state_dir}/shared.pem ${ssh_key}
