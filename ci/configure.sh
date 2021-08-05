#!/usr/bin/env bash

set -eu

script_dir="$( cd "$( dirname "$0" )" && pwd )"

until lpass status;do
    LPASS_DISABLE_PINENTRY=1 lpass ls a
    sleep 1
done

fly -t bosh-ecosystem set-pipeline \
    -p bosh-google-cpi \
    -c ${script_dir}/pipeline.yml \
    -v dockerhub_password=$(lpass show "Docker Hub" --password) \
    -l <(lpass show --notes "google-cpi concourse secrets")
