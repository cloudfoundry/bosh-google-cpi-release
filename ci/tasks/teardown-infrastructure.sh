#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_project
check_param google_region
check_param google_zone
check_param google_json_key_data
check_param google_auto_network
check_param google_target_pool
check_param google_backend_service
check_param google_region_backend_service
check_param google_network
check_param google_subnetwork
check_param google_firewall_internal
check_param google_firewall_external
check_param google_address_director_ubuntu
check_param google_address_bats_ubuntu
check_param google_address_int_ubuntu
check_param google_service_account

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
gcloud -q iam service-accounts delete ${google_service_account}@${google_project}.iam.gserviceaccount.com
gcloud compute instances list --format json | jq -r --arg network ${google_network} '.[] | select(.networkInterfaces[].network==$network) | "\(.name) --zone \(.zone)"' | while read instance; do
  echo "Deleting orphan instance ${instance}..."
  gcloud -q compute instances delete ${instance} --delete-disks all
done

gcloud compute instances list --format json | jq -r --arg network ${google_auto_network} '.[] | select(.networkInterfaces[].network==$network) | "\(.name) --zone \(.zone)"' | while read instance; do
  echo "Deleting orphan instance ${instance}..."
  gcloud -q compute instances delete ${instance} --delete-disks all
done

gcloud -q compute firewall-rules delete ${google_firewall_external}
gcloud -q compute firewall-rules delete ${google_firewall_internal}
gcloud -q compute networks subnets delete ${google_subnetwork}
gcloud -q compute networks delete ${google_network}
gcloud -q compute networks delete ${google_auto_network}
gcloud -q compute addresses delete ${google_address_director_ubuntu}
gcloud -q compute addresses delete ${google_address_bats_ubuntu}
gcloud -q compute addresses delete ${google_address_int_ubuntu}
gcloud -q compute target-pools delete ${google_target_pool}
gcloud -q compute backend-services delete ${google_backend_service}
gcloud -q compute http-health-checks delete ${google_backend_service}
gcloud -q compute instance-groups unmanaged delete ${google_backend_service}
gcloud -q compute backend-services delete ${google_region_backend_service} --region ${google_region}
gcloud -q compute health-checks delete ${google_region_backend_service}
gcloud -q compute instances delete ${google_region_backend_service} --zone ${google_zone}
gcloud -q compute instance-groups unmanaged delete ${google_region_backend_service} --zone ${google_zone}

set -e
