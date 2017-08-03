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
google_service_account: # Name of a service account that will be used in integration tests
google_auto_network: # Name of an auto-configured network in {{google_project}
google_network: # Name of a manually-configured network in {{google_project}
google_subnetwork: # Name of a manually-configured subnetwork in {{google_network}}
google_firewall_internal: # Name of firewall for internal access
google_firewall_external: # Name of firewall for internal access
google_address_int_ubuntu: # Name of an external IP address used in integration tests
google_address_bats_ubuntu: # Name of an external IP address used in BATS tests
google_address_director_ubuntu: # Name of an external IP address used to create a director
google_target_pool: # Name of a network target pool
google_backend_service: # Name of a backend service
google_region_backend_service: # Name of a region backend service

# Networking configuration
google_subnetwork_range: # The CIDR range of {{google_subnetwork}}
google_subnetwork_gw: # The gateway IP of {{google_subnetwork}}

# All of the following IP addresses must be within {{google_subnetwork}}'s CIDR
# and be unique.
google_address_static_int_ubuntu: # Three comma-delimited IP address in {{google_subnetwork}}
google_address_static_director_ubuntu: # A private IP address in {{google_subnetwork}}
google_address_static_bats_ubuntu: # A private IP address in {{google_subnetwork}}
google_address_static_pair_bats_ubuntu: # Two comma-delimited IP address in {{google_subnetwork}}
google_address_static_bats_available_range_ubuntu: # Hyphen-delimited range that contains {{google_address_static_pair_bats_ubuntu}} and {{google_address_static_bats_ubuntu}}

# SSH and auth information
private_key_user: vcap
private_key_data: |
  # Contents of a private key whose public key component is set as a project-wide SSH
  # key in {{google_project}}
bat_vcap_password: # A password to use for bats

# Do not change
director_username: admin
director_password: admin

