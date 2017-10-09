# Deploy BOSH on Google Cloud Platform

These instructions walk you through deploying a BOSH Director on Google Cloud Platform using manual networking and a network that allows private IP addresses with outbound Internet access provided by a NAT instance.

## Overview
Here are a few important facts about the architecture of the BOSH deployment you will create in this guide:

1. An isolated Google Compute Engine subnetwork will be created to contain the BOSH director and all deployments it creates.
1. The BOSH director will be created by a bastion instance (named `bosh-bastion`).
1. The bastion host will have a firewall rule allowing SSH access from the Internet.
1. Both the bastion host and BOSH director will have outbound Internet connectivity.
1. The BOSH director will allow inbound connectivity only from the bastion. All `bosh` commands must be executed from the bastion.
1. Both bastion and BOSH director will be deployed to an isolated subnetwork in the parent network.
1. The BOSH director will have a statically-assigned `10.0.0.6` IP address.

The following diagram
[[source image doc](https://docs.google.com/presentation/d/1iDjWRQqlAfTyDEvkhsn24ZYRK5GS86QhBqon_OFmNhQ)]
provides an overview of the deployment:

![](../img/arch-overview.png)

## Configure your [Google Cloud Platform](https://cloud.google.com/) environment

### Signup

1. [Sign up](https://cloud.google.com/compute/docs/signup) for Google Cloud Platform
1. Create a [new project](https://console.cloud.google.com/iam-admin/projects)
1. Enable the [GCE API](https://console.developers.google.com/apis/api/compute_component/overview) for your project
1. Enable the [IAM API](https://console.cloud.google.com/apis/api/iam.googleapis.com/overview) for your project
1. Enable the [Cloud Resource Manager API](https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview)

### Setup

1. In your new project, open Cloud Shell (the small `>_` prompt icon in the web console menu bar).

1. Configure a few environment variables:

   ```
   export project_id=$(gcloud config get-value project)
   export region=us-east1
   export zone=us-east1-d
   export base_ip=10.0.0.0
   export service_account_email=terraform@${project_id}.iam.gserviceaccount.com
   ```

1. Configure `gcloud` to use your preferred region and zone:

   ```
   gcloud config set compute/zone ${zone}
   gcloud config set compute/region ${region}
   ```

1. Create a service account and key:

   ```
   gcloud iam service-accounts create terraform --display-name terraform
   gcloud iam service-accounts keys create ~/terraform.key.json \
       --iam-account ${service_account_email}
   ```

1. Grant the new service account editor access to your project:

   ```
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email} \
     --role roles/owner
   ```

1. Make your service account's key available in an environment variable to be used by `terraform`:

   ```
   export GOOGLE_CREDENTIALS=$(cat ~/terraform.key.json)
   ```

<a name="deploy-xpn"></a>
### Optional: Setup Shared VPC (Formerly XPN)

   [Shared VPC](https://cloud.google.com/compute/docs/shared-vpc/) uses a host project to manage the network resources and client project(s) to deploy compute resources. An [organization](https://cloud.google.com/resource-manager/docs/quickstart-organizations) is required to use Shared VPC and you must be signed in as an organization admin.

   The host project must have the [GCE API](https://console.developers.google.com/apis/api/compute_component/overview), [IAM API](https://console.cloud.google.com/apis/api/iam.googleapis.com/overview), and the [Cloud Resource Manager API](https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview) enabled. This project will be used to create the bosh network throughout this guide.

1. Modify and export the host project ID:

   ```
   export xpn_host_project_id=<existing project that will host the VPC>
   ```

1. Setup the projects for Shared VPC:

   ```
   export org_id=$(gcloud projects describe ${project_id} --format 'json' | jq -r '.parent.id')
   export email=$(gcloud config get-value account)
   gcloud organizations add-iam-policy-binding ${org_id} \
     --member user:${email} \
     --role roles/compute.xpnAdmin

   gcloud beta compute xpn enable ${xpn_host_project_id}
   gcloud beta compute xpn associated-projects add ${project_id} --host-project=${xpn_host_project_id}

   gcloud projects add-iam-policy-binding ${xpn_host_project_id} \
     --member serviceAccount:${service_account_email} \
     --role roles/owner
   ```

<a name="deploy-automatic"></a>
## Deploy supporting infrastructure
The following instructions offer the fastest path to getting BOSH up and running on Google Cloud Platform. Using [Terraform](https://www.terraform.io/) you will provision all of the infrastructure required to run BOSH in just a few commands.

### Steps
> **Note:** All of these steps should be performed inside the Cloud Shell in your browser.

1. Clone this repository and go into the BOSH docs directory:

   ```
   git clone https://github.com/cloudfoundry-incubator/bosh-google-cpi-release.git
   cd bosh-google-cpi-release/docs/bosh
   ```

1. Initialize the Terraform cloud provider
```
   docker run -i -t \
     -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
     -v `pwd`:/$(basename `pwd`) \
     -w /$(basename `pwd`) \
     hashicorp/terraform:light init
```

1. View the Terraform execution plan to see the resources that will be created:

   ```
   docker run -i -t \
     -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
     -v `pwd`:/$(basename `pwd`) \
     -w /$(basename `pwd`) \
     hashicorp/terraform:light plan \
       -var service_account_email=${service_account_email} \
       -var project_id=${project_id} \
       -var region=${region} \
       -var zone=${zone} \
       -var baseip=${base_ip} \
       -var network_project_id=${xpn_host_project_id-$project_id}
   ```

1. Create the resources (should take between 60-90 seconds):

   ```
   docker run -i -t \
     -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
     -v `pwd`:/$(basename `pwd`) \
     -w /$(basename `pwd`) \
     hashicorp/terraform:light apply \
       -var service_account_email=${service_account_email} \
       -var project_id=${project_id} \
       -var region=${region} \
       -var zone=${zone} \
       -var baseip=${base_ip} \
       -var network_project_id=${xpn_host_project_id-$project_id}
   ```

Now you have the infrastructure ready to deploy a BOSH director.

<a name="deploy-bosh"></a>
## Deploy BOSH

1. SSH to the bastion VM you created in the previous step. You can use Cloud Shell to SSH to the bastion, or you can connect from any workstation with `gcloud` installed. All SSH commands after this should be run from the bastion VM.

   ```
   gcloud compute ssh bosh-bastion
   ```

1. If you see a warning indicating the VM isn't ready, log out, wait a few moments, and log in again.

1. NOTE: During the `gcloud` commands below, if you see a suggestion to update, you can safely ignore it.

1. Create a service account. This service account will be used by BOSH and all VMs it creates:

   ```
   export service_account=bosh-user
   export base_ip=10.0.0.0
   export service_account_email=${service_account}@${project_id}.iam.gserviceaccount.com
   gcloud iam service-accounts create ${service_account}
   ```

1. Grant the new service account editor access to your project:

   ```
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email} \
     --role roles/compute.instanceAdmin
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email} \
     --role roles/compute.storageAdmin
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email} \
     --role roles/storage.admin
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email} \
     --role  roles/compute.networkAdmin
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email} \
     --role roles/iam.serviceAccountActor
   [ -n "${network_project_id}" ] && gcloud projects add-iam-policy-binding ${network_project_id} \
     --member serviceAccount:${service_account_email} \
     --role roles/compute.networkUser
   ```

1. Create a passphrase-less SSH key and upload the public component:

   ```
   ssh-keygen -t rsa -f ~/.ssh/bosh -C bosh -q -N ""
   ```

   ```
   gcloud compute project-info add-metadata --metadata-from-file \
            sshKeys=<( gcloud compute project-info describe --format=json | jq -r '.commonInstanceMetadata.items[] | select(.key ==  "sshKeys") | .value' & echo "bosh:$(cat ~/.ssh/bosh.pub)" )
   ```

1. Confirm that `bosh2` is installed by querying its version:

  ```
  bosh2 -v # == 2.x.y
  ```

  > This is using the updated bosh-cli V2. When viewing docs written against
  the V1 CLI or bosh-init, use [this reference](https://bosh.io/docs/cli-v2-diff.html)
  to translate between commands.

1. Create and `cd` to a directory:

   ```
   mkdir google-bosh-director
   cd google-bosh-director
   ```

1. Use `vim` or `nano` to create a BOSH Director deployment manifest named `manifest.yml.erb`:

   ```
   ---
   <%
   ['zone', 'service_account_email', 'network', 'subnetwork', 'project_id', 'network_project_id', 'base_ip'].each do |val|
     if ENV[val].nil? || ENV[val].empty?
       raise "Missing environment variable: #{val}"
     end
   end
   %>
   name: bosh

   releases:
     - name: bosh
       url: https://bosh.io/d/github.com/cloudfoundry/bosh?v=262.3
       sha1: 31d2912d4320ce6079c190f2218c6053fd1e920f
     - name: bosh-google-cpi
       url: https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-google-cpi-release?v=25.9.0
       sha1: 3fbda22fde33878b54dec77f4182f8044be72687

   resource_pools:
     - name: vms
       network: private
       stemcell:
         url: https://s3.amazonaws.com/bosh-gce-light-stemcells/light-bosh-stemcell-3421.9-google-kvm-ubuntu-trusty-go_agent.tgz
         sha1: 408f78a2091d108bb5418964026e73c822def32d
       cloud_properties:
         zone: <%= ENV['zone'] %>
         machine_type: n1-standard-1
         root_disk_size_gb: 40
         root_disk_type: pd-standard
         service_account: <%= ENV['service_account_email'] %>

   disk_pools:
     - name: disks
       disk_size: 32_768
       cloud_properties:
         type: pd-standard

   networks:
     - name: vip
       type: vip
     - name: private
       type: manual
       subnets:
       - range: <%= ENV['base_ip'] %>/29
         gateway: <%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 1 }.join('.') %>
         static: [<%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 3 }.join('.')%>-<%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 7 }.join('.') %>]
         cloud_properties:
           network_name: <%= ENV['network'] %>
           subnetwork_name: <%= ENV['subnetwork'] %>
           xpn_host_project_id: <%= ENV['network_project_id'] %>
           ephemeral_external_ip: false
           tags:
             - internal
             - no-ip

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
         - name: google_cpi
           release: bosh-google-cpi

       resource_pool: vms
       persistent_disk_pool: disks

       networks:
         - name: private
           static_ips: [<%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 6 }.join('.') %>]
           default:
             - dns
             - gateway

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
           address: <%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 6 }.join('.') %>
           domain_name: microbosh
           db: *db
           recursor: 169.254.169.254

         blobstore:
           address: <%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 6 }.join('.') %>
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
                 - name: admin
                   password: admin
                 - name: hm
                   password: hm-password
         hm:
           director_account:
             user: hm
             password: hm-password
           resurrector_enabled: true

         google: &google_properties
           project: <%= ENV['project_id'] %>

         agent:
           mbus: nats://nats:nats-password@<%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 6 }.join('.') %>:4222
           ntp: *ntp
           blobstore:
              options:
                endpoint: http://<%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 6 }.join('.') %>:25250
                user: agent
                password: agent-password

         ntp: &ntp
           - 169.254.169.254

   cloud_provider:
     template:
       name: google_cpi
       release: bosh-google-cpi

     mbus: https://mbus:mbus-password@<%= ENV['base_ip'].split('.').tap{|i| i[-1] = i[-1].to_i + 6 }.join('.') %>:6868

     properties:
       google: *google_properties
       agent: {mbus: "https://mbus:mbus-password@0.0.0.0:6868"}
       blobstore: {provider: local, path: /var/vcap/micro_bosh/data/cache}
       ntp: *ntp

   misc:
     ca_cert: |
      -----BEGIN CERTIFICATE-----
      MIIDHjCCAgagAwIBAgIRAOExOrAaTV6piWNDsKVW/zwwDQYJKoZIhvcNAQELBQAw
      ODEMMAoGA1UEBhMDVVNBMRYwFAYDVQQKEw1DbG91ZCBGb3VuZHJ5MRAwDgYDVQQD
      DAdib3NoX2NhMB4XDTE3MDcxOTIzMTAyN1oXDTE4MDcxOTIzMTAyN1owODEMMAoG
      A1UEBhMDVVNBMRYwFAYDVQQKEw1DbG91ZCBGb3VuZHJ5MRAwDgYDVQQDDAdib3No
      X2NhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5aKoganPJR7dtFb9
      QDQHsAEj5n4BqMyboEDdbziVjGMWtjshfVGfpdcTi7Bxn2ZVRHVJPnq/REJzPw7Z
      aqP910XqECDBmfgNEugFXqeK8baS0BV/c4Dy09r+zaqDUubxntn8tEOZr8w0kd1O
      wj+aLif27fX7JgSYddv/pKCsw7V8OgmbMYqwuafCqJMDePUnEo+uxYsk3LZR8iZP
      O60/A5wDQiyF5hTAKj+5LEUsfPksoxqiizI41EfJV10sOrRvbeFP/6kvyhZICaYn
      3osowpE+RXe2mlWQhVpBJf9aHJAz0FnYZhZ9zA/vrrwC4fcTaTUaXBy4gJaydgLP
      6fRPEQIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAN
      BgkqhkiG9w0BAQsFAAOCAQEAVOUsSs7TaXrer8harfqTlHEls/hjgRfb4R2vsY5l
      mgK/wHSceo0wH6W++0BHmVvSPKdHUUB6zj+PnMqRPMb1YYNkYaUbjdlnCqpavRy0
      I44aq7s4R9mY0w5cIfeunZHZlKICzbqYjKic8d2TTbfpAJMrlztR6Dn4tNWOqycL
      iZybtwRZthcZ8XMCBnNgvJjtchNQj++lNEqPJqjdGioiNyTrZvqYjCy9ItgM94v5
      4e///uVp3zjb0SSSWWYzFIO/5zZgb5KDYEqoHZAr6qSuu0lEL1Lc6GgNLAjOiyJx
      xQaUArOymjoQtfQGvEHqsjBy3fGvSUsDZU3oPxYi5AclVw==
      -----END CERTIFICATE-----
   ```

1. Use `erb` to substitute variables in the template:

   ```
   erb manifest.yml.erb > manifest.yml
   ```

1. Deploy the new manifest to create a BOSH Director. Note that this can take
15-20 minutes to complete. You may want to consider starting this command in a
terminal multiplexer such as `tmux` or `screen`.

   ```
   bosh2 create-env manifest.yml
   ```

1. Fetch `ca_cert.pem` from the manifest:

    ```
    bosh2 interpolate manifest.yml --path /misc/ca_cert > ca_cert.pem
    ```

1. Target your BOSH environment and login:

   ```
   bosh2 alias-env micro-google --environment 10.0.0.6 --ca-cert ca_cert.pem
   bosh2 login -e micro-google
   ```

Your username is `admin` and password is `admin`.

## Deploy other software

* [Deploying Cloud Foundry on Google Compute Engine](../cloudfoundry/README.md)

## Delete resources

Follow these instructions when you're ready to delete your BOSH deployment.

From your `bosh-bastion` instance, delete your BOSH director and other resources.

   ```
   # Delete BOSH Director
   cd ~/google-bosh-director
   bosh2 delete-env manifest.yml

   # Delete custom SSH key
   boshkey="bosh:$(cat ~/.ssh/bosh.pub)"
   gcloud compute project-info add-metadata --metadata-from-file \
          sshKeys=<( gcloud compute project-info describe --format=json | jq -r '.commonInstanceMetadata.items[] | select(.key ==  "sshKeys") | .value' | sed -e "s|$boshkey||" | grep -v ^$ )

   # Delete IAM service account
   gcloud iam service-accounts delete bosh-user@${project_id}.iam.gserviceaccount.com
   ```

From your Cloud Shell instance, run the following command to delete the infrastructure you created in this lab:

   ```
   # Set a few vars, in case they were forgotten
   export project_id=$(gcloud config list 2>/dev/null | grep project | sed -e 's/project = //g')
   export network_project_id=$(gcloud beta compute xpn get-host-project ${project_id} | grep name | cut -d : -f 2 - | tr -d '\n' | tr  -d ' ')
   export region=us-east1
   export zone=us-east1-d
   export service_account_email=terraform@${project_id}.iam.gserviceaccount.com
   export GOOGLE_CREDENTIALS=$(cat ~/terraform.key.json)

   # Go to the place with the Terraform manifest
   cd bosh-google-cpi-release/docs/bosh/

   # Destroy the deployment
   docker run -i -t \
     -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
     -v `pwd`:/$(basename `pwd`) \
     -w /$(basename `pwd`) \
     hashicorp/terraform:light destroy \
       -var project_id=${project_id} \
       -var region=${region} \
       -var zone=${zone} \
       -var baseip=${base_ip} \
       -var network_project_id=${network_project_id-project_id}

   # Clean up your IAM credentials and key
   gcloud iam service-accounts delete ${service_account_email}
   rm ~/terraform.key.json
   ```

## Submitting an Issue
We use the [GitHub issue tracker](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/issues) to track bugs and features.
Before submitting a bug report or feature request, check to make sure it hasn't already been submitted. You can indicate
support for an existing issue by voting it up. When submitting a bug report, please include a
[Gist](http://gist.github.com/) that includes a stack trace and any details that may be necessary to reproduce the bug,
including your gem version, Ruby version, and operating system. Ideally, a bug report should include a pull request with
 failing specs.


## Submitting a Pull Request

1. Fork the project.
2. Create a topic branch.
3. Implement your feature or bug fix.
4. Commit and push your changes.
5. Submit a pull request.
