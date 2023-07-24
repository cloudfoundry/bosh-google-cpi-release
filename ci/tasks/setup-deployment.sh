#!/usr/bin/env bash

set -e

source ci/ci/tasks/utils.sh

deployment_dir="${PWD}/deployment"
cpi_release_name=bosh-google-cpi

echo "Setting up artifacts..."
cp ./bosh-cpi-release/*.tgz ${deployment_dir}/${cpi_release_name}.tgz
cp ./stemcell/*.tgz ${deployment_dir}/stemcell.tgz
cp ./bosh-cli/bosh-cli-* ${deployment_dir}/bosh && chmod +x ${deployment_dir}/bosh
export BOSH_CLI=${deployment_dir}/bosh
cp -r ./bosh-deployment ${deployment_dir}
