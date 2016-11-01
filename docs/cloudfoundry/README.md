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
  cd /share/bosh-google-cpi-release/docs/cloudfoundry
  ```

1. View the Terraform execution plan to see the resources that will be created:

  ```
  docker run -i -t \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v `pwd`:/$(basename `pwd`) \
    -w /$(basename `pwd`) \
    hashicorp/terraform:light plan \
      -var projectid=${project_id} \
      -var region=${region} \
      -var zone=${zone}
  ```

1. Create the resources:

  ```
  docker run -i -t \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v `pwd`:/$(basename `pwd`) \
    -w /$(basename `pwd`) \
    hashicorp/terraform:light apply \
      -var projectid=${project_id} \
      -var region=${region} \
      -var zone=${zone}
  ```

1. Create a service account and key that will be used by Cloud Foundry VMs:

  ```
  gcloud iam service-accounts create cf
  gcloud iam service-accounts keys create ~/cf.key.json \
      --iam-account cf@${projectid}.iam.gserviceaccount.com
  ```

1. Grant the new service account editor access and logging access to your project:

  ```
  gcloud projects add-iam-policy-binding ${projectid} \
      --member serviceAccount:cf@${projectid}.iam.gserviceaccount.com \
      --role "roles/editor" \
      --role "roles/logging.logWriter" \
      --role "roles/logging.configWriter"
  ```


1. Target and login into your BOSH environment:

  ```
  bosh target <YOUR BOSH IP ADDRESS>
  ```

  > **Note:** Your username is `admin` and password is `admin`.

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
