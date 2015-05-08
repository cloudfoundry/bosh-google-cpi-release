# BOSH Google CPI release

This is a [BOSH](http://bosh.io/) release for the external [BOSH Google CPI](https://github.com/frodenas/bosh-google-cpi/).

## Disclaimer

This is NOT presently a production ready [BOSH Google CPI](https://github.com/frodenas/bosh-google-cpi/) BOSH release. This is a work in progress. It is suitable for experimentation and may not become supported in the future.

## Usage

### Prerequisites:

* A [Google Compute Engine](https://cloud.google.com/compute/) account
* A GCE Static IP Address
* A GCE SSH keypair

### Install the bosh-init CLI:

Install the [bosh-init](https://bosh.io/docs/install-bosh-init.html) tool in your workstation.

### Create a deployment directory

Create a deployment directory to store all artifacts:

```
mkdir google-bosh-deployment
cd google-bosh-deployment
```

### Download the BOSH Google CPI BOSH release

Download the BOSH Google CPI BOSH release inside the previously created deployment directory:

```
TBD
```

### Create a deployment manifest

Create a `google-bosh-manifest.yml` deployment manifest file inside the previously created deployment directory with the following properties:

```
---
name: bosh

releases:
  - name: bosh
    url: https://bosh.io/d/github.com/cloudfoundry/bosh?v=163
    sha1: c0bdfcd479b306c98fdf7e5cb93b882d637c0ec7
  - name: bosh-google-cpi
    url: file://./bosh-google-cpi-1.tgz
    sha1: 6545812c1c8245331b8c420f886dafd24b866eed

resource_pools:
  - name: vms
    network: private
    stemcell:
      url: http://storage.googleapis.com/bosh-stemcells/light-bosh-stemcell-2968-google-kvm-ubuntu-trusty-go_agent.tgz
      sha1: ce5a64c3ecef4fd3e6bd633260dfaa7de76540eb
    cloud_properties:
      machine_type: n1-standard-2
      root_disk_size_gb: 10
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
      - name: public
        static_ips:
          - __STATIC_IP__ # <--- Replace with the static IP

    properties:
      nats:
        address: 127.0.0.1
        user: nats
        password: nats

      redis:
        listen_address: 127.0.0.1
        address: 127.0.0.1
        password: redis

      postgres: &db
        adapter: postgres
        host: 127.0.0.1
        user: postgres
        password: postgres
        database: bosh

      dns:
        address: __STATIC_IP__ # <--- Replace with the static IP
        domain_name: microbosh
        db: *db
        recursor: 8.8.8.8

      blobstore:
        address: __STATIC_IP__ # <--- Replace with the static IP
        provider: dav
        director:
          user: director
          password: director
        agent:
          user: agent
          password: agent

      director:
        address: 127.0.0.1
        name: micro-google
        db: *db
        cpi_job: cpi

      hm:
        http:
          user: hm
          password: hm
        director_account:
          user: admin
          password: admin
        resurrector_enabled: true

      ntp: &ntp
        - 169.254.169.254

      google: &google_properties
        project: __GCE_PROJECT__ # <--- Replace with your GCE project
        json_key: __GCE_JSON_KEY__ # <--- Replace with your GCE JSON key
        default_zone: __GCE_DEFAULT_ZONE__ # <--- Replace with the GCE zone to use by default

      agent:
        mbus: nats://nats:nats@__STATIC_IP__:4222 # <--- Replace with the static IP
        ntp: *ntp

      registry:
        host: __STATIC_IP__ # <--- Replace with the static IP
        username: registry
        password: registry
        port: 25777

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
bosh-init deploy google-bosh-manifest.yml
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
git commit -m "updated BOSH Google CPI"
bosh create release --final
git commit -m "creating vXYZ release"
git tag vXYZ
git push origin master --tags
```

## Copyright

See [LICENSE](https://github.com/frodenas/bosh-google-cpi-boshrelease/blob/master/LICENSE) for details.
Copyright (c) 2015 Ferran Rodenas.
