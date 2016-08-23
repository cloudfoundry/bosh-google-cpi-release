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
