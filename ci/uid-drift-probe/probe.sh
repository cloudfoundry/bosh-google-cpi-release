#!/usr/bin/env bash
# uid-drift probe: dumps the on-disk ownership of non-root-user's home and the
# worker identity, so we can see WHERE and WHEN uid 1001 drifts.
set +e

echo "======================================================================"
echo "PROBE @ $(date -u +%FT%TZ)"
echo "======================================================================"

echo "--- worker identity ---"
echo "hostname: $(hostname)"
# On a Concourse worker container this maps back to the garden host; best-effort.
cat /etc/hostname 2>/dev/null

echo "--- process identity ---"
echo "whoami=$(whoami 2>/dev/null) id=$(id)"

echo "--- passwd entries ---"
echo "non-root-user: $(getent passwd non-root-user)"
echo "ubuntu:        $(getent passwd ubuntu)"

echo "--- baked-in image state (what build time recorded) ---"
cat /home/non-root-user/BAKED_STATE.txt 2>/dev/null || echo "(BAKED_STATE.txt missing!)"

echo "--- CURRENT on-disk state ---"
ls -lan /home/
echo
ls -lan /home/non-root-user/
echo
stat /home/non-root-user

# The single machine-readable line we grep across runs to watch the drift.
home_uid="$(stat -c '%u' /home/non-root-user)"
passwd_uid="$(getent passwd non-root-user | cut -d: -f3)"
echo "DRIFT: passwd_uid=${passwd_uid} home_disk_uid=${home_uid} delta=$(( home_uid - passwd_uid ))"

echo "--- writable check (the thing setup-director actually fails on) ---"
if sudo -u non-root-user bash -c 'mkdir -p /home/non-root-user/.config/gcloud 2>/dev/null && echo yes'; then
  echo "WRITABLE: yes"
else
  echo "WRITABLE: no"
fi

echo "--- who owns which uid (scan) ---"
for u in 1000 1001 1002 1003 1004 1005 2000; do
  hits="$(find / -xdev -uid "$u" 2>/dev/null | head -5 | tr '\n' ' ')"
  [ -n "$hits" ] && echo "uid $u: $hits"
done

echo "--- rootfs mount (driver check) ---"
grep -E ' / ' /proc/self/mountinfo 2>/dev/null | head -1

echo "======================================================================"
