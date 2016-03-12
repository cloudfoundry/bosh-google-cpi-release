#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_project
check_param google_region
check_param google_zone
check_param google_json_key_data
check_param google_network
check_param google_firewall_internal
check_param google_firewall_external
check_param google_address_director
check_param google_address_bats

echo "Creating google json key..."
mkdir -p $HOME/.config/gcloud/
echo "${google_json_key_data}" > $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

echo "Tearing down google infrastructure..."
set +e
for instance in $(gcloud compute instances list --format json | jq --arg network ${google_network} -r '.[] | select(.networkInterfaces[].network==$network) | .name'); do
  echo "Deleting orphan instance ${instance}..."
  gcloud -q compute instances delete ${instance} --delete-disks all
done
gcloud -q compute firewall-rules delete ${google_firewall_external}
gcloud -q compute firewall-rules delete ${google_firewall_internal}
gcloud -q compute networks delete ${google_network}
gcloud -q compute addresses delete ${google_address_bats}
gcloud -q compute addresses delete ${google_address_director}
set -e
