# BOSH Google CPI release

This is a [BOSH](http://bosh.io/) release for the BOSH Google CPI.

## Releases
Releases are available on bosh.io: [https://bosh.io/releases/github.com/cloudfoundry-incubator/bosh-google-cpi-release?all=1](https://bosh.io/releases/github.com/cloudfoundry-incubator/bosh-google-cpi-release?all=1). Please see [CHANGELOG.md] for details of each release.

### Stemcell
Stemcells are available on bosh.io: [http://bosh.io/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent](http://bosh.io/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent)

## Usage
If you are not familiar with [BOSH](http://bosh.io/) and its terminology please take a look at the [BOSH documentation](http://bosh.io/docs).

### Deploy a BOSH Director on Google Cloud Platform
Complete instructions for deploying a BOSH Director are available in the [docs/bosh/README.md](docs/bosh/README.md) file.

### Deploy other software
After you have followed the instructions for deploying a BOSH director in [docs/bosh/README.md](docs/bosh/README.md), you may deploy releases like CloudFoundry by following the links below:

* [Deploying Cloud Foundry on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/cloudfoundry)

## Developing
Contributions to the CPI are welcome. Unit and integration tests for any new features are encouraged.
Developers may find it easier to set the GOPATH to the directory of the check-out repository:

```
cd bosh-google-cpi-release
export GOPATH=$PWD
PATH=$PATH:$GOPATH/bin
```

### Running integration tests
1. Set your project and [credentials](https://developers.google.com/identity/protocols/application-default-credentials) (alternatively, you may set `GOOGLE_CREDENTIALS` to the contents of your JSON credentials file):

  ```
  export GOOGLE_PROJECT=your-project-id
  export GOOGLE_APPLICATION_CREDENTIALS=~/google_creds.json
  ```

1. Create the infrastructure required to run tests:

  ```
  pushd src/bosh-google-cpi/; make configint; popd
  ```

1. Run the integration tests:

  ```
  pushd src/bosh-google-cpi/; make testint; popd
  ```

1. To destroy the infrastructure required to run the integration tests, execute:

  ```
  pushd src/bosh-google-cpi/; make cleanint; popd
  ```

#### Terraform Permissions

The Google account which terraforms the environment should have, at a minimum, the following permissions: _Compute Image User_, _Compute Instance Admin_, _Compute Network Admin_, _Service Account Admin_.

## Contributing
For details on how to contribute to this project - including filing bug reports and contributing code changes - pleasee see [CONTRIBUTING.md].

[CHANGELOG.md]: CHANGELOG.md
[CONTRIBUTING.md]: CONTRIBUTING.md
