#!/usr/bin/env bash

check_param() {
  local name=$1
  local value=$(eval echo '$'$name)
  if [ "$value" == 'replace-me' ]; then
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
  local release_go_version="$(cat "$cpi_release/packages/golang/spec" | \
    grep linux | awk '{print $2}' | sed "s/golang\/go\(.*\)\.linux-amd64.tar.gz/\1/")"

  local current=$(go version)
  if [[ "$current" != *"$release_go_version"* ]]; then
    echo "Go version is incorrect. Required version: $release_go_version"
    exit 1
  fi
}
