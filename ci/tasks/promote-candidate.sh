#!/usr/bin/env bash

set -e

source ci/ci/tasks/utils.sh

check_param release_blobs_json_key

# Version info
semver_version="$(cat release-version-semver/number)"
echo "$semver_version" > promoted/semver_version
echo "v$semver_version" > promoted/prefixed_semver_version
echo "BOSH Google CPI BOSH Release v${semver_version}" > promoted/annotation_message

cp -r bosh-cpi-src promoted/repo

dev_release=$(echo $PWD/bosh-cpi-release/*.tgz)

pushd promoted/repo
  echo "Creating config/private.yml with blobstore secrets"
  set +x
  cat > config/private.yml << EOF
---
blobstore:
  options:
    credentials_source: static
    json_key: '${release_blobs_json_key}'
EOF

  echo "Using BOSH CLI version..."
  bosh --version

  echo "Finalizing CPI BOSH Release..."
  bosh finalize-release ${dev_release} --version ${semver_version}

  rm config/private.yml

  git diff | cat
  git add .

  git config --global user.email cf-bosh-eng@pivotal.io
  git config --global user.name CI
  git commit -m "BOSH Google CPI BOSH Release v${semver_version}"
popd
