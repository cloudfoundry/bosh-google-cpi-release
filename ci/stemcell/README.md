# BOSH Google Stemcell Concourse Pipeline

The official BOSH Google Light Stemcells are produced [here](https://bosh-cpi.ci.cf-app.com/pipelines/light-gce-stemcells?groups=light-gce-stemcells).

In order to create your own BOSH Google Stemcell Concourse Pipeline you must have an existing [Concourse](http://concourse.ci) environment. See [Deploying Concourse on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/deploy_concourse.md) for instructions.

* Target your Concourse CI environment:

```
fly -t google login -c <YOUR CONCOURSE URL>
```

* Update the [credentials.yml](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/ci/stemcell/credentials.yml) file.

* Set the BOSH Google Stemcell pipeline:

```
fly -t google set-pipeline -p bosh-google-stemcell -c pipeline.yml -l credentials.yml
```

* Unpause the BOSH Google Stemcell pipeline:

```
fly -t google unpause-pipeline -p bosh-google-stemcell
```
