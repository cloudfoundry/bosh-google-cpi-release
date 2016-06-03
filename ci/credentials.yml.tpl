---
version_semver_access_key: <YOUR GCS S3 INTER-OP ACCESS KEY ID>
version_semver_secret_key: <YOUR GCS S3 INTER-OP SECRET ACCESS KEY>
version_semver_bucket_name: <AN EXISTING GCS S3 INTER-OP BUCKET TO STORE THE ARTIFACTS>
version_semver_region: <THE GCS S3 INTER-OP BUCKET REGION>
release_blobs_access_key: <THE GCS S3 INTER-OP ACCESS KEY ID FOR THE RELEASE BLOBSTORE>
release_blobs_secret_key: <THE GCS S3 INTER-OP SECRET ACCESS KEY FOR THE RELEASE BLOBSTORE>
github_deployment_key_bosh_google_cpi_release: |
  <THE CONTENT OF THE GITHUB DEPLOY PRIVATE SSH KEY FILE FOR THE BOSH GOOGLE CPI RELEASE REPOSITORY>
google_project: <YOUR GOOGLE PROJECT>
google_region: <THE GOOGLE REGION TO USE>
google_zone: <THE GOOGLE ZONE TO USE>
google_json_key_data: |
  <THE CONTENT OF YOUR GOOGLE ACCOUNT JSON KEY FILE>
google_network: <NAME OF THE GOOGLE NETWORK TO BE CREATED>
google_subnetwork: <NAME OF THE GOOGLE SUBNETWORK NETWORK TO BE CREATED>
google_subnetwork_range: <NAME OF THE GOOGLE SUBNETWORK RANGE NETWORK TO BE CREATED>
google_subnetwork_gw: <NAME OF THE GOOGLE SUBNETWORK GATEWAY TO BE CREATED>
google_firewall_internal: <NAME OF THE GOOGLE FIREWALL (INTERNAL TRAFFIC) TO BE CREATED>
google_firewall_external: <NAME OF THE GOOGLE FIREWALL (EXTERNAL TRAFFIC) TO BE CREATED>
google_address_director_ubuntu: <NAME OF THE GOOGLE ADDRESS (DIRECTOR IP) TO BE CREATED FOR UBUNTU TESTS>
google_address_static_director_ubuntu: <NAME OF THE STATIC DIRECTOR IP TO BE CREATED FOR UBUNTU TESTS>
google_address_bats_ubuntu: <NAME OF THE GOOGLE ADDRESS (BATS IP) TO BE CREATED FOR UBUNTU TESTS>
google_address_static_bats_ubuntu: <NAME OF THE GOOGLE ADDRESS (BATS IP) TO BE CREATED FOR UBUNTU TESTS>
google_address_static_bats_range_ubuntu: <HYPHEN-DELIMITED RANGE OF STATIC IPS BATS DIRECTOR CAN USE>
google_address_static_bats_dir_ubuntu: <STATIC IP THAT WILL BE ASSIGNED TO BATS DIR>
google_stemcells_bucket_name: <NAME OF THE GCS BUCKET TO STORE STEMCELLS>
google_releases_bucket_name: <NAME OF THE GCS BUCKET TO STORE RELEASES>
private_key_user: vcap
private_key_data: |
  <THE CONTENT OF YOUR PRIVATE SSH KEY FILE>
director_username: <THE DIRECTOR USERNAME TO BE CREATED>
director_password: <THE DIRECTOR PASSWORD TO BE CREATED>
bat_vcap_password: <THE BAT VCAP PASSWORD TO BE CREATED>

google_existing_stemcell: ci
google_target_pool: ci
google_backend_service: ci

