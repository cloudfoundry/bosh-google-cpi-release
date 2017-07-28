#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

check_param google_project
check_param google_region
check_param google_zone
check_param google_json_key_data
check_param google_network
check_param google_subnetwork
check_param google_subnetwork_range
check_param google_subnetwork_gw
check_param google_firewall_internal
check_param google_firewall_external
check_param google_address_director
check_param google_address_static_director
check_param private_key_user
check_param private_key_data
check_param director_password
check_param director_username

deployment_dir="${PWD}/deployment"
cpi_release_name=bosh-google-cpi
google_json_key=${deployment_dir}/google_key.json
private_key=${deployment_dir}/private_key.pem
manifest_filename="director-manifest.yml"

echo "Setting up artifacts..."
cp ./bosh-cpi-release/*.tgz ${deployment_dir}/${cpi_release_name}.tgz
cp ./bosh-release/*.tgz ${deployment_dir}/bosh-release.tgz
cp ./stemcell/*.tgz ${deployment_dir}/stemcell.tgz

echo "Creating google json key..."
echo "${google_json_key_data}" > ${google_json_key}
mkdir -p $HOME/.config/gcloud/
cp ${google_json_key} $HOME/.config/gcloud/application_default_credentials.json

echo "Configuring google account..."
gcloud auth activate-service-account --key-file $HOME/.config/gcloud/application_default_credentials.json
gcloud config set project ${google_project}
gcloud config set compute/region ${google_region}
gcloud config set compute/zone ${google_zone}

echo "Looking for director IP..."
director_ip=$(gcloud compute addresses describe ${google_address_director} --format json | jq -r '.address')

echo "Creating private key..."
echo "${private_key_data}" > ${private_key}
chmod go-r ${private_key}
eval $(ssh-agent)
ssh-add ${private_key}

echo "Creating ${manifest_filename}..."
cat > "${deployment_dir}/${manifest_filename}"<<EOF
---
name: bosh
releases:
  - name: bosh
    url: file://bosh-release.tgz
  - name: ${cpi_release_name}
    url: file://${cpi_release_name}.tgz

resource_pools:
  - name: vms
    network: private
    stemcell:
      url: file://stemcell.tgz
    cloud_properties:
      zone: ${google_zone}
      machine_type: n1-standard-2
      root_disk_size_gb: 40
      root_disk_type: pd-standard
      service_scopes:
        - compute
        - devstorage.full_control

disk_pools:
  - name: disks
    disk_size: 32_768
    cloud_properties:
      type: pd-standard

networks:
  - name: private
    type: manual
    subnets:
    - range: ${google_subnetwork_range}
      gateway: ${google_subnetwork_gw}
      cloud_properties:
        network_name: ${google_network}
        subnetwork_name: ${google_subnetwork}
        tags:
          - ${google_firewall_internal}
          - ${google_firewall_external}
  - name: public
    type: vip

jobs:
  - name: bosh
    instances: 1

    templates:
      - name: nats
        release: bosh
      - name: postgres-9.4
        release: bosh
      - name: powerdns
        release: bosh
      - name: blobstore
        release: bosh
      - name: director
        release: bosh
      - name: health_monitor
        release: bosh
      - name: registry
        release: bosh
      - name: google_cpi
        release: bosh-google-cpi

    resource_pool: vms
    persistent_disk_pool: disks

    networks:
      - name: private
        static_ips: [${google_address_static_director}]
        default:
          - dns
          - gateway
      - name: public
        static_ips:
          - ${director_ip}

    properties:
      nats:
        address: 127.0.0.1
        user: nats
        password: nats-password

      postgres: &db
        listen_address: 127.0.0.1
        host: 127.0.0.1
        user: postgres
        password: postgres-password
        database: bosh
        adapter: postgres

      dns:
        address: ${google_address_static_director}
        domain_name: microbosh
        db: *db
        recursor: 169.254.169.254

      registry:
        address: ${google_address_static_director}
        host: ${google_address_static_director}
        db: *db
        http:
          user: registry
          password: registry-password
          port: 25777
        username: registry
        password: registry-password
        port: 25777

      blobstore:
        address: ${google_address_static_director}
        port: 25250
        provider: dav
        director:
          user: director
          password: director-password
        agent:
          user: agent
          password: agent-password

      director:
        address: 127.0.0.1
        name: micro-google
        db: *db
        cpi_job: google_cpi
        ssl:
          key: |
            -----BEGIN RSA PRIVATE KEY-----
            MIIEpAIBAAKCAQEArKtH+JNmEv4osTZwhyHBLauJaQjmNAS5vYDCep6F9AUpW3kL
            jcDYGk+BwyFrpLa7ECkINbYB+5iBEkbBRK0sRMIz5rzUWU7Qv/HA609rF/ynbSUS
            zk1uv9fgTC2BVnb2f7L03H1wqmtghp3bJvJ/LnMqe3x+OVEpkr1A5xqXsyshDSU6
            tgfM+4tNyBBQudXCv3JyyKoCJ/cQpLCOh9nQ0yO4r1fqNPbfuYSf3iGRfuFg4Vgo
            vLKOqRhWxNykGxcB14/uG1jxN1vX+Yg1aZvU075Z00M6NDc0gaYadDBjxNGkjYKz
            icHI0/EoT6qJnR+tGBzT0gq24rcyz2LhLtP6fwIDAQABAoIBAFwbwnjHqFvZWLuv
            3rc3OmWya8qsBKEbJDoCxbvDdJGHb1hsac1kYeMnJoGBAnsLPx6PxOFiBgzAfZnS
            RKbt+f9z2VvsvxolARZjUBY2d1qEXIvMiwuiIsIT1oLMg4IsU7IrNJOqFr/SJ9un
            uZA9K7sLlE3rSyooMZUlf8nIVcQtDVIP4sK57PEZkVcscjlV4MRO5q0cRpOd+IFC
            RDzRMNlZXxLQadbZiGoLmEMp56S6Fkr597k4lh+ijV9xuqcC/2R7yZ/UKOxh+Z6v
            eQ+69oype6EeBtMhrVuo8t11Fh8p5eomNKEW940e2aDvuVInKOTaw6RirKZL7yY/
            tMKqHIECgYEA0BO+h1EVfWYqvA9qVp7jmFke9uW/+FUIpEqh061lx38UHLhpVmvW
            tadcPPGYsCUH9oRcDEqtM56+2OSOSf4mvydZzzMS7q/OS9TGobeUjzH61MrsGp7D
            fU2zx/3yJo3D6pR6eJWitXSFA//tCZbkLoiwTgkHtKT60xwKILDsGtMCgYEA1G/f
            db9rk5L8fp5NkVHV+P/ttNtG69TFy9hHoW+PtDD0Xal9xjq/oLr1o19WidczUMjj
            uLd++Z5DIrWMX8o9MHkuKuPWaDj5aP1wrZgbJMVDhHx+qoeuhxGSYHl3J2RHN06W
            03IacWydavZG20e5Avunvp1/i31ozA9T3h3XPiUCgYEAqJZGucZlffuIRmTLCLGl
            v6r9npdZqa/j15EseqA0JaX9uqNjnYS0KuwVnL82sgje4cot9juPB5LoGD1eV+8W
            n6wXZPyBq2g/4krcQOzH7hlVnJFpKMxXoa+SKUjEqJ4WDXsNm6PJd/GXUD1MZYef
            C2DuT9ubJa7CFsfSINiYA8cCgYBSeEPF0FQQ7ET9WrM+MQjiK2i6h03XC7jl08ar
            E0Y0a7TSD5R2OiReX3YwwDg2NscDG5ncAdBXU2s4tEYUgcyTXtffaqe3ujaI3aq6
            mYwgEDyP2EzMIvRMFzQ+I6lwL2u+OtIur+M4GTRba9RCGGvojo2mYDo9iqf+YAzs
            86S1yQKBgQC4BMHlMw9dd9tmbHZVSIWnoZOYsPiQ8Uuerd4oSh5/w6U3ZZsDNv0S
            Ysqz/bFu22Ov4xb4PlNf/e+7Yx7rFIpTEFpphLxmt99aeebcPh0ShOzySr0y5YD/
            fjauZEns0I511J9Unats0HX7CUGyVLVf0ZU9WatCGINYRAEs6/m8GA==
            -----END RSA PRIVATE KEY-----
          cert: |
            -----BEGIN CERTIFICATE-----
            MIIDQjCCAiqgAwIBAgIRAOciLjtHiiFIpTYuXpA8Mm4wDQYJKoZIhvcNAQELBQAw
            ODEMMAoGA1UEBhMDVVNBMRYwFAYDVQQKEw1DbG91ZCBGb3VuZHJ5MRAwDgYDVQQD
            DAdib3NoX2NhMB4XDTE3MDcxOTIzMTAyN1oXDTE4MDcxOTIzMTAyN1owOTEMMAoG
            A1UEBhMDVVNBMRYwFAYDVQQKEw1DbG91ZCBGb3VuZHJ5MREwDwYDVQQDEwgxMC4w
            LjAuNjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKyrR/iTZhL+KLE2
            cIchwS2riWkI5jQEub2AwnqehfQFKVt5C43A2BpPgcMha6S2uxApCDW2AfuYgRJG
            wUStLETCM+a81FlO0L/xwOtPaxf8p20lEs5Nbr/X4EwtgVZ29n+y9Nx9cKprYIad
            2ybyfy5zKnt8fjlRKZK9QOcal7MrIQ0lOrYHzPuLTcgQULnVwr9ycsiqAif3EKSw
            jofZ0NMjuK9X6jT237mEn94hkX7hYOFYKLyyjqkYVsTcpBsXAdeP7htY8Tdb1/mI
            NWmb1NO+WdNDOjQ3NIGmGnQwY8TRpI2Cs4nByNPxKE+qiZ0frRgc09IKtuK3Ms9i
            4S7T+n8CAwEAAaNGMEQwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUF
            BwMBMAwGA1UdEwEB/wQCMAAwDwYDVR0RBAgwBocECgAABjANBgkqhkiG9w0BAQsF
            AAOCAQEAtjxMoW1duyeo32vMYLHqLU7VnomXYdlLMRCKV/J9pipGSfIum1SuYBOl
            DTFi9pxjw03C4S+qSg13fIHIO3x2eQ2eDotC2QS+ORDgrFXCuxRZBWY7s3B1iLWs
            AWA+G2D9KyNJfsiKwX8SfgOR2dA6ISDobvbCO56BmiOOZJaTMbF4JTsK57bBmUpk
            0B+Z+fwGpVFBfIFnMIcIAkDk21eygHnhEB6DqPPMOP/i2VX+vv3HSUNygRgx7hUn
            ztIPn8EfzDq0kTKmT55M8gmXvbzxXmRRBn5s88xqD3r1MW5KCpy/1EfGFJ3tPnv0
            iwq5UxHGDhtvMynjcWqQhIjf7fdjrw==
            -----END CERTIFICATE-----
        user_management:
          provider: local
          local:
            users:
              - name: ${director_username}
                password: ${director_password}
              - name: hm
                password: hm-password
      hm:
        director_account:
          user: hm
          password: hm-password
        resurrector_enabled: true

      google: &google_properties
        project: ${google_project}

      agent:
        mbus: nats://nats:nats-password@${google_address_static_director}:4222
        ntp: *ntp
        blobstore:
           options:
             endpoint: http://${google_address_static_director}:25250
             user: agent
             password: agent-password

      ntp: &ntp
        - 169.254.169.254

cloud_provider:
  template:
    name: google_cpi
    release: bosh-google-cpi

  mbus: https://mbus:mbus-password@${director_ip}:6868

  properties:
    google: *google_properties
    agent: {mbus: "https://mbus:mbus-password@0.0.0.0:6868"}
    blobstore: {provider: local, path: /var/vcap/micro_bosh/data/cache}
    ntp: *ntp
EOF

pushd ${deployment_dir}
  function finish {
    echo "Final state of director deployment:"
    echo "=========================================="
    cat director-manifest-state.json
    echo "=========================================="

    cp -r $HOME/.bosh_init ./
  }
  trap finish ERR

  chmod +x ../bosh-init/bosh-init*

  echo "Using bosh-init version..."
  ../bosh-init/bosh-init* version

  echo "Deploying BOSH Director..."
  ../bosh-init/bosh-init* deploy ${manifest_filename}

  trap - ERR
  finish
popd
