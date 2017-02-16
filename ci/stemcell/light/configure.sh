#!/bin/bash

fly -t cpi set-pipeline -p light-gce-stemcells \
  -c ./ci/stemcell/light/pipeline.yml \
  -l <(lpass show --note "google stemcell concourse secrets")