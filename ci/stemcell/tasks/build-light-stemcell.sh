#!/bin/bash

set -e

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
  mv image "${raw_stemcell_dir}/${raw_stemcell_name}"

  > image
  echo "  source_url: https://storage.googleapis.com/${BUCKET_NAME}/${raw_stemcell_name}" >> stemcell.MF

  light_stemcell_path="${light_stemcell_dir}/${light_stemcell_name}"
  tar czvf "${light_stemcell_path}" *
  echo -n $(sha1sum ${light_stemcell_path} | awk '{print $1}') > ${light_stemcell_path}.sha1
popd
