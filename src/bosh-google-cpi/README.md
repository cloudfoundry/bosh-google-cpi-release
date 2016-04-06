# BOSH Google Compute Engine CPI [![Build Status](https://travis-ci.org/frodenas/bosh-google-cpi.png)](https://travis-ci.org/frodenas/bosh-google-cpi)

This is an **experimental** external [BOSH CPI](http://bosh.io/docs/bosh-components.html#cpi) for [Google Compute Engine](https://cloud.google.com/compute/).

## Disclaimer

This is **NOT** presently a production ready CPI. This is a work in progress. It is suitable for experimentation and may not become supported in the future.

## Usage

### Deployment
This is the implemention of the CPI, and is part of the [BOSH Google CPI release](https://github.com/cloudfoundry-incubator/bosh-google-cpi-boshrelease) repository. See the README at the root of this repository for instructions on deploying the release with this CPI.

### Installation

The source for this CPI is not intended to be deployed except as a BOSH deployment and is not `go get`able. To build or install the CPI locally for development or test purposes, you can symlink the repository into your Go workspace.

From the root of the `bosh-google-cpi-boshrelease` repository:

```
ln -s $(pwd)/src/bosh-google-cpi $GOPATH/src/
```

You can now `go build` or `go install` the "main" package.

### Configuration

Create a configuration file:

```
{
  "google": {
    "project": "my-gce-project",
    "default_zone": "us-central1-a",
    "json_key": "{\"private_key_id\": \"...\"}",
    "default_root_disk_size_gb": 20,
    "default_root_disk_type": ""
  },
  "actions": {
    "agent": {
      "mbus": "https://mbus:mbus@0.0.0.0:6868",
      "ntp": [
        "169.254.169.254"
      ],
      "blobstore": {
        "type": "local",
        "options": {}
      }
    },
    "registry": {
      "protocol": "http",
      "host": "127.0.0.1",
      "port": 25777,
      "username": "admin",
      "password": "admin",
      "tls": {
        "_comment": "TLS options only apply when using HTTPS protocol",
        "insecure_skip_verify": true,
        "certfile": "/path/to/public.pem",
        "keyfile": "/path/to/private.pem",
        "cacertfile": "/path/to/ca.pem"
      }
    }
  }
}
```

| Option                                    | Required   | Type          | Description
|:------------------------------------------|:----------:|:------------- |:-----------
| google.project                            | Y          | String        | Google Compute Engine [Project](https://cloud.google.com/compute/docs/projects)
| google.default_zone                       | Y          | String        | Google Compute Engine default [Zone](https://cloud.google.com/compute/docs/zones)
| google.json_key                           | N         | String        | Contents of the Google Compute Engine [JSON file](https://developers.google.com/identity/protocols/application-default-credentials). Only required if you are not running the CPI inside a Google Compute Engine VM with `compute` and `devstorage.full_control` service scopes and/or the Google Cloud SDK has not been initialized
| google.default_root_disk_size_gb          | N          | Integer       | The default size (in Gb) of the instance root disk (default is `10Gb`)
| google.default_root_disk_type             | N          | String        | The name of the default [Google Compute Engine Disk Type](https://cloud.google.com/compute/docs/disks/#overview) the CPI will use when creating the instance root disk
| actions.agent.mbus.endpoint               | Y          | String        | [BOSH Message Bus](http://bosh.io/docs/bosh-components.html#nats) URL used by deployed BOSH agents
| actions.agent.ntp                         | Y          | Array&lt;String&gt; | List of NTP servers used by deployed BOSH agents
| actions.agent.blobstore.type              | Y          | String        | Provider type for the [BOSH Blobstore](http://bosh.io/docs/bosh-components.html#blobstore) used by deployed BOSH agents (e.g. dav, s3)
| actions.agent.blobstore.options           | Y          | Hash          | Options for the [BOSH Blobstore](http://bosh.io/docs/bosh-components.html#blobstore) used by deployed BOSH agents
| actions.registry.protocol                 | Y          | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Protocol (`http` or `https`)
| actions.registry.host                     | Y          | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Host
| actions.registry.port                     | Y          | Integer       | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) port
| actions.registry.username                 | Y          | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) username
| actions.registry.password                 | Y          | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) password
| actions.registry.tls.insecure_skip_verify | When https | Boolean       | Skip [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) server's certificate chain and host name verification
| actions.registry.tls.certfile             | When https | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Client Certificate (PEM format) file location
| actions.registry.tls.keyfile              | When https | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Client Key (PEM format) file location
| actions.registry.tls.cacertfile           | When https | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Client Root CA certificates (PEM format) file location

### Run

Run the cpi using the previously created configuration file:

```
$ echo "{\"method\": \"method_name\", \"arguments\": []}" | cpi -configFile="/path/to/configuration_file.json"
```

## Features

### BOSH Network options

The BOSH Google Compute Engine CPI supports these [BOSH Networks Types](http://bosh.io/docs/networks.html):

| Type    | Description
|:-------:|:-----------
| dynamic | To use DHCP assigned IPs by Google Compute Engine
| vip     | To use previously allocated Google Compute Engine Static IPs


These options are specified under `cloud_properties` at the [networks](http://bosh.io/docs/networks.html) section of a BOSH deployment manifest and are only valid for `dynamic` networks:

| Option                | Required | Type          | Description
|:----------------------|:--------:|:------------- |:-----------
| network_name          | N        | String        | The name of the [Google Compute Engine Network](https://cloud.google.com/compute/docs/networking#networks) the CPI will use when creating the instance (if not set, by default it will use the `default` network)
| subnetwork_name       | N        | String        | The name of the [Google Compute Engine Subnet Network](https://cloud.google.com/compute/docs/networking#subnet_network) the CPI will use when creating the instance (if the network is in legacy mode, do not provide this property. If the network is in auto subnet mode, providing the subnetwork is optional. If the network is in custom subnet mode, then this field should be specified)
| ephemeral_external_ip | N        | Boolean       | If instances must have an [ephemeral external IP](https://cloud.google.com/compute/docs/instances-and-network#externaladdresses) (`false` by default)
| ip_forwarding         | N        | Boolean       | If instances must have [IP forwarding](https://cloud.google.com/compute/docs/networking#canipforward) enabled (`false` by default)
| target_pool           | N        | String        | The name of the [Google Compute Engine Target Pool](https://cloud.google.com/compute/docs/load-balancing/network/target-pools) the instances should be added to
| instance_group        | N        | String        | The name of the [Google Compute Engine Instance Group](https://cloud.google.com/compute/docs/instance-groups/unmanaged-groups) the instances should be added to
| tags                  | N        | Array&lt;String&gt; | A list of [tags](https://cloud.google.com/compute/docs/instances/managing-instances#tags) to apply to the instances, useful if you want to apply firewall or routes rules based on tags

### BOSH Resource pool options

These options are specified under `cloud_properties` at the [resource_pools](http://bosh.io/docs/deployment-basics.html#resource-pools) section of a BOSH deployment manifest:

| Option              | Required | Type          | Description
|:--------------------|:--------:|:------------- |:-----------
| machine_type        | Y        | String        | The name of the [Google Compute Engine Machine Type](https://cloud.google.com/compute/docs/machine-types) the CPI will use when creating the instance (required if not using `cpu` and `ram`)
| cpu                 | Y        | Integer       | Number of vCPUs ([Google Compute Engine Custom Machine Types](https://cloud.google.com/custom-machine-types/)) the CPI will use when creating the instance (required if not using `machine_type`)
| ram                 | Y        | Integer       | Amount of memory ([Google Compute Engine Custom Machine Types](https://cloud.google.com/custom-machine-types/)) the CPI will use when creating the instance (required if not using `machine_type`)
| zone                | N        | String        | The name of the [Google Compute Engine Zone](https://cloud.google.com/compute/docs/zones) where the instance must be created
| root_disk_size_gb   | N        | Integer       | The size (in Gb) of the instance root disk (default is `10Gb`)
| root_disk_type      | N        | String        | The name of the [Google Compute Engine Disk Type](https://cloud.google.com/compute/docs/disks/#overview) the CPI will use when creating the instance root disk
| automatic_restart   | N        | Boolean       | If the instances should be [restarted automatically](https://cloud.google.com/compute/docs/instances/setting-instance-scheduling-options#autorestart) if they are terminated for non-user-initiated reasons (`false` by default)
| on_host_maintenance | N        | String        | [Instance behavior](https://cloud.google.com/compute/docs/instances/setting-instance-scheduling-options#onhostmaintenance) on infrastructure maintenance that may temporarily impact instance performance (supported values are `MIGRATE` (default) or `TERMINATE`)
| preemptible         | N        | Boolean       | If the instances should be [preemptible](https://cloud.google.com/preemptible-vms/) (`false` by default)
| service_scopes      | N        | Array&lt;String&gt; | [Authorization scope names](https://cloud.google.com/docs/authentication#oauth_scopes) for your default service account that determine the level of access your instance has to other Google services (no scope is assigned to the instance by default)

### BOSH Persistent Disks options

These options are specified under `cloud_properties` at the [disk_pools](http://bosh.io/docs/persistent-disks.html#persistent-disk-pool) section of a BOSH deployment manifest:

| Option | Required | Type   | Description
|:-------|:--------:|:------ |:-----------
| type   | N        | String | The name of the [Google Compute Engine Disk Type](https://cloud.google.com/compute/docs/disks/#overview)

## Deployment Manifest Example

This is an example of how Google Compute Engine CPI specific properties are used in a BOSH deployment manifest:

```
---
name: example
director_uuid: 38ce80c3-e9e9-4aac-ba61-97c676631b91

...

networks:
  - name: private
    type: dynamic
    dns:
      - 8.8.8.8
      - 8.8.4.4
    cloud_properties:
      network_name: default
      subnetwork_name: my-subnetwork
      ephemeral_external_ip: false
      ip_forwarding: false
      target_pool: my-load-balancer
      tags:
        - bosh

  - name: public
    type: vip
    cloud_properties: {}
...

resource_pools:
  - name: vms
    network: private
    stemcell:
      name: bosh-google-kvm-ubuntu-trusty-go_agent
      version: latest
    cloud_properties:
      instance_type: n1-standard-2
      zone: us-central1-a
      root_disk_size_gb: 20
      root_disk_type: pd-ssd
      automatic_restart: false
      on_host_maintenance: MIGRATE
      service_scopes:
        - compute.readonly
        - devstorage.read_write
...

disk_pools:
  - name: disks
    disk_size: 32_768
    cloud_properties:
      type: pd-ssd
...

```
