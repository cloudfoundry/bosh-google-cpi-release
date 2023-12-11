# BOSH Google CPI Concourse Pipeline

In order to run the BOSH Google CPI Concourse Pipeline you must have an existing [Concourse](http://concourse.ci) environment. See [Deploying Concourse on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/docs/deploy_concourse.md) for instructions.

* Target your Concourse CI environment:

```
fly -t google login -c <YOUR CONCOURSE URL>
```

* Update the [credentials-template.yml](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/master/ci/credentials-template.yml) file. Note that this configuration file requires a SSH private key for the `vcap` user. In order to generate this SSH key run:

 ```
 ssh-keygen -t rsa -f vcap_key -C "vcap@localhost" -m PEM
 ```

 And then:
 * Add the contents of the generated public SSH key-pair (`vcap_key.pub` file) as a Google Compute Engine [project-level key](https://cloud.google.com/compute/docs/instances/adding-removing-ssh-keys#project-wide).

 Alternatively you may store some or all of these parameters in Credhub if your Concourse environment has it set up. Any parameter stored in Credhub should be removed from the credentials file before running the command below.

* Set the BOSH Google CPI pipeline:

```
fly -t google set-pipeline -p bosh-google-cpi -c pipeline.yml -l credentials.yml -v "cpi_source_branch=<branch>"
```

* Unpause the BOSH Google CPI pipeline:

```
fly -t google unpause-pipeline -p bosh-google-cpi
```
