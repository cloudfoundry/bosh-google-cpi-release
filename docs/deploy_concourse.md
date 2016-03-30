# Deploying Concourse on Google Compute Engine

In order to deploy [Concourse](http://concourse.ci/) on [Google Compute Engine](https://cloud.google.com/) follow these steps:

### Prerequisites

* An existing BOSH environment

* Ensure that you have enough [Resource Quotas](https://cloud.google.com/compute/docs/resource-quotas) available:
    - 24 Cores
    - 3 IP addresses
    - 160 Gb persistent disk

### Prepare the Google Compute Engine environment

* Reserve a new [static external IP address](https://cloud.google.com/compute/docs/instances-and-network#reserve_new_static):

```
$ gcloud compute addresses create concourse
```

* Create the following firewalls and [set the appropriate rules](https://cloud.google.com/compute/docs/networking#addingafirewall):

```
$ gcloud compute firewall-rules create concourse \
  --description "Concourse Public Traffic" \
  --network cf \
  --target-tags concourse \
  --allow tcp:8080
```

### Deploying Concourse

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
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.333.0
$ bosh upload release https://bosh.io/d/github.com/concourse/concourse?v=0.74.0
```

* Download the [concourse.yml](https://raw.githubusercontent.com/cloudfoundry-incubator/bosh-google-cpi-boshrelease/master/docs/concourse.yml) deployment manifest file and update it with your properties (at the top of the file):
    - `director_uuid = 'CHANGE-ME'`: replace `CHANGE-ME` with your BOSH UUID (run `bosh status`)
    - `vip_ip = 'CHANGE-ME'`: replace `CHANGE-ME` with the static IP reserved previously (named `concourse`)

* Target the deployment file and deploy:

```
$ bosh deployment concourse.yml
$ bosh deploy
```

* Follow the [Concourse Getting Started](http://concourse.ci/getting-started.html) guide.
