#!/usr/bin/env bash

set -e

source /etc/profile.d/chruby.sh
chruby 2.1.7

# input
input_dir=$(realpath director-state/)
bosh_cli=$(realpath bosh-cli/bosh-cli-*)
chmod +x $bosh_cli

if [ ! -e "${input_dir}/director-state.json" ]; then
  echo "director-state.json does not exist, skipping..."
  exit 0
fi

if [ -d "${input_dir}/.bosh" ]; then
  # reuse compiled packages
  cp -r ${input_dir}/.bosh $HOME/
fi

pushd ${input_dir} > /dev/null
  source director.env
  : ${BOSH_ENVIRONMENT:?}
  : ${BOSH_CLIENT:?}
  : ${BOSH_CLIENT_SECRET:?}
  export BOSH_CA_CERT="${input_dir}/ca_cert.pem"

  echo "deleting all deployments"
  $bosh_cli deployments | awk '{print $1}' | xargs --no-run-if-empty -n 1 $bosh_cli -n delete-deployment --force -d
  echo "cleaning up bosh BOSH Director..."
  $bosh_cli -n clean-up --all
  echo "deleting existing BOSH Director VM..."
  $bosh_cli -n delete-env --vars-store "${input_dir}/creds.yml" -v director_name=bosh director.yml
popd
