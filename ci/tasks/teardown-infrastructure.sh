#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_json_key_data

infrastructure_metadata="${PWD}/infrastructure/metadata"
read_infrastructure

echo "Creating google json key..."
mkdir -p $HOME/.config/gcloud/
echo "${google_json_key_data}" > $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

gcloud compute instances list --format json | jq -r --arg network ${google_auto_network} '.[] | select(.networkInterfaces[].network==$network) | "\(.name) --zone \(.zone)"' | while read instance; do
  echo "Deleting orphan instance ${instance}..."
  gcloud -q compute instances delete ${instance} --delete-disks all &
done

gcloud compute instances list --format json | jq -r --arg network ${google_network} '.[] | select(.networkInterfaces[].network==$network) | "\(.name) --zone \(.zone)"' | while read instance; do
  echo "Deleting orphan instance ${instance}..."
  gcloud -q compute instances delete ${instance} --delete-disks all &
done

wait
set -e
