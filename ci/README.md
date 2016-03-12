# BOSH Google CPI Concourse Pipeline

In order to run the BOSH Google CPI Concourse Pipeline you must have an existing [Concourse](http://concourse.ci) environment. See [Deploying Concourse on Google Compute Engine](https://github.com/frodenas/bosh-google-cpi-boshrelease/blob/master/docs/deploy_concourse.md) for instructions.

* Target your Concourse CI environment:

```
fly -t google login -c <YOUR CONCOURSE URL>
```

* Update the [credentials.yml](https://github.com/frodenas/bosh-google-cpi-boshrelease/blob/master/ci/credentials.yml) file.

* Set the BOSH Google CPI pipeline:

```
fly -t google set-pipeline -p bosh-google-cpi -c pipeline.yml -l credentials.yml
```

* Unpause the BOSH Google CPI pipeline:

```
fly -t google unpause-pipeline -p bosh-google-cpi
```
