# Deploying Cloud Foundry on Google Compute Engine

In order to deploy [Cloud Foundry](https://www.cloudfoundry.org/) on [Google Compute Engine](https://cloud.google.com/) follow these steps:

### Prerequisites

* An existing BOSH environment

* Ensure that you have enough [Resource Quotas](https://cloud.google.com/compute/docs/resource-quotas) available:
    - 100 Cores
    - 25 IP addresses
    - 1 Tb persistent disk

### Prepare the Google Compute Engine environment

* Reserve a new [static external IP address](https://cloud.google.com/compute/docs/instances-and-network#reserve_new_static):

```
$ gcloud compute addresses create cf
```

* Create the following load balancing [health checks](https://cloud.google.com/compute/docs/load-balancing/health-checks):

```
$ gcloud compute http-health-checks create cf-public \
  --description "Cloud Foundry Public Health Check" \
  --timeout "5s" \
  --check-interval "30s" \
  --healthy-threshold "10" \
  --unhealthy-threshold "2" \
  --port 80 \
  --request-path "/info" \
  --host "api.<YOUR CF IP ADDRESS>.xip.io"
```

* Create the following load balancing [target pools](https://cloud.google.com/compute/docs/load-balancing/network/target-pools):

```
$ gcloud compute target-pools create cf-public \
  --description "Cloud Foundry Public Target Pool" \
  --health-check cf-public
```

* Create the following load balancing [forwarding rules](https://cloud.google.com/compute/docs/load-balancing/network/forwarding-rules):

```
$ gcloud compute forwarding-rules create cf-http \
  --description "Cloud Foundry HTTP Traffic" \
  --ip-protocol TCP \
  --port-range 80 \
  --target-pool cf-public \
  --address <YOUR CF IP ADDRESS>
```

```
$ gcloud compute forwarding-rules create cf-https \
  --description "Cloud Foundry HTTPS Traffic" \
  --ip-protocol TCP \
  --port-range 443 \
  --target-pool cf-public \
  --address <YOUR CF IP ADDRESS>
```

```
$ gcloud compute forwarding-rules create cf-ssh \
  --description "Cloud Foundry SSH Traffic" \
  --ip-protocol TCP \
  --port-range 2222 \
  --target-pool cf-public \
  --address <YOUR CF IP ADDRESS>
```

```
$ gcloud compute forwarding-rules create cf-wss \
  --description "Cloud Foundry WSS Traffic" \
  --ip-protocol TCP \
  --port-range 4443 \
  --target-pool cf-public \
  --address <YOUR CF IP ADDRESS>
```

* Create the following firewalls and [set the appropriate rules](https://cloud.google.com/compute/docs/networking#addingafirewall):

```
$ gcloud compute firewall-rules create cf-public \
  --description "Cloud Foundry Public Traffic" \
  --network cf \
  --target-tags cf-public \
  --allow tcp:80,tcp:443,tcp:2222,tcp:4443
```

### Deploying Cloud Foundry

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
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry/cf-mysql-release?v=23
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.333.0
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/etcd-release?v=36
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/diego-release?v=0.1454.0
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry/cf-release?v=231
```

* Download the [cloudfoundry.yml](https://raw.githubusercontent.com/cloudfoundry-incubator/bosh-google-cpi-boshrelease/master/docs/cloudfoundry.yml) deployment manifest file and update it with your properties (at the top of the file):
    - `director_uuid = 'CHANGE-ME'`: replace `CHANGE-ME` with your BOSH UUID (run `bosh status`)
    - `vip_ip = 'CHANGE-ME'`: replace `CHANGE-ME` with the static IP reserved previously (named `cf`)

* Target the deployment file and deploy:

```
$ bosh deployment cloudfoundry.yml
$ bosh deploy
```

* Once deployed, you can target your Cloud Foundry environment using the [CF CLI](http://docs.cloudfoundry.org/cf-cli/). Your CF API endpoint is `https://api.<YOUR CF IP ADDRESS>.xip.io`, your username is `admin` and your password is `c1oudc0w`.
