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

### Running unit tests

There is a Makefile target for running unit tests on the Golang code for the
CPI.

```
cd src/bosh-google-cpi
make test
```

### Running integration tests
1. Set your project:

  ```
  export GOOGLE_PROJECT=your-project-id
  export GOOGLE_JSON_KEY=your-file-location-of-gcp-json-key
  export GOOGLE_APPLICATION_CREDENTIALS=your-file-location-of-gcp-json-key
  ```

2. Create the infrastructure required to run tests:

  ```
  make configint
  ```

3. Run the integration tests:

  ```
  make testint
  ```

To destroy the infrastructure required to run the integration tests, execute:

  ```
  make cleanint
  ```

### Running ERB job templates unit tests

The ERB templates rendered by the jobs of this Bosh Release have unit tests
using Ruby. The required Ruby version is specified in `.ruby-version` as per
convention with `chruby` or similar tools. A script will help you to install
the correct Ruby version if necessary and run the ERB unit tests:

```
./scripts/test-unit-erb
```

## Contributing
For detailes on how to contribute to this project - including filing bug reports and contributing code changes - pleasee see [CONTRIBUTING.md].

[CHANGELOG.md]: CHANGELOG.md
[CONTRIBUTING.md]: CONTRIBUTING.md
