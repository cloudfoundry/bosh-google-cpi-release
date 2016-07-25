#!/bin/bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
check_param IAAS
check_param HYPERVISOR
check_param OS_NAME
check_param OS_VERSION
check_param CANDIDATE_BUILD_NUMBER

export TASK_DIR=$PWD

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

permit_device_control

# Also copied from baggageclaim_ctl.erb creates 64 loopback mappings. This fixes failures with losetup --show --find ${disk_image}
for i in $(seq 0 64); do
  if ! mknod -m 0660 /dev/loop$i b 7 $i; then
    break
  fi
done

chown -R ubuntu:ubuntu bosh-src
sudo --preserve-env --set-home --user ubuntu -- /bin/bash --login -i <<SUDO
  pushd bosh-src
    source /etc/profile.d/chruby.sh
    chruby "ruby-2.1.7"
    cd bosh-src

    bundle install --local
    CANDIDATE_BUILD_NUMBER=${CANDIDATE_BUILD_NUMBER} bundle exec rake   stemcell:build[$IAAS,$HYPERVISOR,$OS_NAME,$OS_VERSION,go,bosh-os-images,bosh-ubuntu-trusty-os-image.tgz]
  popd
SUDO

image_name=bosh-stemcell-${CANDIDATE_BUILD_NUMBER}-google-kvm-${OS_NAME}-${OS_VERSION}-go_agent.tgz
image_path=stemcell/$image_name
mv /mnt/stemcells/google/kvm/${OS_NAME}/work/work/$image_name $image_path

echo -n $(sha1sum $image_path | awk '{print $1}') > $image_path.sha1
