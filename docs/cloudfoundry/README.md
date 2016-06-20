# Deploying Cloud Foundry on Google Compute Engine

This guide describes how to deploy [Cloud Foundry](https://www.cloudfoundry.org/) on [Google Compute Engine](https://cloud.google.com/) using BOSH. The BOSH director must have been created following the steps in the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide. You must use the same deployment steps here as you did in the BOSH instructions (that is, either __deploy automatically__ or __deploy manually__).


## Prerequisites

* You must have an existing BOSH director and bastion host created by following the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide.

* Ensure that you have enough [Resource Quotas](https://cloud.google.com/compute/docs/resource-quotas) available:
    - 100 Cores
    - 25 IP addresses
    - 1 Tb persistent disk

<a name="deploy-automatic"></a>
## Deploy supporting infrastructure automatically

> **Note:** Use this section if you also used the automatic steps in the BOSH deployment instructions.

The following instructions offer the fastest path to getting Cloud Foundry up and running on Google Cloud Platform. Using [Terraform](terraform.io), you will provision all of the infrastructure required to run CloudFoundry injust a few commands.

If you would like a detailed understanding of what is being created, you may
follow the instructions in [Deploy supporting infrastructure manually](#deploy-manual).

### Requirements
You must have followed the [Deploy supporting infrastructure automatically](../bosh/README.md#deploy-automatic) steps in the Deploy BOSH on Google Cloud Platform, which used `terraform` to deploy the BOSH director.

### Steps
1. Download the Cloud Foundry Terraform file - [cloudfoundry.tf](cloudfoundry.tf) - and save it to the same directory where you saved the BOSH director's `main.tf` file.
1. Use Terraform's `plan` feature to confirm that the new resources will be created:

  ```
  $ tf plan
  ```

1. Create the resources

  ```
  $ tf apply
  ```

Now you have the infrastructure ready to deploy Cloud Foundry. Go ahead to the [Deploy Cloud Foundry](#deploy-cloudfoundry) section to do that. 

<a name="deploy-manually"></a>
## Deploy supporting infrastructure manually

> **Note:** Use this section if you also used the manual steps in the BOSH deployment instructions.

> **Note:** Do not follow these steps if you've already completed the [Deploy supporting infrastructure automatically](#deploy-automatic) instructions. Instead, skip ahead to the [Deploy Cloud Foundry](#deploy-cloudfoundry) section.

The following instructions offer the most detailed path to getting Cloud Foundry up and running on Google Cloud Platform using the `gcloud` CLI. Although it is recommended you use the [Deploy supporting infrastructure automatically](#deploy-automatic) instructions for a production environment, these steps can be helpful in understanding exactly what is being provisioned to support your Cloud Foundry environment.

### Requirements
You must have the `gcloud` CLI installed on your workstation. See
[https://cloud.google.com/sdk/](https://cloud.google.com/sdk/) for more details.

1. Create a new subnetwork for public CloudFoundry components:

  ```
  $ gcloud compute networks subnets create cf-public \
      --network cf \
      --range 10.200.0.0/16 \
      --description "Subnet for public CloudFoundry components" \
      --region us-central1
  ```

1. Create a new subnetwork for private CloudFoundry components:

  ```
  $ gcloud compute networks subnets create cf-private \
      --network cf \
      --range 192.168.0.0/16 \
      --description "Subnet for private CloudFoundry components" \
      --region us-central1
  ```

1. Reserve a new [static external IP address](https://cloud.google.com/compute/docs/instances-and-network#reserve_new_static):

  ```
  $ gcloud compute addresses create cf
  ```

1. Capture the address in a variable:

  ```
  $ address=`gcloud compute addresses describe cf | grep ^address: | cut -f2 -d' '`
  ```

1. Create the following load balancing [health check](https://cloud.google.com/compute/docs/load-balancing/health-checks):

  ```
  $ gcloud compute http-health-checks create cf-public \
    --description "Cloud Foundry Public Health Check" \
    --timeout "5s" \
    --check-interval "30s" \
    --healthy-threshold "10" \
    --unhealthy-threshold "2" \
    --port 80 \
    --request-path "/info" \
    --host "api.${address}.xip.io"
  ```

1. Create the following load balancing [target pool](https://cloud.google.com/compute/docs/load-balancing/network/target-pools):

  ```
  $ gcloud compute target-pools create cf-public \
    --description "Cloud Foundry Public Target Pool" \
    --health-check cf-public
  ```

1. Create the following load balancing [forwarding rules](https://cloud.google.com/compute/docs/load-balancing/network/forwarding-rules):

  ```
  $ gcloud compute forwarding-rules create cf-http \
    --description "Cloud Foundry HTTP Traffic" \
    --ip-protocol TCP \
    --port-range 80 \
    --target-pool cf-public \
    --address ${address}
  ```

  ```
  $ gcloud compute forwarding-rules create cf-https \
    --description "Cloud Foundry HTTPS Traffic" \
    --ip-protocol TCP \
    --port-range 443 \
    --target-pool cf-public \
    --address ${address}
  ```

  ```
  $ gcloud compute forwarding-rules create cf-ssh \
    --description "Cloud Foundry SSH Traffic" \
    --ip-protocol TCP \
    --port-range 2222 \
    --target-pool cf-public \
    --address ${address}
  ```

  ```
  $ gcloud compute forwarding-rules create cf-wss \
    --description "Cloud Foundry WSS Traffic" \
    --ip-protocol TCP \
    --port-range 4443 \
    --target-pool cf-public \
    --address ${address}
  ```

1. Create the following firewalls and [set the appropriate rules](https://cloud.google.com/compute/docs/networking#addingafirewall):

  ```
  $ gcloud compute firewall-rules create cf-public \
    --description "Cloud Foundry public traffic" \
    --network cf \
    --target-tags cf-public \
    --allow tcp:80,tcp:443,tcp:2222,tcp:4443
  ```

<a name="deploy-cloudfoundry"></a>
### Deploy Cloud Foundry

1. Target and login into your BOSH environment:

  ```
  $ bosh target <YOUR BOSH IP ADDRESS>
  ```

  > **Note:** Your username is `admin` and password is `admin`.

1. Upload the required [Google BOSH Stemcell](http://bosh.io/docs/stemcell.html):

  ```
  $ bosh upload stemcell https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3218-google-kvm-ubuntu-trusty-go_agent.tgz
  ```

1. Upload the required [BOSH Releases](http://bosh.io/docs/release.html):

  ```
  $ bosh upload release https://bosh.io/d/github.com/cloudfoundry/cf-mysql-release?v=23
  $ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.333.0
  $ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/etcd-release?v=36
  $ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/diego-release?v=0.1454.0
  $ bosh upload release https://bosh.io/d/github.com/cloudfoundry/cf-release?v=231
  ```

1. Download the [cloudfoundry.yml](cloudfoundry.yml) deployment manifest file and update it with your properties (at the top of the file):
    - `director_uuid = '{{DIRECTOR_UUID}}'`: replace `{{DIRECTOR_UUID}}` with your BOSH UUID (run `bosh status --uuid`)
    - `vip_ip = 'VIP_IP'`: replace `VIP_IP` with the static IP reserved previously (named `cf`)

    * Alternatively, use `sed` to replace these values automatically:

      ```
      $ sed -i s#{{VIP_IP}}#`gcloud compute addresses describe cf | grep ^address: | cut -f2 -d' '`# cloudfoundry.yml
      $ sed -i s#{{DIRECTOR_UUID}}#`bosh status --uuid 2>/dev/null`# cloudfoundry.yml
      ```

1. Target the deployment file and deploy:

  ```
  $ bosh deployment cloudfoundry.yml
  $ bosh deploy
  ```

Once deployed, you can target your Cloud Foundry environment using the [CF CLI](http://docs.cloudfoundry.org/cf-cli/). Your CF API endpoint is `https://api.<YOUR CF IP ADDRESS>.xip.io`, your username is `admin` and your password is `c1oudc0w`.
