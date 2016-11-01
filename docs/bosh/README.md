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

The following diagram provides an overview of the deployment:

![](../img/arch-overview.png)

## Configure your [Google Cloud Platform](https://cloud.google.com/) environment

### Signup
1. [Sign up](https://cloud.google.com/compute/docs/signup) for Google Cloud Platform
1. Create a [new project](https://console.cloud.google.com/iam-admin/projects)
1. Enable the [IAM API](https://console.cloud.google.com/apis/api/iam.googleapis.com/overview) for your project
1. Enable the [Project API](https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview)
### Setup

1. In your new project, open Cloud Shell

1. Configure a few environment variables:

  ```
  export project_id=$(gcloud config list 2>/dev/null | grep project | sed -e 's/project = //g')
  export region=us-east1
  export zone=us-east1-d
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

<a name="deploy-automatic"></a>
## Deploy supporting infrastructure
The following instructions offer the fastest path to getting BOSH up and running on Google Cloud Platform. Using [Terraform](terraform.io) you will provision all of the infrastructure required to run BOSH in just a few commands.

### Steps
> **Note:** All of these steps should be performed inside the Cloud Shell in your browser.

1. Clone this repository and go into the BOSH docs directory:

  ```
  git clone https://github.com/cloudfoundry-incubator/bosh-google-cpi-release.git
  cd bosh-google-cpi-release/docs/bosh
  ```

1. In a terminal from the same directory where `main.tf` is located, view the Terraform execution plan to see the resources that will be created:

  ```
  docker run -i -t \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v `pwd`:/$(basename `pwd`) \
    -w /$(basename `pwd`) \
    hashicorp/terraform:light plan \
      -var service_account_email=${service_account_email} \
      -var projectid=${project_id} \
      -var region=${region} \
      -var zone=${zone}
  ```

1. Create the resources:

  ```
  docker run -i -t \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v `pwd`:/$(basename `pwd`) \
    -w /$(basename `pwd`) \
    hashicorp/terraform:light apply \
      -var service_account_email=${service_account_email} \
      -var projectid=${project_id} \
      -var region=${region} \
      -var zone=${zone}
  ```

Now you have the infrastructure ready to deploy a BOSH director.

<a name="deploy-bosh"></a>
## Deploy BOSH

1. SSH to the bastion VM you created in the previous step. You can use Cloud Shell to SSH to the bastion, or you can connect from any workstation with `gcloud` installed. All SSH commands after this should be run from the bastion VM.

  ```
  gcloud compute ssh bosh-bastion
  ```

1. If you see a warning indicating the VM isn't ready, log out, wait a few moments, and log in again.

1. Create a service account. This service account will be used by BOSH and all VMs it creates:

  ```
  export service_account=bosh-user
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
  gcloud iam service-accounts add-iam-policy-binding ${service_account_email} \
    --member serviceAccount:${service_account_email} \
    --role roles/iam.serviceAccountActor
  ```

1. Create a **password-less** SSH key and upload the public component:

  ```
  ssh-keygen -t rsa -f ~/.ssh/bosh -C bosh
  gcloud compute project-info add-metadata --metadata-from-file \
           sshKeys=<( gcloud compute project-info describe --format=json | jq -r '.commonInstanceMetadata.items[] | select(.key == "sshKeys") | .value' & echo "bosh:$(cat ~/.ssh/bosh.pub)" )
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
  name: bosh

  releases:
    - name: bosh
      url: https://bosh.io/d/github.com/cloudfoundry/bosh?v=257.3
      sha1: e4442afcc64123e11f2b33cc2be799a0b59207d0
    - name: bosh-google-cpi
      url: https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-google-cpi-release?v=25.6.0
      sha1: 3e01539a1228d62b8015feb388df2234978efaf6

  resource_pools:
    - name: vms
      network: private
      stemcell:
        url: https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent?v=3263.7
        sha1: a09ce8b4acfa9876f52ee7b4869b4b23f27d5ace
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
      - range: 10.0.0.0/29
        gateway: 10.0.0.1
        static: [10.0.0.3-10.0.0.7]
        cloud_properties:
          network_name: <%= ENV['network'] %>
          subnetwork_name: <%= ENV['subnetwork'] %>
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
          static_ips: [10.0.0.6]
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
          address: 10.0.0.6
          domain_name: microbosh
          db: *db
          recursor: 169.254.169.254

        blobstore:
          address: 10.0.0.6
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
          mbus: nats://nats:nats-password@10.0.0.6:4222
          ntp: *ntp
          blobstore:
             options:
               endpoint: http://10.0.0.6:25250
               user: agent
               password: agent-password

        ntp: &ntp
          - 169.254.169.254

  cloud_provider:
    template:
      name: google_cpi
      release: bosh-google-cpi

    ssh_tunnel:
      host: 10.0.0.6
      port: 22
      user: bosh
      private_key: <%= ENV['ssh_key_path'] %>

    mbus: https://mbus:mbus-password@10.0.0.6:6868

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

1. Deploy the new manifest to create a BOSH Director:

  ```
  bosh-init deploy manifest.yml
  ```

1. Target your BOSH environment:

  ```
  bosh target 10.0.0.6
  ```

Your username is `admin` and password is `admin`.

### Deploy other software

* [Deploying Cloud Foundry on Google Compute Engine](../cloudfoundry/README.md)

### Delete resources

From your `bosh-bastion` instance, delete your BOSH director:

  ```
  bosh-init delete manifest.yml
  ```

From your Cloud Shell instance, run the following command to delete the infrastructure you created in this lab:

  ```
  docker run -i -t \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v `pwd`:/$(basename `pwd`) \
    -w /$(basename `pwd`) \
    hashicorp/terraform:light destroy \
      -var projectid=${project_id} \
      -var region=${region} \
      -var zone=${zone}
  ```
 
### Submitting an Issue
We use the [GitHub issue tracker](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/issues) to track bugs and features.
Before submitting a bug report or feature request, check to make sure it hasn't already been submitted. You can indicate
support for an existing issue by voting it up. When submitting a bug report, please include a
[Gist](http://gist.github.com/) that includes a stack trace and any details that may be necessary to reproduce the bug,
including your gem version, Ruby version, and operating system. Ideally, a bug report should include a pull request with
 failing specs.

 
### Submitting a Pull Request

1. Fork the project.
2. Create a topic branch.
3. Implement your feature or bug fix.
4. Commit and push your changes.
5. Submit a pull request.
