variable "project" {
  type = "string"
}

variable "google_credentials" {
  type = "string"
}

variable "region" {
  default = "us-central1"
}

variable "zone_a" {
  default = "us-central1-a"
}

variable "zone_b" {
  default = "us-central1-b"
}

variable "env_name" {
  type = "string"
}

provider "google" {
  credentials = "${file(var.google_credentials)}"
  project     = "${var.project}"
  region      = "${var.region}"
}

// If you see `Error creating service account: googleapi: Error 403: Google Identity and Access Management (IAM) API has not been used in project`
// then follow the instructions to fix
// gcloud --project=$(GOOGLE_PROJECT) iam service-accounts create cfintegration
resource "google_service_account" "cfintegration" {
  account_id   = "${var.env_name}"
}

resource "google_compute_network" "cfintegration" {
  name = "${var.env_name}"
  auto_create_subnetworks = "true"
}

resource "google_compute_network" "cfintegration-custom" {
  name = "${var.env_name}-custom"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "cfintegration-custom-subnet" {
  name          = "${var.env_name}-custom-subnet"
  ip_cidr_range = "192.168.0.0/16"
  network       = "${var.env_name}-custom"
  region        = "${var.region}"
  depends_on    = [ "google_compute_network.cfintegration-custom" ]
}

resource "google_compute_address" "cfintegration" {
  name       = "${var.env_name}"
  region     = "${var.region}"
  depends_on = [ "google_compute_network.cfintegration" ]
}

resource "google_compute_target_pool" "cfintegration" {
  name = "${var.env_name}"
}

resource "google_compute_instance_group" "cfintegration_a" {
  name        = "${var.env_name}"
  zone        = "${var.zone_a}"
}

resource "google_compute_instance_group" "cfintegration_b" {
  name        = "${var.env_name}"
  zone        = "${var.zone_b}"
}

resource "google_compute_http_health_check" "cfintegration" {
  name         = "${var.env_name}"
}

resource "google_compute_backend_service" "cfintegration" {
  name        = "${var.env_name}"
  port_name   = "http"
  timeout_sec = 30

  backend {
    group = "${google_compute_instance_group.cfintegration_a.self_link}"
    balancing_mode = "UTILIZATION"
    capacity_scaler = 1 # default 1
    max_utilization = 0.8 # default 0.8
  }

  backend {
    group = "${google_compute_instance_group.cfintegration_b.self_link}"
    balancing_mode = "UTILIZATION"
    capacity_scaler = 1 # default 1
    max_utilization = 0.8 # default 0.8
  }

  health_checks = [ "${google_compute_http_health_check.cfintegration.self_link}" ]
  depends_on    = [
    "google_compute_http_health_check.cfintegration",
    "google_compute_instance_group.cfintegration_a",
    "google_compute_instance_group.cfintegration_b"
  ]
}

resource "google_compute_instance_group" "cfintegration-ilb-a" {
  name        = "${var.env_name}-ilb"
  zone        = "${var.zone_a}"

  instances = [
    "${google_compute_instance.cfintegration-ilb-a.self_link}"
  ]
}

resource "google_compute_instance_group" "cfintegration-ilb-b" {
  name        = "${var.env_name}-ilb"
  zone        = "${var.zone_b}"

  instances = [
    "${google_compute_instance.cfintegration-ilb-b.self_link}"
  ]
}

resource "google_compute_health_check" "cfintegration" {
  name = "${var.env_name}"
  tcp_health_check { }
}

resource "google_compute_instance" "cfintegration-ilb-a" {
  name         = "${var.env_name}-ilb-a"
  zone         = "${var.zone_a}"
  machine_type = "f1-micro"

  network_interface {
    subnetwork = "${google_compute_subnetwork.cfintegration-custom-subnet.name}"
  }

  disk {
    image = "debian-cloud/debian-8"
  }
}

resource "google_compute_instance" "cfintegration-ilb-b" {
  name         = "${var.env_name}-ilb-b"
  zone         = "${var.zone_b}"
  machine_type = "f1-micro"

  network_interface {
    subnetwork = "${google_compute_subnetwork.cfintegration-custom-subnet.name}"
  }

  disk {
    image = "debian-cloud/debian-8"
  }
}

resource "google_compute_region_backend_service" "region-cfintegration" {
  name             = "region-${var.env_name}"
  protocol         = "TCP"
  timeout_sec      = 30

  backend {
    group = "${google_compute_instance_group.cfintegration-ilb-a.self_link}"
  }
  backend {
    group = "${google_compute_instance_group.cfintegration-ilb-b.self_link}"
  }

  health_checks = [ "${google_compute_health_check.cfintegration.self_link}" ]
}

output "external_ip" {
  value  = "${google_compute_address.cfintegration.address}"
}
