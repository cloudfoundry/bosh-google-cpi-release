---
google_raw_stemcells_bucket_name: <AN EXISTING GCS BUCKET TO STORE THE RAW MACHINE IMAGES>
google_raw_stemcells_json_key_data: |
  <THE CONTENT OF YOUR GOOGLE ACCOUNT JSON KEY FILE>
google_light_stemcells_bucket_name: <AN EXISTING S3 or GCS BUCKET TO STORE THE LIGHT STEMCELLS>
google_light_stemcells_access_key_id: <ACCESS KEY FOR STEMCELL BUCKET>
google_light_stemcells_secret_access_key: <SECRET KEY FOR STEMCELL BUCKET>
google_light_stemcells_endpoint: <ENTER s3.amazonaws.com FOR AWS, storage.googleapis.com for GCS>
google_light_stemcells_region: <REGION CONTAINING THE STEMCELL BUCKET>
google_boshio_checksum_token: "" # <SET TO A VALID TOKEN TO POST STEMCELL CHECKSUMS TO BOSH.IO, LEAVE EMPTY TO SKIP>
ssh_private_key: |
  <GCE SSH KEY, USED TO VERIFY STEMCELL BOOTS>
gce_credentials_json: |
  <GCE SERVICE ACCOUNT KEY WITH PERMISSION TO CREATE NETWORKS + INSTANCES>
gce_project_id: <GCE PROJECT ID USED TO APPLY TERRAFORM TEMPLATE>
terraform_bucket_name: <GCS BUCKET TO STORE TERRAFORM STATE FILES>
terraform_bucket_access_key: <GCS S3-COMPATIBLE KEY TO ACCESS BUCKET>
terraform_bucket_secret_key: <GCS S3-COMPATIBLE SECRET TO ACCESS BUCKET>
