#!/usr/bin/env bash

set -eu

script_dir="$( cd "$( dirname "$0" )" && pwd )"

fly -t bosh-ecosystem set-pipeline \
    -p bosh-google-cpi \
    -c ${script_dir}/pipeline.yml
