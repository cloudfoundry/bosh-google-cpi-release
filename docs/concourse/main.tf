variable "projectid" {
    type = "string"
    default = "REPLACE-WITH-YOUR-GOOGLE-PROJECT-ID"
}

variable "region" {
    type = "string"
    default = "us-east1"
}

variable "zone-1" {
    type = "string"
    default = "us-east1-d"
}

variable "zone-2" {
    type = "string"
    default = "us-east1-b"
}

variable "name" {
    type = "string"
    default = "concourse"
}

provider "google" {
    project = "${var.projectid}"
    region = "${var.region}"
}

resource "google_compute_network" "network" {
  name       = "${var.name}"
}

// Subnet for the BOSH director
resource "google_compute_subnetwork" "bosh-subnet-1" {
  name          = "bosh-${var.name}-${var.region}"
  ip_cidr_range = "10.0.0.0/24"
  network       = "${google_compute_network.network.self_link}"
}

// Allow SSH to BOSH bastion
resource "google_compute_firewall" "bosh-bastion" {
  name    = "bosh-bastion-${var.name}"
  network = "${google_compute_network.network.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  target_tags = ["bosh-bastion"]
}

// Allow open access between internal MVs
resource "google_compute_firewall" "bosh-internal" {
  name    = "bosh-internal-${var.name}"
  network = "${google_compute_network.network.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
  }

  allow {
    protocol = "udp"
  }
  target_tags = ["bosh-internal"]
  source_tags = ["bosh-internal"]
}

// BOSH bastion host
resource "google_compute_instance" "bosh-bastion" {
  name         = "bosh-bastion-${var.name}"
  machine_type = "n1-standard-1"
  zone         = "${var.zone-1}"

  tags = ["bosh-bastion", "bosh-internal"]

  boot_disk {
   initialize_params {
     image = "ubuntu-1404-trusty-v20180122"
   }
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.bosh-subnet-1.name}"
    access_config {
      // Ephemeral IP
    }
  }

  metadata_startup_script = <<EOT
#!/bin/bash
apt-get update -y
apt-get install -y build-essential zlibc zlib1g-dev ruby ruby-dev openssl libxslt-dev libxml2-dev libssl-dev libreadline6 libreadline6-dev libyaml-dev libsqlite3-dev sqlite3
gem install bosh_cli
curl -o /tmp/cf.tgz https://s3.amazonaws.com/go-cli/releases/v6.20.0/cf-cli_6.20.0_linux_x86-64.tgz
tar -zxvf /tmp/cf.tgz && mv cf /usr/bin/cf && chmod +x /usr/bin/cf
curl -o /usr/bin/bosh-init https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.96-linux-amd64
chmod +x /usr/bin/bosh-init
EOT

  service_account {
    scopes = ["cloud-platform"]
  }
}
