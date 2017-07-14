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
   gcloud iam service-accounts create terraform
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

1. Create a **password-less** SSH key and upload the public component:

   ```
   ssh-keygen -t rsa -f ~/.ssh/bosh -C bosh
   ```

   ```
   gcloud compute project-info add-metadata --metadata-from-file \
            sshKeys=<( gcloud compute project-info describe --format=json | jq -r '.commonInstanceMetadata.items[] | select(.key ==  "sshKeys") | .value' & echo "bosh:$(cat ~/.ssh/bosh.pub)" )
   ```

1. Confirm that `bosh-init` is installed by querying its version:

   ```
   bosh-init -v
   ```

1. Create and `cd` to a directory:

   ```
   mkdir google-bosh-director
   cd google-bosh-director
   ```

1. Use `vim` or `nano` to create a BOSH Director deployment manifest named `manifest.yml.erb`:

   ```
   ---
   <%
   ['zone', 'service_account_email', 'network', 'subnetwork', 'project_id', 'network_project_id', 'ssh_key_path', 'base_ip'].each do |val|
     if ENV[val].nil? || ENV[val].empty?
       raise "Missing environment variable: #{val}"
     end
   end
   %>
   name: bosh

   releases:
     - name: bosh
       url: https://bosh.io/d/github.com/cloudfoundry/bosh?v=260.3
       sha1: 22c79db2a785efa9cbc32c62b8094500e952e170
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
         - name: postgres
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
   ```

1. Use `erb` to substitute variables in the template:

   ```
   erb manifest.yml.erb > manifest.yml
   ```

1. Deploy the new manifest to create a BOSH Director. Note that this can take
15-20 minutes to complete. You may want to consider starting this command in a
terminal multiplexer such as `tmux` or `screen`.

   ```
   bosh-init deploy manifest.yml
   ```

1. Target your BOSH environment:

   ```
   bosh target 10.0.0.6
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
   bosh-init delete manifest.yml

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
