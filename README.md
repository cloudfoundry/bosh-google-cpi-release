# BOSH Google CPI release

This is a [BOSH](http://bosh.io/) release for the BOSH Google CPI.

## Usage

If you are not familiar with [BOSH](http://bosh.io/) and its terminology please take a look at the [BOSH documentation](http://bosh.io/docs).

## Deploy a BOSH Director on Google Cloud Platform
These instructions walk  you through deploy a BOSH Director on Google Cloud Platform using manual networking.

### Configure your [Google Cloud Platform](https://cloud.google.com/) environment

#### Signup and Setup
1. [Sign up](https://cloud.google.com/compute/docs/signup) and activate Google Compute Engine
1. Create a [service account](https://developers.google.com/identity/protocols/OAuth2ServiceAccount), download the JSON Key, and securely store it
1. [Download and install](https://cloud.google.com/sdk/) the Google Cloud SDK command line tool.

#### Configure gcloud

```
$ gcloud auth login
$ gcloud config set project REPLACE_WITH_YOUR_PROJECT_ID
$ gcloud config set compute/zone us-central1-f
```

#### Provision infrastructure
1. Create a new [network with custom subnetwork ranges](https://cloud.google.com/compute/docs/networking):

  ```
  $ gcloud compute networks create cf --mode custom
  ```

1. Create a new subnetwork for BOSH:

  ```
  $ gcloud compute networks subnets create bosh \
      --network cf \
      --range 10.0.0.0/24 \
      --description "Subnet for BOSH Director and bastion" \
      --region us-central1
  ```

1. Create a new subnetwork for BOSH deployments:

  ```
  $ gcloud compute networks subnets create deployments \
      --network cf \
      --range 10.100.0.0/16 \
      --description "Subnet for BOSH deployments" \
      --region us-central1
  ```

1. Create a [firewall](https://cloud.google.com/compute/docs/networking#addingafirewall) to allow all internal traffic between VMs with the `cf-internal` tag:

  ```
  $ gcloud compute firewall-rules create cf-internal \
    --description "Cloud Foundry Internal traffic" \
    --network cf \
    --source-tags cf-internal \
    --target-tags cf-internal \
    --allow tcp,udp,icmp
  ```

1. Create a [firewall](https://cloud.google.com/compute/docs/networking#addingafirewall) to allow all SSH access to the bastion host that you will deploy the BOSH Director from:

  ```
  $ gcloud compute firewall-rules create bosh-ssh \
    --description "BOSH bastion" \
    --network cf \
    --target-tags bosh-ssh \
    --allow tcp:22
  ```

1. Create a bastion VM that you will use to run `bosh-init` and deply the Director:

  ```
  $ gcloud compute instances create bosh-bastion \
      --image ubuntu-14-04 \
      --subnet bosh \
      --private-network-ip 10.0.0.200 \
      --tags bosh-ssh,cf-internal \
      --scopes cloud-platform \
      --metadata "startup-script=apt-get update -y ; apt-get upgrade -y ; apt-get install -y build-essential zlibc zlib1g-dev ruby ruby-dev openssl libxslt-dev libxml2-dev libssl-dev libreadline6 libreadline6-dev libyaml-dev libsqlite3-dev sqlite3 ; gem install bosh_cli ; curl -o /usr/bin/bosh-init https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.94-linux-amd64; chmod +x /usr/bin/bosh-init"
  ```

#### Deploy BOSH Director

1. SSH to the bastion VM you created in the previous step. All SSH commands after this should be run from the VM:

  ```
  $ gcloud compute ssh bosh-bastion
  ```

1. Create a **password-less** SSH key:

  ```
  $ ssh-keygen -t rsa -f ~/.ssh/bosh -C bosh
  ```

1. Navigate to your [project's web console](https://console.cloud.google.com/compute/metadata/sshKeys) and add the new SSH public key by pasting the contents of ~/.ssh/bosh.pub:

  ![](img/add-ssh.png)

  > **Important:** The username field should auto-populate the value `bosh` after you paste the public key. If it does not, be sure there are no newlines or carriage returns being pasted; the value you paste should be a single line.


1. Confirm that `bosh-init` is installed by querying its version:

  ```
  $ bosh-init -v
  ```

1. Create and `cd` to a directory:

  ```
  $ mkdir google-bosh-director
  $ cd google-bosh-director
  ```

1. Use `vim` or `nano` to create a BOSH Director deployment manifest named `manifest.yml`:

  ```
  ---
  name: bosh

  releases:
    - name: bosh
      url: https://bosh.io/d/github.com/cloudfoundry/bosh?v=256.2
      sha1: ff2f4e16e02f66b31c595196052a809100cfd5a8
    - name: bosh-google-cpi
      url: https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-20.tgz
      sha1: 9b4ada4267f1523f5d278e0813eb055e91ea15ed 

  resource_pools:
    - name: vms
      network: private
      stemcell:
        url: https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3218-google-kvm-ubuntu-trusty-go_agent.tgz
        sha1: 3626815e3fabe8afa7d91362eba1e4c0540795b2 
      cloud_properties:
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
        static: [10.0.0.2-10.0.0.7]
        cloud_properties:
          network_name: cf
          subnetwork_name: bosh
          ephemeral_external_ip: true
          tags:
            - cf-internal

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
          default_zone: us-central1-f

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
      agent:
        mbus: https://mbus:mbus-password@0.0.0.0:6868
        blobstore:
          provider: local
          options:
            blobstore_path: /var/vcap/micro_bosh/data/cache
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

1. Deploy the new manifest to create a BOSH Director:

```
$ bosh-init deploy manifest.yml
  ```

1. Target your BOSH environment:

```
$ bosh target 10.0.0.6
```

Your username is `admin` and password is `admin`.

### Deploy other software

* [Deploying Cloud Foundry on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/deploy_cf.md)
* [Deploying Cloud Foundry MySQL Service on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/deploy_mysql.md)
* [Deploying Cloud Foundry Redis Service on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/deploy_redis.md)
* [Deploying Concourse on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/deploy_concourse.md)

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
