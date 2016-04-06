# Deploying Cloud Foundry MySQL Service on Google Compute Engine

In order to deploy the [Cloud Foundry MySQL Service](https://github.com/cloudfoundry/cf-mysql-release) on [Google Compute Engine](https://cloud.google.com/) follow these steps:

### Prerequisites

* An existing BOSH environment

* An existing [Cloud Foundry environment](https://github.com/cloudfoundry-incubator/bosh-google-cpi-boshrelease/blob/master/docs/deploy_cf.md)

* Ensure that you have enough [Resource Quotas](https://cloud.google.com/compute/docs/resource-quotas) available:
    - 24 Cores
    - 5 IP addresses
    - 1.5 Tb persistent disk

### Deploying Cloud Foundry MySQL Service

* Target and login into your BOSH environment:

```
$ bosh target <YOUR BOSH IP ADDRESS>
```

Your username is `admin` and password is `admin`.

* Upload the required [Google BOSH Stemcell](http://bosh.io/docs/stemcell.html):

```
$ bosh upload stemcell https://storage.googleapis.com/bosh-stemcells/light-bosh-stemcell-3202-google-kvm-ubuntu-trusty-go_agent.tgz
```

* Upload the required [BOSH Releases](http://bosh.io/docs/release.html):

```
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry/cf-mysql-release?v=26
```

* Download the [mysql.yml](https://raw.githubusercontent.com/cloudfoundry-incubator/bosh-google-cpi-boshrelease/master/docs/mysql.yml) deployment manifest file and update it with your properties (at the top of the file):
    - `director_uuid = 'CHANGE-ME'`: replace `CHANGE-ME` with your BOSH UUID (run `bosh status`)
    - `vip_ip = 'CHANGE-ME'`: replace `CHANGE-ME` with the static IP assigned to your Cloud Foundry environment (named `cf`)

* Target the deployment file and deploy:

```
$ bosh deployment mysql.yml
$ bosh deploy
```

* Register the broker within your Cloud Foundry environment:

```
$ bosh run errand broker-registrar
```

* Now your applications will be able to use MySQL:

```
$ cf marketplace
Getting services from marketplace in org frodenas / space dev as frodenas...
OK

service   plans        description
mysql     100mb, 1gb   MySQL databases on demand

TIP:  Use 'cf marketplace -s SERVICE' to view descriptions of individual plans of a given service.
```

Refer to the [Services Overview](http://docs.cloudfoundry.org/devguide/services/) guide to known how you can integrate this service with applications that have been pushed to Cloud Foundry.
