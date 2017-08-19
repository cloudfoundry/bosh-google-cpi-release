#!/bin/bash

fly -t cpi sp -p light-gce-stemcells \
  -c pipeline.yml \
  -l <(lpass show --note "google stemcell concourse secrets")
