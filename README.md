# BOSH Google CPI release

This is a [BOSH](http://bosh.io/) release for the BOSH Google CPI.

## Releases
Releases are available on bosh.io: [https://bosh.io/releases/github.com/cloudfoundry-incubator/bosh-google-cpi-release?all=1](https://bosh.io/releases/github.com/cloudfoundry-incubator/bosh-google-cpi-release?all=1). Please see [CHANGELOG.md] for details of each release.

### Stemcell
Stemcells are available on bosh.io: [http://bosh.io/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent](http://bosh.io/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent)

## Usage
If you are not familiar with [BOSH](http://bosh.io/) and its terminology please take a look at the [BOSH documentation](http://bosh.io/docs).

### Deploy a BOSH Director on Google Cloud Platform
[BOSH Bootloader](https://github.com/cloudfoundry/bosh-bootloader) is the recommended way to deploy a BOSH director on GCP. Detailed instructions are available [here](https://github.com/cloudfoundry/bosh-bootloader/blob/master/docs/getting-started-gcp.md).

### Deploy other software
After you have followed the instructions for deploying a BOSH director in [docs/bosh/README.md](docs/bosh/README.md), you may deploy releases like Cloud Foundry or Concourse by following the links below:

* [Deploying Cloud Foundry on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/cloudfoundry)
* [Deploying Concourse on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/concourse)

## Developing
Contributions to the CPI are welcome. Unit and integration tests for any new features are encouraged.
Developers may find it easier to set the GOPATH to the directory of the check-out repository:

```
cd bosh-google-cpi-release
export GOPATH=$pwd
PATH=$PATH:$GOPATH/bin
```

### Running integration tests
1. Set your project:

  ```
  export GOOGLE_PROJECT=your-project-id
  ```

1. Create the infrastructure required to run tests:

  ```
  make configint
  ```

1. Run the integration tests:

  ```
  make testint
  ```

To destroy the infrastructure required to run the integration tests, execute:

  ```
  make cleanint
  ```

## Contributing
For detailes on how to contribute to this project - including filing bug reports and contributing code changes - pleasee see [CONTRIBUTING.md].

[CHANGELOG.md]: CHANGELOG.md
[CONTRIBUTING.md]: CONTRIBUTING.md
