# Deploy BOSH on Google Cloud Platform

These instructions walk  you through deploying a BOSH Director on Google Cloud Platform using manual networking.

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
1. [Sign up](https://cloud.google.com/compute/docs/signup) and activate Google Compute Engine
1. [Download and install](https://cloud.google.com/sdk/) the Google Cloud SDK command line tool.

### Setup

1. Set your project ID:

  ```
  export projectid=REPLACE_WITH_YOUR_PROJECT_ID
  ```

1. Export your preferred compute region and zone:

  ```
  export region=us-east1
  export zone=us-east1-d
  ```

1. Configure `gcloud`:

  ```
  gcloud auth login
  gcloud config set project ${projectid}
  gcloud config set compute/zone ${zone}
  gcloud config set compute/region ${region}
  ```
  
1. Create a service account and key:

  ```
  gcloud iam service-accounts create terraform-bosh
  gcloud iam service-accounts keys create /tmp/terraform-bosh.key.json \
      --iam-account terraform-bosh@${projectid}.iam.gserviceaccount.com
  ```

1. Grant the new service account editor access to your project:

  ```
  gcloud projects add-iam-policy-binding ${projectid} \
      --member serviceAccount:terraform-bosh@${projectid}.iam.gserviceaccount.com \
      --role roles/editor
  ```

1. Make your service account's key available in an environment variable to be used by `terraform`:

  ```
  export GOOGLE_CREDENTIALS=$(cat /tmp/terraform-bosh.key.json)
  ```

<a name="deploy-automatic"></a>
## Deploy supporting infrastructure automatically
The following instructions offer the fastest path to getting BOSH up and running on Google Cloud Platform. Using [Terraform](terraform.io) you will provision all of the infrastructure required to run BOSH in just a few commands.

If you would like a detailed understanding of what is being created, you may instead follow the instructions in [Deploy supporting infrastructure manually](#deploy-manual).

### Requirements
You must have the `terraform` CLI installed on your workstation. See
[Download Terraform](https://www.terraform.io/downloads.html) for more details.

### Steps
1. Download the BOSH terraform file - [main.tf](main.tf) - and save it on your workstation where `terraform` is installed.

1. In a terminal from the same directory where `main.tf` is located, view the Terraform execution plan to see the resources that will be created:

  ```
  terraform plan -var projectid=${projectid} -var region=${region} -var zone=${zone}
  ```

1. Create the resources:

  ```
  terraform apply -var projectid=${projectid} -var region=${region} -var zone=${zone}
  ```

Now you have the infrastructure ready to deploy a BOSH director. Go ahead to the [Deploy BOSH](#deploy-bosh) section to do that. 

<a name="deploy-manual"></a>
## Deploy supporting infrastructure manually

> **Note:** Do not follow the steps if you've already completed the [Deploy supporting infrastructure automatically](#deploy-automatic) instructions. Instead, skip ahead to the [Deploy BOSH](#deploy-bosh) section.

The following instructions offer the most detailed path to getting BOSH up and
running on Google Cloud Platform using the `gcloud` CLI. Although it is
recommended you use the [Deploy supporting infrastructure automatically](#deploy-automatic)
instructions for a production environment, these steps can be helpful in
understanding exactly what is being provisioned to support your BOSH environment.

### Requirements
You must have the `gcloud` CLI installed on your workstation. See
[https://cloud.google.com/sdk/](https://cloud.google.com/sdk/) for more details.

### Steps

1. Create a new [network with custom subnetwork ranges](https://cloud.google.com/compute/docs/networking):

  ```
  gcloud compute networks create cf --mode custom
  ```

1. Create a new subnetwork for BOSH:

  ```
  gcloud compute networks subnets create bosh-us-east1 \
      --network cf \
      --range 10.0.0.0/24 \
      --description "Subnet for BOSH Director and bastion"
  ```

1. Create a [firewall](https://cloud.google.com/compute/docs/networking#addingafirewall) to allow all internal traffic between VMs with the `bosh-internal` tag:

  ```
  gcloud compute firewall-rules create bosh-internal \
    --description "BOSH internal traffic" \
    --network cf \
    --source-tags bosh-internal \
    --target-tags bosh-internal \
    --allow tcp,udp,icmp
  ```

1. Create a [firewall](https://cloud.google.com/compute/docs/networking#addingafirewall) to allow all SSH access to the bastion host that you will deploy the BOSH Director from:

  ```
  gcloud compute firewall-rules create bosh-bastion \
    --description "BOSH bastion" \
    --network cf \
    --target-tags bosh-bastion \
    --allow tcp:22
  ```

1. Create a bastion VM that you will use to run `bosh-init` and deply the Director:

  ```
  gcloud compute instances create bosh-bastion \
      --image ubuntu-14-04 \
      --subnet bosh-us-east1 \
      --private-network-ip 10.0.0.200 \
      --tags bosh-bastion,bosh-internal \
      --scopes cloud-platform \
      --metadata startup-script="apt-get update -y
apt-get upgrade -y
apt-get install -y build-essential zlibc zlib1g-dev ruby ruby-dev openssl libxslt-dev libxml2-dev libssl-dev libreadline6 libreadline6-dev libyaml-dev libsqlite3-dev sqlite3
gem install bosh_cli 
curl -o /tmp/cf.tgz https://s3.amazonaws.com/go-cli/releases/v6.20.0/cf-cli_6.20.0_linux_x86-64.tgz
tar -zxvf /tmp/cf.tgz && mv cf /usr/bin/cf && chmod +x /usr/bin/cf
curl -o /usr/bin/bosh-init https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.96-linux-amd64
chmod +x /usr/bin/bosh-init"
  ```

<a name="deploy-bosh"></a>
## Deploy BOSH
Before working this section, you must have deployed the supporting infrastructure on Google Cloud Platform using either the [automatic](#deploy-automatic) or [manual](deploy-manual) steps provided earlier.

1. SSH to the bastion VM you created in the previous step. All SSH commands after this should be run from the VM:

  ```
  gcloud compute ssh bosh-bastion
  ```

1. Configure `gcloud` to use the correct zone and region:

  ```
  zone=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone)
  zone=${zone##*/}
  region=${zone%-*}
  gcloud config set compute/zone ${zone}
  gcloud config set compute/region ${region}
  ```

1. Create a **password-less** SSH key:

  ```
  ssh-keygen -t rsa -f ~/.ssh/bosh -C bosh
  ```

1. Navigate to your [project's web console](https://console.cloud.google.com/compute/metadata/sshKeys) and add the new SSH public key by pasting the contents of ~/.ssh/bosh.pub:

  ![](../img/add-ssh.png)

  > **Important:** The username field should auto-populate the value `bosh` after you paste the public key. If it does not, be sure there are no newlines or carriage returns being pasted; the value you paste should be a single line.


1. Confirm that `bosh-init` is installed by querying its version:

  ```
  bosh-init -v
  ```

1. Create and `cd` to a directory:

  ```
  mkdir google-bosh-director
  cd google-bosh-director
  ```

1. Use `vim` or `nano` to create a BOSH Director deployment manifest named `manifest.yml`:

  ```
  ---
  name: bosh

  releases:
    - name: bosh
      url: https://bosh.io/d/github.com/cloudfoundry/bosh?v=257.3
      sha1: e4442afcc64123e11f2b33cc2be799a0b59207d0
    - name: bosh-google-cpi
      url: https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-25.1.0.tgz
      sha1: f99dff6860731921282dd1bcd097a74beaeb72a4

  resource_pools:
    - name: vms
      network: private
      stemcell:
        url: https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3262.5-google-kvm-ubuntu-trusty-go_agent.tgz
        sha1: b7ed64f1a929b9a8e906ad5faaed73134dc68c53
      cloud_properties:
        zone: {{ZONE}}
        machine_type: n1-standard-4
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
    - name: vip
      type: vip
    - name: private
      type: manual
      subnets:
      - range: 10.0.0.0/29
        gateway: 10.0.0.1
        static: [10.0.0.3-10.0.0.7]
        cloud_properties:
          network_name: cf
          subnetwork_name: bosh-{{REGION}}
          ephemeral_external_ip: true
          tags:
            - bosh-internal

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
          project: {{PROJECT_ID}}

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
      private_key: {{SSH_KEY_PATH}}

    mbus: https://mbus:mbus-password@10.0.0.6:6868

    properties:
      google: *google_properties
      agent: {mbus: "https://mbus:mbus-password@0.0.0.0:6868"}
      blobstore: {provider: local, path: /var/vcap/micro_bosh/data/cache}
      ntp: *ntp
  ```

1. Run this `sed` command to insert your Google Cloud Platform project ID into the manifest:

  ```
  sed -i s#{{PROJECT_ID}}#`curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/project/project-id`# manifest.yml
  ```

1. Run this `sed` command to insert the full path of the SSH private key you created earlier:

  ```
  sed -i s#{{SSH_KEY_PATH}}#$HOME/.ssh/bosh# manifest.yml
  ```

1. Run this `sed` command to insert the region for your director:

  ```
  sed -i s#{{REGION}}#$region# manifest.yml
  ```

1. Run this `sed` command to insert the zone for your director:

  ```
  sed -i s#{{ZONE}}#$zone# manifest.yml
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
