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
$ gcloud compute addresses create concourse --global
```

* Create the following firewalls and [set the appropriate rules](https://cloud.google.com/compute/docs/networking#addingafirewall):

```
$ gcloud compute firewall-rules create concourse-public \
  --description "Concourse public traffic" \
  --network cf \
  --target-tags concourse-public \
  --source-ranges 0.0.0.0/0 \
  --allow tcp:8080
```

```
$ gcloud compute firewall-rules create concourse-internal \
  --description "Concourse internal traffic" \
  --network cf \
  --target-tags concourse-internal \
  --source-tags concourse-internal \
  --allow tcp:0-65535,udp:0-65535,icmp
```

* Create a load balancer for Concourse

1. Create an unmanaged instance group and named port:
  ```
  gcloud compute instance-groups unmanaged create concourse-us-central1-f --zone us-central1-f
  gcloud compute instance-groups unmanaged set-named-ports concourse-us-central1-f --named-ports "http:8080"
  ```

2. Create a health check:
  ```
  gcloud compute http-health-checks create concourse --port 8080 --request-path="/login"
  ```
3. Create a backend service:

  ```
  gcloud compute backend-services create concourse --http-health-check "concourse" --port 8080 --timeout "30"

  gcloud compute backend-services add-backend "concourse" --instance-group "concourse-us-central1-f" --zone "us-central1-f" --balancing-mode "UTILIZATION" --capacity-scaler "1" --max-utilization "0.8"
  ```

1. Create a URL Map:

  ```
  gcloud compute url-maps create concourse-http --default-service concourse
  ```

1. Create a target proxy:
  ```
  gcloud compute target-http-proxies create concourse-http --url-map concourse-http
  ```
1. Create a global forwarding rule

  ```
  gcloud compute forwarding-rules create concourse-http-fw --target-http-proxy concourse-http --global --address=REPLACE --port-range=80-80
  ```

### Deploying Concourse

* Target and login into your BOSH environment:

```
$ bosh target <YOUR BOSH IP ADDRESS>
```

Your username is `admin` and password is `admin`.

* Upload the required [Google BOSH Stemcell](http://bosh.io/docs/stemcell.html):

```
$ bosh upload stemcell https://storage.googleapis.com/bosh-cpi-artifacts/o/light-bosh-stemcell-3218-google-kvm-ubuntu-trusty-go_agent.tgz
```

* Upload the required [BOSH Releases](http://bosh.io/docs/release.html):

```
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.337.0
$ bosh upload release https://bosh.io/d/github.com/concourse/concourse?v=1.0.0
```

* Download the [concourse.yml](https://raw.githubusercontent.com/cloudfoundry-incubator/bosh-google-cpi-release/master/docs/concourse.yml) deployment manifest file and update it with your properties (at the top of the file):
    - `director_uuid = 'CHANGE-ME'`: replace `CHANGE-ME` with your BOSH UUID (run `bosh status`)
    - `vip_ip = 'CHANGE-ME'`: replace `CHANGE-ME` with the static IP reserved previously (named `concourse`)

* Target the deployment file and deploy:

```
$ bosh deployment concourse.yml
$ bosh deploy
```

* Follow the [Concourse Getting Started](http://concourse.ci/using-concourse.html) guide.
