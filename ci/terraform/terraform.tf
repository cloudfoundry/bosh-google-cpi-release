variable "google_project" {
  type = "string"
}
variable "google_region" {
  type = "string"
}
variable "google_zone" {
  type = "string"
}
variable "google_json_key_data" {
  type = "string"
}
variable "google_network" {
  type = "string"
}
variable "google_auto_network" {
  type = "string"
}
variable "google_subnetwork" {
  type = "string"
}
variable "google_subnetwork_range" {
  type = "string"
}
variable "google_firewall_internal" {
  type = "string"
}
variable "google_firewall_external" {
  type = "string"
}
variable "google_address_director_ubuntu" {
  type = "string"
}
variable "google_address_bats_ubuntu" {
  type = "string"
}
variable "google_address_int_ubuntu" {
  type = "string"
}
variable "google_target_pool" {
  type = "string"
}
variable "google_backend_service" {
  type = "string"
}
variable "google_region_backend_service" {
  type = "string"
}
variable "google_service_account" {
  type = "string"
}

provider "google" {
  credentials = "${var.google_json_key_data}"
  project     = "${var.google_project}"
  region      = "${var.google_region}"
}

resource "google_service_account" "google_service_account" {
  account_id   = "${var.google_service_account}"
}

resource "google_compute_address" "google_address_director_ubuntu" {
  name = "${var.google_address_director_ubuntu}"
}

resource "google_compute_address" "google_address_bats_ubuntu" {
  name = "${var.google_address_bats_ubuntu}"
}

resource "google_compute_address" "google_address_int_ubuntu" {
  name = "${var.google_address_int_ubuntu}"
}

resource "google_compute_network" "google_auto_network" {
  name = "${var.google_auto_network}"
}

resource "google_compute_network" "google_network" {
  name = "${var.google_network}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "google_subnetwork" {
  name          = "${var.google_subnetwork}"
  ip_cidr_range = "${var.google_subnetwork_range}"
  network       = "${google_compute_network.google_network.self_link}"
}

resource "google_compute_firewall" "google_firewall_internal" {
  name    = "${var.google_firewall_internal}"
  network = "${google_compute_network.google_network.name}"

  description = "BOSH CI Internal traffic"

  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
  }
  allow {
    protocol = "udp"
  }

  source_tags = ["${var.google_firewall_internal}"]
  target_tags = ["${var.google_firewall_internal}"]
}

resource "google_compute_firewall" "google_firewall_external" {
  name    = "${var.google_firewall_external}"
  network = "${google_compute_network.google_network.name}"

  description = "BOSH CI External traffic"

  allow {
    protocol = "tcp"
    ports = ["22", "443", "4222", "6868", "25250", "25555", "25777"]
  }
  allow {
    protocol = "udp"
    ports = ["53"]
  }

  target_tags = ["${var.google_firewall_external}"]
}

resource "google_compute_target_pool" "google_target_pool" {
  name = "${var.google_target_pool}"
}

resource "google_compute_instance_group" "google_backend_service" {
  name = "${var.google_backend_service}"
  zone = "${var.google_zone}"
}

resource "google_compute_http_health_check" "google_backend_service" {
  name         = "${var.google_backend_service}"
}

resource "google_compute_backend_service" "google_backend_service" {
  name        = "${var.google_backend_service}"
  port_name   = "http"
  timeout_sec = 30

  backend {
    group = "${google_compute_instance_group.google_backend_service.self_link}"
  }

  health_checks = ["${google_compute_http_health_check.google_backend_service.self_link}"]
}

resource "google_compute_instance_group" "google_region_backend_service" {
  name = "${var.google_region_backend_service}"
  zone = "${var.google_zone}"

  instances = ["${google_compute_instance.google_region_backend_service.self_link}"]
}

resource "google_compute_health_check" "google_region_backend_service" {
  name         = "${var.google_region_backend_service}"

  tcp_health_check {}
}

resource "google_compute_instance" "google_region_backend_service" {
  name = "${var.google_region_backend_service}"
  zone         = "${var.google_zone}"
  machine_type = "f1-micro"

  network_interface {
    subnetwork = "${google_compute_subnetwork.google_subnetwork.name}"

    address = "${cidrhost(var.google_subnetwork_range, -3)}"
  }

  disk {
    image = "debian-cloud/debian-8"
  }

  depends_on = ["google_compute_subnetwork.google_subnetwork"]
}

resource "google_compute_region_backend_service" "google_region_backend_service" {
  name        = "${var.google_region_backend_service}"
  protocol    = "TCP"
  timeout_sec = 30

  backend {
    group = "${google_compute_instance_group.google_region_backend_service.self_link}"
  }

  health_checks = ["${google_compute_health_check.google_region_backend_service.self_link}"]
}
