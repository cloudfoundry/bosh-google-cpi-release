# BOSH Google CPI Concourse Pipeline

In order to run the BOSH Google CPI Concourse Pipeline you must have an existing [Concourse](http://concourse.ci) environment. See [Deploying Concourse on Google Compute Engine](https://github.com/cloudfoundry-incubator/bosh-google-cpi-boshrelease/blob/master/docs/deploy_concourse.md) for instructions.

* Target your Concourse CI environment:

```
fly -t google login -c <YOUR CONCOURSE URL>
```

* Update the [credentials.yml](https://github.com/cloudfoundry-incubator/bosh-google-cpi-boshrelease/blob/master/ci/credentials.yml) file. Note that this configuration file requires a SSH private key for the `vcap` user. In order to generate this SSH key run:
 
 ```
 ssh-keygen -t rsa -f vcap_key -C "vcap@localhost"
 ```
 
 And then:
 * Add the contents of the generated public SSH key-pair (`vcap_key.pub` file) as a Google Compute Engine [project-level key](https://cloud.google.com/compute/docs/instances/adding-removing-ssh-keys#project-wide).
 * Paste the contents of the generated private SSH key-pair (`vcap_key.pub`) into the `private_key_data` property at the `credentials.yml` file.

* Set the BOSH Google CPI pipeline:

```
fly -t google set-pipeline -p bosh-google-cpi -c pipeline.yml -l credentials.yml
```

* Unpause the BOSH Google CPI pipeline:

```
fly -t google unpause-pipeline -p bosh-google-cpi
```
