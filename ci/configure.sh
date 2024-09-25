#!/usr/bin/env bash
set -eu

REPO_ROOT="$( cd "$( dirname "$0" )/.." && pwd )"

fly -t "${CONCOURSE_TARGET:-bosh}" \
    set-pipeline -p bosh-google-cpi \
    -c "${REPO_ROOT}/ci/pipeline.yml"
