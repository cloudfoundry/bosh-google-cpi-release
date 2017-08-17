---
version_semver_access_key: # GCS interop access key
version_semver_secret_key: # GCS interop secret key
version_semver_bucket_name: # GCS bucket for semver storage
version_semver_region: # GCS bucket's region
release_blobs_access_key: # GCS interop access key for release blobs
release_blobs_secret_key: # GCS interop secret key key for release blobs
google_stemcells_bucket_name: # GCS bucket that contains stsemcells
google_releases_bucket_name: # GCS bucket that releases are stored in
google_test_bucket_name: # GCS bucket used by the director in tests
github_deployment_key_bosh_google_cpi_release: |
  # GitHub deployment key for release artifacts
github_pr_access_token: # An access token with repo:status access, used to test PRs

# Google Cloud Platform configuration
google_project: # Google project ID
google_region: # Default Google Compute Engine region
google_zone: # Default Google Compute Engine zone (must be in {{google_region}}
google_json_key_data: |
  # Google Compute Engine Service Account JSON (created in {{google_project}}


# The following configuration values are names of resources that will be
# automatically created and destroyed in the pipeline. They must not conflict
# with existing resources in {{google_project}}
# Name of a service account that will be used in integration tests
google_service_account: google-cpi-ci-service-account
# Name of an auto-configured network in {{google_project}}
google_auto_network: google-cpi-ci-auto-network
# Name of a manually-configured network in {{google_project}}
google_network: google-cpi-ci-network
# Name of a manually-configured subnetwork in {{google_network}}
google_subnetwork: google-cpi-ci-subnetwork
# Name of firewall for internal access
google_firewall_internal: google-cpi-ci-firewall-internal
# Name of firewall for external access
google_firewall_external: google-cpi-ci-firewall-external
# Name of an external IP address used in integration tests
google_address_int_ubuntu: google-cpi-ci-ip-int-ubuntu
# Name of an external IP address used in BATS tests
google_address_bats_ubuntu: google-cpi-ci-ip-bats-ubuntu
# Name of an external IP address used to create a director
google_address_director_ubuntu: google-cpi-ci-ip-director-ubuntu
# Name of a network target pool
google_target_pool: google-cpi-ci-target-pool
# Name of a backend service
google_backend_service: google-cpi-ci-backend-service
# Name of a region backend service
google_region_backend_service: google-cpi-ci-regional-backend-service

# Networking configuration
# The CIDR range of {{google_subnetwork}}
google_subnetwork_range: 10.0.0.0/24
# The gateway IP of {{google_subnetwork}}
google_subnetwork_gw: 10.0.0.1

# All of the following IP addresses must be within {{google_subnetwork}}'s CIDR
# and be unique.
# Three comma-delimited IP address in {{google_subnetwork}}
google_address_static_int_ubuntu: 10.0.0.100,10.0.0.101,10.0.0.102
# A private IP address in {{google_subnetwork}}
google_address_static_director_ubuntu: 10.0.0.6
# A private IP address in {{google_subnetwork}}
google_address_static_bats_ubuntu: 10.0.0.20
# Two comma-delimited IP address in {{google_subnetwork}}
google_address_static_pair_bats_ubuntu: 10.0.0.20,10.0.0.21
# Hyphen-delimited range that contains {{google_address_static_pair_bats_ubuntu}} and {{google_address_static_bats_ubuntu}}
google_address_static_bats_available_range_ubuntu: 10.0.0.20-10.0.0.30

# SSH and auth information
private_key_user: vcap
private_key_data: |
  # Contents of a private key whose public key component is set as a project-wide SSH
  # key in {{google_project}}
bat_vcap_password: # A password to use for bats

# Do not change
director_username: admin
director_password: admin

