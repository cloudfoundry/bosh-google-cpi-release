#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param aws_access_key_id
check_param aws_secret_access_key

# Creates an integer version number from the semantic version format
# May be changed when we decide to fully use semantic versions for releases
integer_version=`cut -d "." -f1 release-version-semver/number`
echo $integer_version > promoted/integer_version

cp -r bosh-cpi-src promoted/repo

dev_release=$(echo $PWD/bosh-cpi-release/*.tgz)

pushd promoted/repo
  echo "Creating config/private.yml with blobstore secrets"
  set +x
  cat > config/private.yml << EOF
---
blobstore:
  s3:
    access_key_id: ${aws_access_key_id}
    secret_access_key: ${aws_secret_access_key}
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
  git commit -m "New Final CPI BOSH Release v${integer_version}"
popd
