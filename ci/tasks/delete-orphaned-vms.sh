#!/usr/bin/env bash

source ci/ci/tasks/utils.sh

check_param google_json_key_data
check_param google_project
check_param google_region
check_param google_zone
check_param google_auto_network
check_param google_manual_network

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
  gcloud -q compute instances delete ${instance} --delete-disks all
done

gcloud compute instances list --format json | jq -r --arg network ${google_manual_network} '.[] | select(.networkInterfaces[].network==$network) | "\(.name) --zone \(.zone)"' | while read instance; do
  echo "Deleting orphan instance ${instance}..."
  gcloud -q compute instances delete ${instance} --delete-disks all
done
