# BOSH Google Compute Engine CPI

This is an external [BOSH CPI](http://bosh.io/docs/bosh-components.html#cpi) for [Google Compute Engine](https://cloud.google.com/compute/) that is jointly developed by Pivotal and Google.

## Usage

### Deployment
This is the implemention of the CPI, and is part of the [BOSH Google CPI release](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release) repository. See the README at the root of this repository for instructions on deploying the release with this CPI.

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
| manual  | To use manually- or BOSH-assigned private IPs
| dynamic | To use DHCP-assigned private IPs from Google Compute Engine
| vip     | To use previously allocated Google Compute Engine Static IPs


These options are specified under `cloud_properties` at the [networks](http://bosh.io/docs/networks.html) section of a BOSH deployment manifest and are only valid for `manual` or `dynamic` networks:

| Option                | Required | Type          | Description
|:----------------------|:--------:|:------------- |:-----------
| network_name          | N        | String        | The name of the [Google Compute Engine Network](https://cloud.google.com/compute/docs/networking#networks) the CPI will use when creating the instance (if not set, by default it will use the `default` network)
| subnetwork_name       | N        | String        | The name of the [Google Compute Engine Subnet Network](https://cloud.google.com/compute/docs/networking#subnet_network) the CPI will use when creating the instance (if the network is in legacy mode, do not provide this property. If the network is in auto subnet mode, providing the subnetwork is optional. If the network is in custom subnet mode, then this field is required)
| ephemeral_external_ip | N        | Boolean       | If instances must have an [ephemeral external IP](https://cloud.google.com/compute/docs/instances-and-network#externaladdresses) (`false` by default). Can be overridden in resource_pools.
| ip_forwarding         | N        | Boolean       | If instances must have [IP forwarding](https://cloud.google.com/compute/docs/networking#canipforward) enabled (`false` by default). Can be overridden in resource_pools.
| tags                  | N        | Array&lt;String&gt; | A list of [tags](https://cloud.google.com/compute/docs/instances/managing-instances#tags) to apply to the instances, useful if you want to apply firewall or routes rules based on tags. Will be merged with tags in resource_pools.

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
| service_account     | N        | String        | The full service account address (e.g., service-account-name@project-name.iam.gserviceaccount.com) of the service account to launch the VM with. If a value is provided, `service_scopes` will default to `https://www.googleapis.com/auth/cloud-platform` unless it is explicitly set. See [service account permissions](https://cloud.google.com/compute/docs/access/service-accounts#service_account_permissions) for more details. To use the default service account, leave this field empty and specify `service_scopes`.
| service_scopes      | N        | Array&lt;String&gt; | Optional. If this value is specified and `service_account` is empty, `default` will be used for `service_account`. This value supports both short (e.g., `cloud-platform`) and fully-qualified (e.g., `https://www.googleapis.com/auth/cloud-platform` formats. See [Authorization scope names](https://cloud.google.com/docs/authentication#oauth_scopes) for more details.
| target_pool         | N        | String        | The name of the [Google Compute Engine Target Pool](https://cloud.google.com/compute/docs/load-balancing/network/target-pools) the instances should be added to
| backend_service     | N        | String        | The name of the [Google Compute Engine Backend Service](https://cloud.google.com/compute/docs/instance-groups/unmanaged-groups) the instances should be added to. The backend service must already be configured with an [Instance Group](https://cloud.google.com/compute/docs/instance-groups/unmanaged-groups)in the same zone as this instance
| ephemeral_external_ip | N        | Boolean       | Overrides the equivalent option in the networks section
| ip_forwarding         | N        | Boolean       | Overrides the equivalent option in the networks section
| tags                  | N        | Array&lt;String&gt; | Merged with tags from the networks section 

### BOSH Persistent Disks options

These options are specified under `cloud_properties` at the [disk_pools](http://bosh.io/docs/persistent-disks.html#persistent-disk-pool) section of a BOSH deployment manifest:

| Option | Required | Type   | Description
|:-------|:--------:|:------ |:-----------
| type   | N        | String | The name of the [Google Compute Engine Disk Type](https://cloud.google.com/compute/docs/disks/#overview)

## Deployment Manifest Example - Dynamic Networking

This is an example of how Google Compute Engine CPI specific properties are used in a BOSH deployment manifest with dynamic networking:

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

## Deployment Manifest Example - Manual Networking

This is an example of how Google Compute Engine CPI specific properties are used in a BOSH deployment manifest with manual networking. This assumes you've created a networked named `custom-network` and a subnetwork named `custom-subnetwork` with a CIDR of 10.0.0.0/24:

```
---
name: example
director_uuid: 38ce80c3-e9e9-4aac-ba61-97c676631b91

...

networks:
  - name: private
      type: manual
      subnets:
      - range: 10.0.0.0/24
        gateway: 10.0.0.1
        static: [10.0.0.2-10.0.0.100]
        cloud_properties:
          network_name: custom-network
          subnetwork_name: custom-subnetwork
          ephemeral_external_ip: true
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
