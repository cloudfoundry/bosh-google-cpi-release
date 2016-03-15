#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby.sh
chruby "ruby-2.1.7"

check_param build_number
check_param os_name
check_param os_version
check_param os_image_bucket
check_param os_image_file

TASK_DIR=$PWD

# This is copied from https://github.com/concourse/concourse/blob/3c070db8231294e4fd51b5e5c95700c7c8519a27/jobs/baggageclaim/templates/baggageclaim_ctl.erb#L23-L54
# helps the /dev/mapper/control issue and lets us actually do scary things with the /dev mounts
# This allows us to create device maps from partition tables in image_create/apply.sh
function permit_device_control() {
  local devices_mount_info=$(cat /proc/self/cgroup | grep devices)

  local devices_subsytems=$(echo $devices_mount_info | cut -d: -f2)
  local devices_subdir=$(echo $devices_mount_info | cut -d: -f3)

  cgroup_dir=/mnt/tmp-todo-devices-cgroup

  if [ ! -e ${cgroup_dir} ]; then
    # mount our container's devices subsystem somewhere
    mkdir ${cgroup_dir}
  fi

  if ! mountpoint -q ${cgroup_dir}; then
    mount -t cgroup -o $devices_subsytems none ${cgroup_dir}
  fi

  # permit our cgroup to do everything with all devices
  # ignore failure in case something has already done this; echo appears to
  # return EINVAL, possibly because devices this affects are already in use
  echo a > ${cgroup_dir}${devices_subdir}/devices.allow || true
}

# Also copied from baggageclaim_ctl.erb creates 64 loopback mappings. This fixes failures with losetup --show --find ${disk_image}
function create_loopback_mappings() {
  for i in $(seq 0 64); do
    if ! mknod -m 0660 /dev/loop$i b 7 $i; then
      break
    fi
  done
}

permit_device_control
create_loopback_mappings

chown -R ubuntu:ubuntu bosh-src
sudo --preserve-env --set-home --user ubuntu -- /bin/bash --login -i <<SUDO
  pushd bosh-src
    echo "Installing gems..."
    bundle install --local

    echo "Creating stemcell..."
    CANDIDATE_BUILD_NUMBER=${build_number} bundle exec rake stemcell:build[google,kvm,${os_name},${os_version},go,${os_image_bucket},${os_image_file}]
  popd
SUDO

echo "Copying stemcell..."
mv /mnt/stemcells/google/kvm/${os_name}/work/work/bosh-stemcell-${build_number}-google-kvm-${os_name}-${os_version}-go_agent.tgz stemcell/
mv /mnt/stemcells/google/kvm/${os_name}/work/work/stemcell/image stemcell/bosh-stemcell-${build_number}-google-kvm-${os_name}-${os_version}-go_agent-raw.tar.gz
