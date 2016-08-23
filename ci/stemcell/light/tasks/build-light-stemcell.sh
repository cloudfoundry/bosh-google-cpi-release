#!/bin/bash

set -e

: ${BUCKET_NAME:?}
: ${BOSHIO_TOKEN:=""}

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
  echo -n $(sha1sum ${raw_stemcell_path} | awk '{print $1}') > ${raw_stemcell_path}.sha1

  > image
  echo "  source_url: https://storage.googleapis.com/${BUCKET_NAME}/${raw_stemcell_name}" >> stemcell.MF

  light_stemcell_path="${light_stemcell_dir}/${light_stemcell_name}"
  tar czvf "${light_stemcell_path}" *

  checksum="$(sha1sum ${light_stemcell_path} | awk '{print $1}')"
  echo -n "${checksum}" > ${light_stemcell_path}.sha1

  if [ -n "${BOSHIO_TOKEN}" ]; then
    curl -X POST \
      --fail \
      -d "sha1=${checksum}" \
      -H "Authorization: bearer ${BOSHIO_TOKEN}" \
      "https://bosh.io/checksums/${light_stemcell_name}"
  fi
popd
