#!/bin/bash

set -e

src_dir="$( cd "$( dirname "$0" )" && cd ../../../.. && pwd )"
workspace_dir="$( cd "${src_dir}/.." && pwd )"

# inputs
deployment_state="$( cd "${workspace_dir}/deployment-state" && pwd )"

# make ruby available for BOSH CLI template rendering
source /etc/profile.d/chruby.sh
chruby ruby-2.1.2

pushd "${deployment_state}" > /dev/null
  echo "Destroying skeletal instance..."

  set +e
  ./assets/bosh -n delete-env ./skeletal-deployment.yml
  exit_code=$?
  set -e

  if [ "${exit_code}" == "0" ]; then
    echo "Successfully destroyed!"
  else
    echo "Failed to destroy deployment, continuing..."
  fi

popd > /dev/null
