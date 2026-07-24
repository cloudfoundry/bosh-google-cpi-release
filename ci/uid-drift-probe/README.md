# uid-drift bisection harness

Empirically locate where `non-root-user`'s home directory uid drifts (1001 → 1002 →
1003 → 1004 …) in the GCP CPI CI, causing `setup-director` to fail with
`mkdir: Permission denied`.

## What we already know (facts, not theory)

- The drift is **driver-independent** — observed on overlay (DIAG mountinfo shows
  `overlay`), so it is NOT concourse/registry-image-resource#404 (that bug is btrfs
  extraction only).
- It is **uid-only, gid-preserved** (`Uid: 1004 / Gid: 1001(non-root-user)`), which
  rules out a user-namespace remap (those move both uid and gid).
- It **compounds ~+1 per run over days**, while the image only rebuilds weekly — so
  something re-owns the *volume on the worker*, not the built image.
- Deployed Concourse is **8.2.4 / worker 2.5**, `cache_streamed_volumes: true`.
- baggageclaim's uid mapper in 8.2.4 is **identity for uid 1001** — so the documented
  translation path does not explain it. That is exactly why we reproduce instead of
  reason further.
- Pinning `non-root-user` to uid **2000** avoids it (verified). This harness is about
  root cause, not the fix.

## Hypothesis under test

The drift is introduced Concourse-side when the image is consumed as a **volume**
(`image: <get-resource>`) and re-streamed/cached between workers
(`cache_streamed_volumes: true`), and it compounds each time the cached volume is
re-owned. The `stat-inline` job (inline `image_resource:`) is the control that should
stay clean.

## Files

| file | role |
|------|------|
| `Dockerfile`      | probe image: `useradd non-root-user` with NO uid pin → lands at 1001, like the real image before the fix. Bakes build-time state into `BAKED_STATE.txt`. |
| `probe.sh`        | dumps home ownership, passwd, worker identity, writable check, and a machine-readable `DRIFT:` line. |
| `local-bisect.sh` | proves the image is clean through build + `docker save` + registry round-trip, BEFORE Concourse. |
| `pipeline.yml`    | `build-probe` → `stat-inline` (control) + `stat-via-get` (suspect), all re-triggerable. |

## Procedure

### Step 1 — local control (rules build + registry in or out)

```bash
./local-bisect.sh <youruser>/uid-drift-probe:latest
```

Look at each hop's `home_disk_uid=… passwd_uid=…`. Expected: `1001 == 1001` at every
hop. If it's already drifted here → bug is build/registry (stop, that's the answer).
If clean → continue to Concourse.

### Step 2 — set up the probe pipeline

Push these files to a branch, then:

```bash
fly -t bosh set-pipeline -p uid-drift-probe \
  -c ci/uid-drift-probe/pipeline.yml \
  -v probe_branch=<your-branch> \
  -v probe_image_repository=<youruser>/uid-drift-probe
fly -t bosh unpause-pipeline -p uid-drift-probe
```

### Step 3 — build once, then bisect the two consumption modes

```bash
fly -t bosh trigger-job -j uid-drift-probe/build-probe -w
```

Then run BOTH and compare the `DRIFT:` line:

```bash
fly -t bosh trigger-job -j uid-drift-probe/stat-inline  -w   # control
fly -t bosh trigger-job -j uid-drift-probe/stat-via-get -w   # suspect
```

- `stat-inline` `delta=0`, `stat-via-get` `delta>0` → confirmed: the volume/`image:`
  path is the culprit.
- both `delta=0` → not reproduced on first hit; go to Step 4 (it compounds).

### Step 4 — watch it compound (the key step)

Re-trigger `stat-via-get` repeatedly, noting the `selected worker` and `DRIFT:` each
time. The real failures grew +1 over days, so force volume streaming between workers:

```bash
for i in 1 2 3 4 5 6; do
  echo "=== run $i ==="
  fly -t bosh trigger-job -j uid-drift-probe/stat-via-get -w 2>&1 | grep -E 'selected worker|DRIFT:|WRITABLE'
done
```

Watch whether `delta` climbs and whether it correlates with the volume being streamed
`from` a different worker (the build log line `streaming volume … from <uuid>`).

### Step 5 — pin the hop

Once the drift appears, capture the full `probe.sh` output for that run (worker uuid,
mount driver, uid scan). That output — a clean image that drifts only after N Concourse
volume hops — is the root-cause evidence for the PR.

## Cleanup

```bash
fly -t bosh destroy-pipeline -p uid-drift-probe
```
