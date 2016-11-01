# Deploying Cloud Foundry on Google Compute Engine

This guide describes how to deploy a minimal [Cloud Foundry](https://www.cloudfoundry.org/) on [Google Compute Engine](https://cloud.google.com/) using BOSH. The BOSH director must have been created following the steps in the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide.

## Prerequisites

* You must have an existing BOSH director and bastion host created by following the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide.

* Ensure that you have enough [Resource Quotas](https://cloud.google.com/compute/docs/resource-quotas) available:
    - 25 Cores
    - 2 IP addresses
    - 1 Tb persistent disk

## Deploy supporting infrastructure automatically

The following instructions use [Terraform](terraform.io) to provision all of the infrastructure required to run CloudFoundry.

### Steps to perform in `bosh-bastion`

1. SSH to the `bosh-bastion` VM. You can SSH form Cloud Shell or any workstation that has `gcloud` installed:

  ```
  gcloud compute ssh bosh-bastion
  ```

1. `cd` into the Cloud Foundry docs directory that you cloned when you created the BOSH bastion:

  ```
  cd /share/docs/cloudfoundry
  ```

1. View the Terraform execution plan to see the resources that will be created:

  ```
  terraform plan \
    -var bosh_network=${network} \
    -var projectid=${project_id} \
    -var region=${region} \
    -var zone=${zone}
  ```

1. Create the resources:

  ```
  terraform apply \
    -var bosh_network=${network} \
    -var projectid=${project_id} \
    -var region=${region} \
    -var zone=${zone}
  ```

1. Create a service account and key that will be used by Cloud Foundry VMs:

  ```
  gcloud iam service-accounts create cf-component
  gcloud iam service-accounts keys create ~/cf-component.key.json \
      --iam-account cf-component@${project_id}.iam.gserviceaccount.com
  ```

1. Grant the new service account editor access and logging access to your project:

  ```
  gcloud projects add-iam-policy-binding ${project_id} \
      --member serviceAccount:cf-component@${project_id}.iam.gserviceaccount.com \
      --role "roles/editor" \
      --role "roles/logging.logWriter" \
      --role "roles/logging.configWriter"
  ```


1. Target and login into your BOSH environment:

  ```
  bosh target 10.0.0.6
  ```

  > **Note:** Your username is `admin` and password is `admin`.

1. Export several environment variables:

  ```
  export vip=$(terraform output ip)
  export director=$(bosh status --uuid)
  ```

1. Upload the stemcell:

  ```
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

1. Target the deployment file and deploy:

  ```
  bosh deployment manifest.yml
  bosh deploy
  ```

Once deployed, you can target your Cloud Foundry environment using the [CF CLI](http://docs.cloudfoundry.org/cf-cli/). Your CF API endpoint is `https://api.<YOUR CF IP ADDRESS>.xip.io`, your username is `admin` and your password is `c1oudc0w`.

### Delete resources

From your `bosh-bastion` instance, delete your Cloud Foundry deployment:

  ```
  bosh delete deployment cf
  ```

Then delete the infrastructure you created with terraform:
  ```
  cd /share/docs/cloudfoundry
  terraform destroy \
    -var projectid=${project_id} \
    -var region=${region} \
    -var zone=${zone}
  ```
