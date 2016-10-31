# Deploying Cloud Foundry on Google Compute Engine

This guide describes how to deploy a minimal [Cloud Foundry](https://www.cloudfoundry.org/) on [Google Compute Engine](https://cloud.google.com/) using BOSH. The BOSH director must have been created following the steps in the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide.


## Prerequisites

* You must have an existing BOSH director and bastion host created by following the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide.

* Ensure that you have enough [Resource Quotas](https://cloud.google.com/compute/docs/resource-quotas) available:
    - 25 Cores
    - 2 IP addresses
    - 1 Tb persistent disk

## Deploy supporting infrastructure automatically

The following instructions offer the fastest path to getting Cloud Foundry up and running on Google Cloud Platform. Using [Terraform](terraform.io), you will provision all of the infrastructure required to run CloudFoundry injust a few commands.

### Requirements
You must have followed the [Deploy supporting infrastructure automatically](../bosh/README.md#deploy-automatic) steps in the Deploy BOSH on Google Cloud Platform, which used `terraform` to deploy the BOSH director.

### Steps
1. Download the Cloud Foundry Terraform file - [cloudfoundry.tf](cloudfoundry.tf) - and save it to the same directory where you saved the BOSH director's `main.tf` file.

1. Export your project id, preferred region, and zone (you may skip this if already set from the previous BOSH deployment instructions)

  ```
  export projectid=REPLACE_WITH_YOUR_PROJECT_ID
  export region=us-east1
  export zone=us-east1-d
  ```

1. Use Terraform's `plan` feature to confirm that the new resources will be created:

  ```
  terraform plan -var projectid=${projectid} -var region=${region} -var zone=${zone}
  ```

1. Create the resources

  ```
  terraform apply -var projectid=${projectid} -var region=${region} -var zone=${zone}
  ```

1. Create a service account and key:

  ```
  gcloud iam service-accounts create cf-component
  gcloud iam service-accounts keys create /tmp/cf-component.key.json \
      --iam-account cf-component@${projectid}.iam.gserviceaccount.com
  ```

1. Grant the new service account editor access and logging access to your project:

  ```
  gcloud projects add-iam-policy-binding ${projectid} \
      --member serviceAccount:cf-component@${projectid}.iam.gserviceaccount.com \
      --role "roles/editor" \
      --role "roles/logging.logWriter" \
      --role "roles/logging.configWriter"
  ```


Now you have the infrastructure ready to deploy Cloud Foundry. Go ahead to the [Deploy Cloud Foundry](#deploy-cloudfoundry) section to do that. 

<a name="deploy-cloudfoundry"></a>
## Deploy Cloud Foundry
Before working this section, you must have deployed the supporting infrastructure on Google Cloud Platform using the [automatic](#deploy-automatic) steps provided earlier.

1. SSH to your bastion instance:

  ```
  gcloud compute ssh bosh-bastion
  ```

1. Export the region and zone on the bastion host:

  ```
  zone=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone)
  zone=${zone##*/}
  region=${zone%-*}
  ```

1. Target and login into your BOSH environment:

  ```
  bosh target <YOUR BOSH IP ADDRESS>
  ```

  > **Note:** Your username is `admin` and password is `admin`.

  bosh upload stemcell https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent?v=3263.7
  ```

1. Upload the required [BOSH Releases](http://bosh.io/docs/release.html):

  ```
  bosh upload release https://bosh.io/d/github.com/cloudfoundry/cf-mysql-release?v=23
  bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.340.0
  bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/etcd-release?v=36
  bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/diego-release?v=0.1454.0
  bosh upload release https://bosh.io/d/github.com/cloudfoundry/cf-release?v=231
  ```

1. Download the [cloudfoundry.yml](cloudfoundry.yml) deployment manifest file and use `sed` to modify a few values in it:

  ```
  sed -i s#{{VIP_IP}}#`gcloud compute addresses describe cf | grep ^address: | cut -f2 -d' '`# cloudfoundry.yml
  sed -i s#{{DIRECTOR_UUID}}#`bosh status --uuid 2>/dev/null`# cloudfoundry.yml
  sed -i s#{{REGION}}#$region# cloudfoundry.yml
  sed -i s#{{ZONE}}#$zone# cloudfoundry.yml
  sed -i s#{{PROJECT_ID}}#$projectid# cloudfoundry.yml
  ```

1. Target the deployment file and deploy:

  ```
  bosh deployment cloudfoundry.yml
  bosh deploy
  ```

Once deployed, you can target your Cloud Foundry environment using the [CF CLI](http://docs.cloudfoundry.org/cf-cli/). Your CF API endpoint is `https://api.<YOUR CF IP ADDRESS>.xip.io`, your username is `admin` and your password is `c1oudc0w`.
