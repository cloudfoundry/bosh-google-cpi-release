#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby.sh
chruby "ruby-2.1.7"

check_param build_number
check_param os_name
check_param os_version
check_param os_image

# Set the proper permissions
sudo chown -R ubuntu:ubuntu bosh-src

pushd bosh-src
  echo "Installing gems..."
  bundle install --local

  echo "Creating stemcell..."
  CANDIDATE_BUILD_NUMBER=${build_number} bundle exec rake stemcell:build[google,kvm,${os_name},${os_version},go,bosh-os-images,${os_image}]
popd

echo "Copying stemcell..."
mv /mnt/stemcells/google/kvm/${os_name}/work/work/bosh-stemcell-${build_number}-google-kvm-${os_name}-${os_version}-go_agent.tgz stemcell/
