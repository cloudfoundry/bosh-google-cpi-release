#!/usr/bin/env bash
set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_project
check_param google_region
check_param google_zone
check_param google_json_key_data
check_param google_network
check_param google_auto_network
check_param google_subnetwork
check_param google_subnetwork_range
check_param google_firewall_internal
check_param google_firewall_external
check_param google_address_director_ubuntu
check_param google_address_bats_ubuntu
check_param google_target_pool
check_param google_backend_service 
check_param google_region_backend_service
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

echo "Setting up google infrastructure..."
gcloud -q iam service-accounts create ${google_service_account}
gcloud -q compute addresses create ${google_address_director_ubuntu} --region ${google_region}
gcloud -q compute addresses create ${google_address_bats_ubuntu} --region ${google_region}
gcloud -q compute addresses create ${google_address_int_ubuntu} --region ${google_region}
gcloud -q compute networks create ${google_auto_network}
gcloud -q compute networks create ${google_network} --mode custom
gcloud -q compute networks subnets create ${google_subnetwork} --network=${google_network} --range=${google_subnetwork_range}
gcloud -q compute firewall-rules create ${google_firewall_internal} --description "BOSH CI Internal traffic" --network ${google_network} --source-tags ${google_firewall_internal} --target-tags ${google_firewall_internal} --allow tcp,udp,icmp
gcloud -q compute firewall-rules create ${google_firewall_external} --description "BOSH CI External traffic" --network ${google_network} --target-tags ${google_firewall_external} --allow tcp:22,tcp:443,tcp:4222,tcp:6868,tcp:25250,tcp:25555,tcp:25777,udp:53

# Target pool
gcloud -q compute target-pools create ${google_target_pool} --region=${google_region}

# Backend service
gcloud -q compute instance-groups unmanaged create ${google_backend_service} --zone ${google_zone}
gcloud -q compute http-health-checks create ${google_backend_service}
gcloud -q compute backend-services create ${google_backend_service} --http-health-checks ${google_backend_service} --port-name "http" --timeout "30"
gcloud -q compute backend-services add-backend ${google_backend_service} --instance-group ${google_backend_service} --zone ${google_zone} --balancing-mode "UTILIZATION" --capacity-scaler "1" --max-utilization "0.8"

# Region Backend service
gcloud -q compute instance-groups unmanaged create ${google_region_backend_service} --zone ${google_zone}
gcloud -q compute health-checks create tcp ${google_region_backend_service}

# This is a hack required to give the instance group a network association
gcloud -q compute instances create ${google_region_backend_service} --zone ${google_zone} --network ${google_network} --subnet ${google_subnetwork} --machine-type f1-micro

gcloud -q compute instance-groups unmanaged add-instances ${google_region_backend_service} --instances ${google_region_backend_service} --zone ${google_zone}
gcloud -q compute backend-services create ${google_region_backend_service} --region ${google_region} --health-checks ${google_region_backend_service} --protocol "TCP" --load-balancing-scheme "INTERNAL" --timeout "30"
gcloud -q compute backend-services add-backend ${google_region_backend_service} --instance-group ${google_region_backend_service} --zone ${google_zone} --region ${google_region}
gcloud -q compute instances delete ${google_region_backend_service} --zone ${google_zone}
