#!/usr/bin/env bash
# Local control: prove the probe image is CLEAN (non-root-user home owned by 1001)
# at every hop BEFORE Concourse touches it. If the uid is already drifted here,
# the bug is in build/registry, not Concourse streaming. If it's clean here but
# drifts in the pipeline, the bug is Concourse-side (baggageclaim / cache_streamed_volumes).
#
# Usage: ./local-bisect.sh <registry/repo:tag>   e.g. ./local-bisect.sh youruser/uid-drift-probe:latest
set -euo pipefail

REF="${1:?usage: ./local-bisect.sh <registry/repo:tag>}"
HERE="$(cd "$(dirname "$0")" && pwd)"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

hr(){ printf '\n========== %s ==========\n' "$1"; }

hr "HOP A: docker build (local daemon)"
docker build -t "$REF" "$HERE"
echo "--- stat inside freshly built image ---"
docker run --rm "$REF" bash -c 'echo "home_disk_uid=$(stat -c %u /home/non-root-user) passwd_uid=$(getent passwd non-root-user | cut -d: -f3)"'

hr "HOP B: docker save -> tar header inspection (what the layer ACTUALLY encodes)"
docker save "$REF" -o "$WORK/img.tar"
mkdir -p "$WORK/x"; tar -xf "$WORK/img.tar" -C "$WORK/x"
echo "--- searching layer tarballs for the home/non-root-user entry ownership ---"
found=0
# OCI/docker-archive: layers are blobs; scan each tar-like blob for the path.
while IFS= read -r layer; do
  if tar -tvf "$layer" 2>/dev/null | grep -q 'home/non-root-user'; then
    echo ">>> layer: $layer"
    tar -tvf "$layer" 2>/dev/null | grep 'home/non-root-user' | head -20
    found=1
  fi
done < <(find "$WORK/x" -type f \( -name '*.tar' -o -path '*/blobs/*' \))
[ "$found" = 1 ] || echo "(no home/non-root-user entries found in layers — check base image layer)"

hr "HOP C: registry round-trip (push, remove local, pull fresh, stat)"
docker push "$REF"
docker rmi "$REF" >/dev/null 2>&1 || true
docker pull "$REF"
echo "--- stat after registry round-trip ---"
docker run --rm "$REF" bash -c 'echo "home_disk_uid=$(stat -c %u /home/non-root-user) passwd_uid=$(getent passwd non-root-user | cut -d: -f3)"'

hr "LOCAL VERDICT"
echo "If all hops above show home_disk_uid == passwd_uid (1001 == 1001), the image is"
echo "CLEAN through build+registry. Any drift must then be introduced Concourse-side."
echo "Next: run the probe PIPELINE (see README) and re-trigger the get+stat job repeatedly."
