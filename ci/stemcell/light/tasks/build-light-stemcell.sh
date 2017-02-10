#!/bin/bash

set -eu

: ${BUCKET_NAME:?}

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
