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
# Name of an auto-configured network in {{google_project}}
google_auto_network: google-cpi-ci-auto-network
# Name of a manually-configured network in {{google_project}}
google_network: google-cpi-ci-network

# Networking configuration
# The CIDR range of {{google_subnetwork}}
google_subnetwork_range: 10.0.0.0/24
# The gateway IP of {{google_subnetwork}}
google_subnetwork_gw: 10.0.0.1

# All of the following IP addresses must be within {{google_subnetwork}}'s CIDR
# and be unique.
# Three comma-delimited IP address in {{google_subnetwork}}
google_address_static_int: 10.0.0.100,10.0.0.101,10.0.0.102
# A private IP address in {{google_subnetwork}}
google_address_static_director: 10.0.0.6
# A private IP address in {{google_subnetwork}}
google_address_static_bats: 10.0.0.20
# Two comma-delimited IP address in {{google_subnetwork}}
google_address_static_pair_bats: 10.0.0.20,10.0.0.21
# Hyphen-delimited range that contains {{google_address_static_pair_bats}} and {{google_address_static_bats}}
google_address_static_bats_available_range: 10.0.0.20-10.0.0.30

# SSH and auth information
private_key_user: vcap
private_key_data: |
  # Contents of a private key whose public key component is set as a project-wide SSH
  # key in {{google_project}}
