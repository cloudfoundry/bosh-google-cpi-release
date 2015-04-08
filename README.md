# BOSH Google CPI release

This is a [BOSH](http://bosh.io/) release for the external [BOSH Google CPI](https://github.com/frodenas/bosh-google-cpi/).

## Disclaimer

This is NOT presently a production ready [BOSH Google CPI](https://github.com/frodenas/bosh-google-cpi/) BOSH release. This is a work in progress. It is suitable for experimentation and may not become supported in the future.

## Usage

### Install the BOSH micro CLI:

Install the *experimental* [BOSH micro CLI](https://github.com/cloudfoundry/bosh-init/blob/master/docs/build.md) tool in your workstation.

### Create and download the required BOSH micro CLI artifacts

#### Create a deployment directory

Create a deployment directory:

```
mkdir google-micro-deployment
cd google-micro-deployment
```

#### Create a BOSH micro CLI deployment manifest

Create a `google-cpi-manifest.yml` deployment manifest file inside the previously created deployment directory with the following properties:

```
---
name: micro

networks:
  - name: default
    type: dynamic
    cloud_properties:
      network_name: __GCE_NETWORK_NAME__
      tags:
        - bosh
  - name: vip
    type: vip

resource_pools:
  - name: default
    network: default
    cloud_properties:
      machine_type: __GCE_MACHINE_TYPE__

cloud_provider:
  name: micro-google

  template:
    name: cpi
    release: bosh-google-cpi

  ssh_tunnel:
    host: __STATIC_IP__
    port: 22
    user: __USER__
    private_key: __PRIVATE_KEY_PATH__

  registry: &registry
    schema: http
    host: 127.0.0.1
    port: 25777
    username: admin
    password: admin

  mbus: https://mbus-user:mbus-password@__STATIC_IP__:6868

  properties:
    cpi:
      google: &google_properties
        project: __GCE_PROJECT__
        json_key: __GCE_JSON_KEY__
        default_zone: __GCE_DEFAULT_ZONE__
        access_key_id: __GCS_ACCESS_KEY_ID__
        secret_access_key: __GCS_SECRET_ACCESS_KEY__
      agent:
        mbus: https://mbus-user:mbus-password@0.0.0.0:6868
        ntp: ["169.254.169.254"]
        blobstore:
          provider: local
          options:
            blobstore_path: /var/vcap/micro_bosh/data/cache
      registry: *registry

jobs:
  - name: bosh
    instances: 1
    persistent_disk: 32768
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
      - name: registry
        release: bosh
    networks:
      - name: default
      - name: vip
        static_ips:
          - __STATIC_IP__

    properties:
      nats:
        user: "nats"
        password: "nats"
        auth_timeout: 3
        address: __STATIC_IP__
      redis:
        address: "127.0.0.1"
        password: "redis"
        port: 25255
      postgres: &bosh_db
        adapter: "postgres"
        user: "postgres"
        password: "postgres"
        host: "127.0.0.1"
        database: "bosh"
        port: 5432
      blobstore:
        address: __STATIC_IP__
        director:
          user: "director"
          password: "director"
        agent:
          user: "agent"
          password: "agent"
        provider: "dav"
      director:
        address: "127.0.0.1"
        name: "micro"
        port: 25555
        db: *bosh_db
        backend_port: 25556
        cpi_job: cpi
      registry:
        address: __STATIC_IP__
        db: *bosh_db
        http:
          user: "admin"
          password: "admin"
          port: 25777
      hm:
        http:
          user: "hm"
          password: "hm"
        director_account:
          user: "admin"
          password: "admin"
      dns:
        address: __STATIC_IP__
        domain_name: "microbosh"
        db: *bosh_db
      ntp: []
      google: *google_properties
```

#### Download the BOSH Google Stemcell

Download the BOSH Google Stemcell inside the previously created deployment directory:

```
wget http://storage.googleapis.com/bosh-stemcells/light-bosh-stemcell-2479-google-kvm-ubuntu-trusty.tgz
```

#### Download the BOSH Google CPI BOSH release

TBD

#### Download the BOSH release

Download the BOSH release inside the previously created deployment directory:

```
curl -L -J -O https://bosh.io/d/github.com/cloudfoundry/bosh?v=158
```

### Deploy

Using the previous created deployment manifest and the downloaded artifacts, now we can deploy it:

```
bosh-init deploy google-cpi-manifest.yml light-bosh-stemcell-2479-google-kvm-ubuntu-trusty.tgz bosh-google-cpi-1.tgz bosh-158.tgz
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
