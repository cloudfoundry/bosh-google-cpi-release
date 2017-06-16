#!/usr/bin/env bash

: ${GOOGLE_APPLICATION_CREDENTIALS:?ERROR: you must set GOOGLE_APPLICATION_CREDENTIALS}

jq -e --raw-output '.modules[0].outputs | map_values(.value)' ../../ci/terraform/terraform.tfstate > /tmp/output.json
export METADATA_FILE=/tmp/output.json

wget -O /tmp/stemcell.tgz https://s3.amazonaws.com/bosh-gce-light-stemcells/light-bosh-stemcell-3421.6-google-kvm-ubuntu-trusty-go_agent.tgz
export STEMCELL=/tmp/stemcell.tgz

export google_json_key_data="$( cat $GOOGLE_APPLICATION_CREDENTIALS )"

../../ci/tasks/run-int.sh
