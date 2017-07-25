# Deploying Cloud Foundry on Google Compute Engine

This guide describes how to deploy a minimal [Cloud Foundry](https://www.cloudfoundry.org/) on [Google Compute Engine](https://cloud.google.com/) using BOSH. The BOSH director must have been created following the steps in the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide.

## Prerequisites

* You must have an existing BOSH director and bastion host created by following the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide.

## Cost and resource requirements

* This CF deployment (and the BOSH director/bastion pre-requisite) is small enough to fit in a default project [Resource Quota](https://cloud.google.com/compute/docs/resource-quotas). It consumes:
    - 23 cores in a single region
    - 24 pre-emptible cores in a second region **during compilation only**
    - 2 IP addresses
    - 660 Gb persistent disk

You can view an estimate of the cost to run this deployment for an entire month at [this link](https://cloud.google.com/products/calculator/#id=8de9b03b-79b9-4b0a-9d26-01dc4b40937f).

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

1. Export a few vars to specify the location of compilation VMs:

    ```
    # You may be tempted to set these to the same value as your BOSH deployment (eg `us-east1`).
    # However, this may cause you to exceed your regional quotas.
    # By breaking apart the regions, you avoid this problem.

    export region_compilation=us-central1
    export zone_compilation=us-central1-b
    ```

1. View the Terraform execution plan to see the resources that will be created:

    ```
    terraform plan \
      -var network=${network} \
      -var project_id=${project_id} \
      -var network_project_id=${network_project_id} \
      -var region=${region} \
      -var region_compilation=${region_compilation} \
      -var zone=${zone} \
      -var zone_compilation=${zone_compilation}
    ```

1. Create the resources:

    ```
    terraform apply \
      -var network=${network} \
      -var project_id=${project_id} \
      -var network_project_id=${network_project_id} \
      -var region=${region} \
      -var region_compilation=${region_compilation} \
      -var zone=${zone} \
      -var zone_compilation=${zone_compilation}
    ```

1. Create a service account and key that will be used by Cloud Foundry VMs:

    ```
    gcloud iam service-accounts create cf-component
    ```

1. Grant the new service account editor access and logging access to your project:

    ```
    gcloud projects add-iam-policy-binding ${project_id} \
        --member serviceAccount:cf-component@${project_id}.iam.gserviceaccount.com \
        --role "roles/editor"
      
    gcloud projects add-iam-policy-binding ${project_id} \
        --member serviceAccount:cf-component@${project_id}.iam.gserviceaccount.com \
        --role "roles/logging.logWriter"
      
    gcloud projects add-iam-policy-binding ${project_id} \
        --member serviceAccount:cf-component@${project_id}.iam.gserviceaccount.com \
        --role "roles/logging.configWriter"
    ```

1. Target and login into your BOSH environment:

    ```
    bosh2 alias-env micro-google --environment 10.0.0.6 --ca-cert ../ca_cert.pem
    bosh2 login -e micro-google
    ```

    > **Note:** Your username is `admin` and password is `admin`.

1. Retrieve the outputs of your Terraform run to be used in your Cloud Foundry deployment:

    ```
    export vip=$(terraform output ip)
    export tcp_vip=$(terraform output tcp_ip)
    export zone=$(terraform output zone)
    export zone_compilation=$(terraform output zone_compilation)
    export region=$(terraform output region)
    export region_compilation=$(terraform output region_compilation)
    export private_subnet=$(terraform output private_subnet)
    export compilation_subnet=$(terraform output compilation_subnet)
    export network=$(terraform output network)

    export director=$(bosh2 env -e micro-google | sed -n 2p)
    ```

1. Upload the stemcell:

    ```
    bosh2 upload-stemcell -e micro-google https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent?v=3312.15
    ```

1. Upload the required [BOSH Releases](http://bosh.io/docs/release.html):

    ```
    bosh2 upload-release -e micro-google  https://bosh.io/d/github.com/cloudfoundry/cf-mysql-release?v=23
    bosh2 upload-release -e micro-google  https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.340.0
    bosh2 upload-release -e micro-google  https://bosh.io/d/github.com/cloudfoundry-incubator/etcd-release?v=43
    bosh2 upload-release -e micro-google  https://bosh.io/d/github.com/cloudfoundry-incubator/diego-release?v=0.1463.0
    bosh2 upload-release -e micro-google  https://bosh.io/d/github.com/cloudfoundry/cf-release?v=249
    bosh2 upload-release -e micro-google  https://bosh.io/d/github.com/cloudfoundry-incubator/cf-routing-release?v=0.142.0
    ```

1. Use `erb` to substitute variables in the template:

    ```
    erb manifest.yml.erb > manifest.yml
    ```

1. Target the deployment file and deploy:

    ```
    bosh2 -e micro-google deploy -d cf manifest.yml
    ```

    > **Note:** If package compilation fails, consider disabling VM preemption in
`manifest.yml.erb` under `compilation/cloud_properties/preemptible`


Once deployed, you can target your Cloud Foundry environment using the [CF CLI](http://docs.cloudfoundry.org/cf-cli/):

  ```
  cf api https://api.${vip}.xip.io --skip-ssl-validation
  cf login
  ```

Your username is `admin` and your password is `c1oudc0w`.

### Optional: Setup TCP routing
Setup your tcp router domain using the CLI:
    ```
    cf create-shared-domain tcp.${tcp_vip}.xip.io --router-group default-tcp
    cf update-quota default --reserved-route-ports 10
    ```

### Delete resources

From your `bosh-bastion` instance, delete your Cloud Foundry deployment:

```
bosh2 -e micro-google delete-deployment -d cf
```

Then delete the infrastructure you created with terraform:
```
cd /share/docs/cloudfoundry
terraform destroy \
    -var project_id=${project_id} \
    -var network_project_id=${network_project_id} \
    -var region=${region} \
    -var zone=${zone} \
    -var network=${network}
```

**Important:** The BOSH bastion and director you created must also be destroyed. Follow the **Delete resources** instructions in the [Deploy BOSH on Google Cloud Platform](../bosh/README.md) guide.
