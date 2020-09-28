#!/usr/bin/env bash

check_param() {
  local name=$1
  local value=$(eval echo '$'$name)
  if [ "$value" == 'replace-me' ] || [ "$value" == '' ]; then
    echo "environment variable $name must be set"
    exit 1
  fi
}

print_git_state() {
  echo "--> last commit..."
  TERM=xterm-256color git log -1
  echo "---"
  echo "--> local changes (e.g., from 'fly execute')..."
  TERM=xterm-256color git status --verbose
  echo "---"
}

declare -a on_exit_items
on_exit_items=()

function on_exit {
  echo "Running ${#on_exit_items[@]} on_exit items..."
  for i in "${on_exit_items[@]}"
  do
    for try in $(seq 0 9); do
      sleep $try
      echo "Running cleanup command $i (try: ${try})"
        eval $i || continue
      break
    done
  done
}

function add_on_exit {
  local n=${#on_exit_items[@]}
  on_exit_items=("${on_exit_items[@]}" "$*")
  if [[ $n -eq 0 ]]; then
    trap on_exit EXIT
  fi
}


function check_go_version {
  local cpi_release=$1
  local release_go_version="$(cat $cpi_release/packages/golang-*/spec.lock | \
    grep linux | awk '{print $2}' | sed "s/golang-\(.*\)-linux/\1/")"

  local current=$(go version)
  if [[ "$current" != *"$release_go_version"* ]]; then
    echo "Go version is incorrect. Required version: $release_go_version"
    exit 1
  fi
}

function read_infrastructure {
  echo "Reading infrastructure values..."
  export google_project=$(cat ${infrastructure_metadata} | jq -r .google_project)
  export google_region=$(cat ${infrastructure_metadata} | jq -r .google_region)
  export google_zone=$(cat ${infrastructure_metadata} | jq -r .google_zone)
  export google_json_key_data=$(cat ${infrastructure_metadata} | jq -r .google_json_key_data)
  export google_auto_network=$(cat ${infrastructure_metadata} | jq -r .google_auto_network)
  export google_network=$(cat ${infrastructure_metadata} | jq -r .google_network)
  export google_subnetwork=$(cat ${infrastructure_metadata} | jq -r .google_subnetwork)
  export google_subnetwork_gateway=$(cat ${infrastructure_metadata} | jq -r .google_subnetwork_gateway)
  export google_firewall_internal=$(cat ${infrastructure_metadata} | jq -r .google_firewall_internal)
  export google_firewall_external=$(cat ${infrastructure_metadata} | jq -r .google_firewall_external)
  export google_backend_service=$(cat ${infrastructure_metadata} | jq -r .google_backend_service)
  export google_region_backend_service=$(cat ${infrastructure_metadata} | jq -r .google_region_backend_service)
  export google_target_pool=$(cat ${infrastructure_metadata} | jq -r .google_target_pool)
  export google_address_director_ip=$(cat ${infrastructure_metadata} | jq -r .google_address_director_ip)
  export google_address_director_internal_ip=$(cat ${infrastructure_metadata} | jq -r .google_address_director_internal_ip)
  export google_address_bats_ip=$(cat ${infrastructure_metadata} | jq -r .google_address_bats_ip)
  export google_address_bats_internal_ip=$(cat ${infrastructure_metadata} | jq -r .google_address_bats_internal_ip)
  export google_address_bats_internal_ip_pair=$(cat ${infrastructure_metadata} | jq -r .google_address_bats_internal_ip_pair)
  export google_address_bats_internal_ip_static_range=$(cat ${infrastructure_metadata} | jq -r .google_address_bats_internal_ip_static_range)

  export google_address_int_ip=$(cat ${infrastructure_metadata} | jq -r .google_address_int_ip)
  export google_address_int_internal_ip=$(cat ${infrastructure_metadata} | jq -r .google_address_int_internal_ip)
  export google_node_group=$(cat  ${infrastructure_metadata} | jq -r .google_node_group)
  export google_service_account=$(cat ${infrastructure_metadata} | jq -r .google_service_account)
}
