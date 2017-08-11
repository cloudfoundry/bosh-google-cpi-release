#!/bin/bash

set -eu

: ${BUCKET_NAME:?}
: ${BOSH_IO_BUCKET_NAME:?} # used to check if current stemcell already exists

# inputs
stemcell_dir="$PWD/stemcell"

# outputs
light_stemcell_dir="$PWD/light-stemcell"
raw_stemcell_dir="$PWD/raw-stemcell"

echo "Creating light stemcell..."

original_stemcell="$(echo ${stemcell_dir}/*.tgz)"
original_stemcell_name="$(basename "${original_stemcell}")"
raw_stemcell_name="$(basename "${original_stemcell}" .tgz)-raw.tar.gz"
light_stemcell_name="light-${original_stemcell_name}"


bosh_io_light_stemcell_url="https://s3.amazonaws.com/$BOSH_IO_BUCKET_NAME/$light_stemcell_name"
wget --spider "$bosh_io_light_stemcell_url"
if [[ "$?" != 0 ]]; then
  echo "Google light stemcell '$light_stemcell_name' already exists!"
  echo "You can download here: $bosh_io_light_stemcell_url"
  exit 1
fi

mkdir working_dir
pushd working_dir
  tar xvf "${original_stemcell}"

  raw_stemcell_path="${raw_stemcell_dir}/${raw_stemcell_name}"
  mv image "${raw_stemcell_path}"
  raw_disk_sha1="$(sha1sum ${raw_stemcell_path} | awk '{print $1}')"
  echo -n "${raw_disk_sha1}" > ${raw_stemcell_path}.sha1

  > image
  light_stemcell_sha1=$(sha1sum image | awk '{print $1}')
  sed -i '/^sha1: .*/c\sha1: '${light_stemcell_sha1}'' stemcell.MF
  echo "  source_url: https://storage.googleapis.com/${BUCKET_NAME}/${raw_stemcell_name}" >> stemcell.MF
  echo "  raw_disk_sha1: ${raw_disk_sha1}" >> stemcell.MF

  light_stemcell_path="${light_stemcell_dir}/${light_stemcell_name}"
  tar czvf "${light_stemcell_path}" *
popd
