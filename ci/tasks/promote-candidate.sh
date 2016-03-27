#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param aws_access_key_id_release
check_param aws_secret_access_key_release

cpi_release_name="bosh-google-cpi"
integer_version=`cut -d "." -f1 release-version-semver/number`
echo $integer_version > promoted/integer_version
echo "BOSH Google CPI BOSH Release v${integer_version}" > promoted/annotation_message

cp -r bosh-cpi-src promoted/repo

dev_release=$(echo $PWD/bosh-cpi-release/*.tgz)

pushd promoted/repo
  echo "Creating config/private.yml with blobstore secrets"
  set +x
  cat > config/private.yml << EOF
---
blobstore:
  s3:
    access_key_id: ${aws_access_key_id_release}
    secret_access_key: ${aws_secret_access_key_release}
EOF

  echo "Using BOSH CLI version..."
  bosh version

  echo "Finalizing CPI BOSH Release..."
  bosh finalize release ${dev_release} --version ${integer_version}

  rm config/private.yml

  git diff | cat
  git add .

  git config --global user.email cf-bosh-eng@pivotal.io
  git config --global user.name CI
  git commit -m "BOSH Google CPI BOSH Release v${integer_version}"
popd

mv promoted/repo/releases/${cpi_release_name}/${cpi_release_name}-${integer_version}.tgz promoted/
