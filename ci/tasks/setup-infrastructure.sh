#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_project
check_param google_default_zone
check_param google_json_key_data
check_param google_network

echo "Creating google json key..."
mkdir -p $HOME/.config/gcloud/
echo "${google_json_key_data}" > $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set project ${google_default_zone}

echo "Setting up google infrastructure..."
gcloud compute addresses create bosh-ci-director
gcloud compute networks create ${google_network} --mode auto
gcloud compute firewall-rules create bosh-ci-intenal  --description "BOSH CI Internal traffic" --network ${google_network} --source-tags bosh-ci-intenal --target-tags bosh-ci-intenal --allow tcp,udp,icmp
gcloud compute firewall-rules create bosh-ci-external --description "BOSH CI External traffic" --network ${google_network} --target-tags bosh-ci-external --allow tcp:22,tcp:443,tcp:4222,tcp:6868,tcp:25250,tcp:25555,tcp:25777,udp:53
