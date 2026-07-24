FROM ubuntu

# Reproduce the real bosh-google-cpi-boshrelease user setup WITHOUT the uid 2000 pin.
# On noble the `ubuntu` user pre-exists at uid 1000, so this useradd lands non-root-user
# at uid 1001 -- the value that drifts in the failing setup-director runs.
RUN useradd non-root-user --shell /bin/bash --create-home \
    && echo "non-root-user ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# Bake the expected state into the image so we can compare against it after each hop.
RUN { \
      echo "=== BAKED-IN IMAGE STATE (build time) ==="; \
      echo "passwd non-root-user: $(getent passwd non-root-user)"; \
      echo "passwd ubuntu:        $(getent passwd ubuntu)"; \
      stat -c 'home owner uid=%u gid=%g mode=%a' /home/non-root-user; \
    } > /home/non-root-user/BAKED_STATE.txt \
    && cat /home/non-root-user/BAKED_STATE.txt

COPY probe.sh /probe.sh
RUN chmod +x /probe.sh
