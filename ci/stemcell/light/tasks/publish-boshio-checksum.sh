#!/bin/bash

set -eux

: ${BOSHIO_TOKEN:=""}

# inputs
light_stemcell_dir="$PWD/light-stemcell"

if [ -n "${BOSHIO_TOKEN}" ]; then
  light_stemcell_path="$(echo ${light_stemcell_dir}/*.tgz)"
  light_stemcell_name="$(basename "${light_stemcell_path}")"

  echo "Publishing light stemcell checksum to bosh.io"

  checksum="$(sha1sum ${light_stemcell_dir}/*.tgz | awk '{print $1}')"

  curl -X POST \
    --fail \
    -d "sha1=${checksum}" \
    -H "Authorization: bearer ${BOSHIO_TOKEN}" \
    "https://bosh.io/checksums/${light_stemcell_name}"

  echo "Successfully published checksum!"
else
  echo "Checksum not provided, skipping publish checksum."
fi
