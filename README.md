# BOSH Google CPI release

This is a [BOSH](http://bosh.io/) release for the BOSH Google CPI.

## Releases
Please see [CHANGELOG.md] for details of each release.
<!--The Releases section is automatically generated. Do not edit-->
### CPI

|Version|SHA1|Date|
|---|---|---|
|[24.0.0](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-24.tgz)|e2f77a0a8696b29fdb676cf447cfd9bc6841b648|2016-07-22|
|[24.1.0](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-24.1.0.tgz)|dda781faa6673430ce77d708ac1c4be3cb763886|2016-07-25|
|[24.2.0](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-24.2.0.tgz)|80d3ef039cb0ed014e97eeea10569598804659d3|2016-07-26|
|[24.3.0](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-24.3.0.tgz)|f62cebd284d682121a5b4075d0c60a47dd3981ca|2016-07-27|
|[24.4.0](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-24.4.0.tgz)|0c8c8efc316e5d1e0b2e4665b88dfb044b4b87a3|2016-08-10|
[//]: # (new-cpi)

### Stemcell

|Version|SHA1|Date|
|---|---|---|
|[3262.2 (Light)](https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3262.2-google-kvm-ubuntu-trusty-go_agent.tgz)|f46d82a6ae6e89a5635cb3122389f0c8459a82e0|2016-07-22|
|[3262.2 (Heavy)](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-stemcell-3262.2-google-kvm-ubuntu-trusty-go_agent.tgz)|f294226d3ade9e27b67e4ef82b8cd5d3b4e23af7|2016-07-25|
|[3262.4 (Light)](https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3262.4-google-kvm-ubuntu-trusty-go_agent.tgz)|1f44ee6fc5fd495113694aa772d636bf1a8d645a|2016-07-26|
|[3262.5 (Light)](https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3262.5-google-kvm-ubuntu-trusty-go_agent.tgz)|b7ed64f1a929b9a8e906ad5faaed73134dc68c53|2016-08-10|
[//]: # (new-stemcell)

## Usage
If you are not familiar with [BOSH](http://bosh.io/) and its terminology please take a look at the [BOSH documentation](http://bosh.io/docs).

## Deploy a BOSH Director on Google Cloud Platform
Complete instructions for deploying a BOSH Director are available in the [docs/bosh/README.md](docs/bosh/README.md) file.


## Deploy other software
After you have followed the instructions for deploying a BOSH director in [docs/bosh/README.md](docs/bosh/README.md), you may deploy releases like CloudFoundry by following the links below:

* [Deploying Cloud Foundry on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/cloudfoundry)

## Contributing
For detailes on how to contribute to this project - including filing bug reports and contributing code changes - pleasee see [CONTRIBUTING.md].

[CHANGELOG.md]: CHANGELOG.md
[CONTRIBUTING.md]: CONTRIBUTING.md
