#!/usr/bin/env bash

set -e

until lpass status;do
    LPASS_DISABLE_PINENTRY=1 lpass ls a
    sleep 1
done

until fly -t cpi status;do
    fly -t cpi login
    sleep 1
done

fly -t cpi set-pipeline \
    -p bosh-google-cpi \
    -c pipeline.yml \
    -l <(lpass show --notes "google-cpi concourse secrets")
