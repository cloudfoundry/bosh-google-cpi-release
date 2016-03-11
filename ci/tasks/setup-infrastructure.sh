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

echo "Setting up google infrastructure..."
gcloud -q compute addresses create ${google_address_director}
gcloud -q compute addresses create ${google_address_bats}
gcloud -q compute networks create ${google_network} --mode auto
gcloud -q compute firewall-rules create ${google_firewall_internal} --description "BOSH CI Internal traffic" --network ${google_network} --source-tags ${google_firewall_internal} --target-tags ${google_firewall_internal} --allow tcp,udp,icmp
gcloud -q compute firewall-rules create ${google_firewall_external} --description "BOSH CI External traffic" --network ${google_network} --target-tags ${google_firewall_external} --allow tcp:22,tcp:443,tcp:4222,tcp:6868,tcp:25250,tcp:25555,tcp:25777,udp:53
