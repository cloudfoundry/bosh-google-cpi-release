#!/bin/bash

set -eu

: ${AWS_ACCESS_KEY_ID:?}
: ${AWS_SECRET_ACCESS_KEY:?}
: ${AWS_DEFAULT_REGION:?}
: ${AWS_ENDPOINT:?}
: ${OUTPUT_BUCKET:?}
: ${BOSHIO_TOKEN:=""}

# inputs
light_stemcell_dir="$PWD/light-stemcell"

light_stemcell_path="$(echo ${light_stemcell_dir}/*.tgz)"
light_stemcell_name="$(basename "${light_stemcell_path}")"

echo "Uploading light stemcell ${light_stemcell_name} to ${OUTPUT_BUCKET}..."
aws --endpoint-url=${AWS_ENDPOINT} s3 cp "${light_stemcell_path}" "s3://${OUTPUT_BUCKET}"

if [ -n "${BOSHIO_TOKEN}" ]; then
  echo "Publishing light stemcell checksum to bosh.io..."

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
