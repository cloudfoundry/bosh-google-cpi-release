#!/usr/bin/env bash

set -eu

task_dir=$PWD
repo_output=$task_dir/bosh-cpi-src-out

git config --global user.email "ci@localhost"
git config --global user.name "CI Bot"

git clone bosh-cpi-src-in "$repo_output"

cd "$repo_output"

cat > config/private.yml << EOF
---
blobstore:
  options:
    credentials_source: static
    json_key: '${release_blobs_json_key}'
EOF

bosh vendor-package golang-1-linux "$task_dir/golang-release"
bosh vendor-package golang-1-darwin "$task_dir/golang-release"

if [ -z "$(git status --porcelain)" ]; then
  exit
fi

git add -A

git commit -m "Update golang packages to $(cat "$task_dir/golang-release/packages/golang-1-linux/version") from golang-release"
