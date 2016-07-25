# BOSH Google CPI release

This is a [BOSH](http://bosh.io/) release for the BOSH Google CPI.

## Releases
<!--The Releases section is automatically generated. Do not edit-->
### CPI

|Version|SHA1|Date|
|---|---|---|
|[22.0.0](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-22.tgz)|d39ac2bc02fe5a2287c30e6c728729d2e68b8e1d|2016-07-19|
|[23.0.0](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-23.tgz)|f4ec50c628af404f60870e933e558bf9b2c36424|2016-07-21|
|[24.0.0](https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-24.tgz)|e2f77a0a8696b29fdb676cf447cfd9bc6841b648|2016-07-22|
[//]: # (new-cpi)

### Stemcell

|Version|SHA1|Date|
|---|---|---|
|[3262.2 (Light)](https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3262.2-google-kvm-ubuntu-trusty-go_agent.tgz)|f46d82a6ae6e89a5635cb3122389f0c8459a82e0|2016-07-22|
[//]: # (new-stemcell)

## Usage
If you are not familiar with [BOSH](http://bosh.io/) and its terminology please take a look at the [BOSH documentation](http://bosh.io/docs).

## Deploy a BOSH Director on Google Cloud Platform
Complete instructions for deploying a BOSH Director are available in the [docs/bosh/README.md](docs/bosh/README.md) file.


## Deploy other software
After you have followed the instructions for deploying a BOSH director in [docs/bosh/README.md](docs/bosh/README.md), you may deploy releases like CloudFoundry by following the links below:

* [Deploying Cloud Foundry on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/cloudfoundry)

## Submitting an Issue
We use the [GitHub issue tracker](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/issues) to track bugs and features.
Before submitting a bug report or feature request, check to make sure it hasn't already been submitted. You can indicate
support for an existing issue by voting it up. When submitting a bug report, please include a
[Gist](http://gist.github.com/) that includes a stack trace and any details that may be necessary to reproduce the bug,
including your gem version, Ruby version, and operating system. Ideally, a bug report should include a pull request with
 failing specs.

## Submitting a Pull Request
1. Fork the project.
1. Create a topic branch.
1. Implement your feature or bug fix.
1. Commit and push your changes.
1. Submit a pull request.
