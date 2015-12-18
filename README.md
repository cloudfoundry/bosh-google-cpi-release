# BOSH Google CPI release

This is a [BOSH](http://bosh.io/) release for the external [BOSH Google CPI](https://github.com/frodenas/bosh-google-cpi/).

## Disclaimer

This is NOT presently a production ready [BOSH Google CPI](https://github.com/frodenas/bosh-google-cpi/) BOSH release. This is a work in progress. It is suitable for experimentation and may not become supported in the future.

## Usage

I am assuming you are familiar with [BOSH](http://bosh.io/) and its terminology. If not, please take a look at the [BOSH documentation](http://bosh.io/docs) before running this procedure.

### Setup the [Google Cloud Platform](https://cloud.google.com/) environment

* [Sign up](https://cloud.google.com/compute/docs/signup) and activate Google Compute Engine, if you haven't already.
* Create a [service account](https://developers.google.com/console/help/new/#serviceaccounts) and secure store the downloaded **JSON Key**.
* [Download and Install](https://cloud.google.com/compute/docs/gcloud-compute/#install) the gcloud command line tool.
* Reserve a new [static external IP address](https://cloud.google.com/compute/docs/instances-and-network#reserve_new_static):

```
$ gcloud compute addresses create bosh --region us-central1
```

* Create a new firewall and [set the appropriate rules](https://cloud.google.com/compute/docs/networking#addingafirewall):

```
$ gcloud compute firewall-rules create bosh --description "BOSH" --target-tags bosh --allow tcp:22 tcp:4222 tcp:6868 tcp:25250 tcp:25555 tcp:25777 udp:53
```

* Create your [SSH keys](https://cloud.google.com/compute/docs/instances#sshing) if you haven't already.

### Install the bosh-init CLI

Install the [bosh-init](https://bosh.io/docs/install-bosh-init.html) tool in your workstation.

### Create a deployment directory

Create a deployment directory to store all `bosh-init` artifacts:

```
$ mkdir google-bosh-deployment
$ cd google-bosh-deployment
```

### Create a deployment manifest

Create a `google-bosh-manifest.yml` deployment manifest file inside the previously created deployment directory with the following properties:

```
---
name: bosh

releases:
  - name: bosh
    url: https://bosh.io/d/github.com/cloudfoundry/bosh?v=236
    sha1: 88dd60313dbd7dd832faa44c90493ffa6cd85448
  - name: bosh-google-cpi
    url: http://storage.googleapis.com/bosh-stemcells/bosh-google-cpi-5.tgz
    sha1: c5de3053f233e6ef42c2a4228fa94179d955cc84

resource_pools:
  - name: vms
    network: private
    stemcell:
      url: https://storage.googleapis.com/bosh-stemcells/light-bosh-stemcell-3153-google-kvm-ubuntu-trusty-go_agent.tgz
      sha1: 7f0e4da9507edb122a25fcc1c17bf1d46651842c
    cloud_properties:
      machine_type: n1-standard-2
      root_disk_size_gb: 40
      root_disk_type: pd-standard

disk_pools:
  - name: disks
    disk_size: 32_768
    cloud_properties:
      type: pd-standard

networks:
  - name: private
    type: dynamic
    cloud_properties:
      network_name: default
      tags:
        - bosh

  - name: public
    type: vip

jobs:
  - name: bosh
    instances: 1

    templates:
      - name: nats
        release: bosh
      - name: redis
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
      - name: cpi
        release: bosh-google-cpi
      - name: registry
        release: bosh-google-cpi

    resource_pool: vms
    persistent_disk_pool: disks

    networks:
      - name: private
        default:
          - dns
          - gateway
      - name: public
        static_ips:
          - __STATIC_IP__ # <--- Replace with the static IP

    properties:
      nats:
        address: 127.0.0.1
        user: nats
        password: nats-password

      redis:
        listen_address: 127.0.0.1
        address: 127.0.0.1
        password: redis-password

      postgres: &db
        listen_address: 127.0.0.1
        host: 127.0.0.1
        user: postgres
        password: postgres-password
        database: bosh
        adapter: postgres

      dns:
        address: __STATIC_IP__ # <--- Replace with the static IP
        domain_name: microbosh
        db: *db
        recursor: 8.8.8.8

      blobstore:
        address: __STATIC_IP__ # <--- Replace with the static IP
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
        cpi_job: cpi
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

      ntp: &ntp
        - 169.254.169.254

      google: &google_properties
        project: __GCE_PROJECT__ # <--- Replace with your GCE project
        json_key: __GCE_JSON_KEY__ # <--- Replace with your GCE JSON key content
        default_zone: __GCE_DEFAULT_ZONE__ # <--- Replace with the GCE zone to use by default

      agent:
        mbus: nats://nats:nats@__STATIC_IP__:4222 # <--- Replace with the static IP
        ntp: *ntp
        blobstore:
           options:
             endpoint: http://__STATIC_IP__:25250 # <--- Replace with the static IP
             user: agent
             password: agent-password

      registry:
        host: __STATIC_IP__ # <--- Replace with the static IP
        username: registry
        password: registry-password

cloud_provider:
  template:
    name: cpi
    release: bosh-google-cpi

  ssh_tunnel:
    host: __STATIC_IP__ # <--- Replace with the static IP
    port: 22
    user: __SSH_USER__ # <--- Replace with the user corresponding to your private SSH key
    private_key: __PRIVATE_KEY_PATH__ # <--- Replace with the location of your google_compute_engine SSH private key

  mbus: https://mbus:mbus@__STATIC_IP__:6868 # <--- Replace with the static IP

  properties:
    google: *google_properties
    agent:
      mbus: https://mbus:mbus@0.0.0.0:6868
      ntp: *ntp
      blobstore:
        provider: local
        options:
          blobstore_path: /var/vcap/micro_bosh/data/cache

```

### Deploy

Using the previously created deployment manifest, now we can deploy it:

```
$ bosh-init deploy google-bosh-manifest.yml
```

## Contributing

In the spirit of [free software](http://www.fsf.org/licensing/essays/free-sw.html), **everyone** is encouraged to help improve this project.

Here are some ways *you* can contribute:

* by using alpha, beta, and prerelease versions
* by reporting bugs
* by suggesting new features
* by writing or editing documentation
* by writing specifications
* by writing code (**no patch is too small**: fix typos, add comments, clean up inconsistent whitespace)
* by refactoring code
* by closing [issues](https://github.com/frodenas/bosh-google-cpi-boshrelease/issues)
* by reviewing patches

### Submitting an Issue
We use the [GitHub issue tracker](https://github.com/frodenas/bosh-google-cpi-boshrelease/issues) to track bugs and features.
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

### Create new release

#### Creating a final release

If you need to create a new final release, you will need to get read/write API credentials to the [@cloudfoundry-community](https://github.com/cloudfoundry-community) s3 account.

Please email [Dr Nic Williams](mailto:&#x64;&#x72;&#x6E;&#x69;&#x63;&#x77;&#x69;&#x6C;&#x6C;&#x69;&#x61;&#x6D;&#x73;&#x40;&#x67;&#x6D;&#x61;&#x69;&#x6C;&#x2E;&#x63;&#x6F;&#x6D;) and he will create unique API credentials for you.

Create a `config/private.yml` file with the following contents:

``` yaml
---
blobstore:
  s3:
    access_key_id:     ACCESS
    secret_access_key: PRIVATE
```

You can now create final releases for everyone to enjoy!

```
bosh create release
# test this dev release
git commit -m "updated BOSH Google CPI release"
bosh create release --final
git commit -m "creating vXYZ release"
git tag vXYZ
git push origin master --tags
```

## Copyright

See [LICENSE](https://github.com/frodenas/bosh-google-cpi-boshrelease/blob/master/LICENSE) for details.
Copyright (c) 2015 Ferran Rodenas.
