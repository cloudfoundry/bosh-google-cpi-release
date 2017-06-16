variable "env_name" {
  type = "string"
}
variable "google_project" {
  type = "string"
}
variable "google_region" {
  default = "us-central1"
}
variable "google_zone" {
  default = "us-central1-a"
}
variable "google_json_key_data" {
  type = "string"
}

provider "google" {
  credentials = "${var.google_json_key_data}"
  project     = "${var.google_project}"
  region      = "${var.google_region}"
}

resource "google_service_account" "google_service_account" {
  account_id   = "${var.env_name}"
}

data "google_iam_policy" "service_account_actor" {
  binding {
    role = "roles/iam.serviceAccountActor"

    members = [
      "serviceAccount:${google_service_account.google_service_account.email}",
    ]
  }
}

resource "google_compute_address" "google_address_director_ubuntu" {
  name = "${var.env_name}-director-ubuntu"
}

resource "google_compute_address" "google_address_bats_ubuntu" {
  name = "${var.env_name}-bats-ubuntu"
}

resource "google_compute_address" "google_address_int_ubuntu" {
  name = "${var.env_name}-int-ubuntu"
}

resource "google_compute_network" "google_auto_network" {
  name = "${var.env_name}-auto"
  auto_create_subnetworks = true
}

resource "google_compute_network" "google_custom_network" {
  name = "${var.env_name}-custom"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "google_subnetwork" {
  name          = "${var.env_name}-${var.google_region}"
  ip_cidr_range = "10.0.0.0/24"
  network       = "${google_compute_network.google_custom_network.self_link}"
}

resource "google_compute_firewall" "google_firewall_internal" {
  name    = "${var.env_name}-internal"
  network = "${google_compute_network.google_custom_network.name}"

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

  source_tags = ["${var.env_name}-internal"]
  target_tags = ["${var.env_name}-internal"]
}

resource "google_compute_firewall" "google_firewall_external" {
  name    = "${var.env_name}-external"
  network = "${google_compute_network.google_custom_network.name}"

  description = "BOSH CI External traffic"

  allow {
    protocol = "tcp"
    ports = ["22", "443", "4222", "6868", "25250", "25555", "25777"]
  }
  allow {
    protocol = "udp"
    ports = ["53"]
  }

  target_tags = ["${var.env_name}-external"]
}

resource "google_compute_target_pool" "google_target_pool" {
  name = "${var.env_name}"
}

resource "google_compute_instance_group" "google_backend_service" {
  name = "${var.env_name}"
  zone = "${var.google_zone}"
}

resource "google_compute_http_health_check" "google_backend_service" {
  name = "${var.env_name}"
}

resource "google_compute_backend_service" "google_backend_service" {
  name        = "${var.env_name}"
  port_name   = "http"
  timeout_sec = 30

  backend {
    group = "${google_compute_instance_group.google_backend_service.self_link}"
  }

  health_checks = ["${google_compute_http_health_check.google_backend_service.self_link}"]
}

resource "google_compute_instance_group" "google_region_backend_service" {
  name = "${var.env_name}-region"
  zone = "${var.google_zone}"

  instances = ["${google_compute_instance.google_region_backend_service.self_link}"]
}

resource "google_compute_health_check" "google_region_backend_service" {
  name         = "${var.env_name}"

  tcp_health_check {}
}

resource "google_compute_instance" "google_region_backend_service" {
  name         = "${var.env_name}"
  zone         = "${var.google_zone}"
  machine_type = "f1-micro"

  network_interface {
    subnetwork = "${google_compute_subnetwork.google_subnetwork.name}"

    address = "${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, -3)}"
  }

  disk {
    image = "debian-cloud/debian-8"
  }

  depends_on = ["google_compute_subnetwork.google_subnetwork"]
}

resource "google_compute_region_backend_service" "google_region_backend_service" {
  name        = "${var.env_name}"
  protocol    = "TCP"
  timeout_sec = 30

  backend {
    group = "${google_compute_instance_group.google_region_backend_service.self_link}"
  }

  health_checks = ["${google_compute_health_check.google_region_backend_service.self_link}"]
}

output "ProjectID" {
  value = "${var.google_project}"
}
output "Region" {
  value = "${var.google_region}"
}
output "Zone" {
  value = "${var.google_zone}"
}
output "AutoNetwork" {
  value = "${google_compute_network.google_auto_network.name}"
}
output "CustomNetwork" {
  value = "${google_compute_network.google_custom_network.name}"
}
output "Subnetwork" {
  value = "${google_compute_subnetwork.google_subnetwork.name}"
}
output "SubnetworkCIDR" {
  value = "${google_compute_subnetwork.google_subnetwork.ip_cidr_range}"
}
output "InternalTag" {
  value = "${google_compute_firewall.google_firewall_internal.name}"
}
output "ExternalTag" {
  value = "${google_compute_firewall.google_firewall_external.name}"
}
output "DirectorExternalIP" {
  value = "${google_compute_address.google_address_director_ubuntu.address}"
}
output "DirectorInternalIP" {
  value = "${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 6)}"
}
output "SubnetworkGateway" {
  value = "${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 1)}"
}
output "ReservedRange" {
  value = "${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 2)}-${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 15)}"
}
output "StaticRange" {
  value = "${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 16)}-${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 40)}"
}
output "IntegrationStaticIPs" {
  value = "${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 10)},${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 11)},${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 12)}"
}
output "IntegrationExternalIP" {
  value = "${google_compute_address.google_address_int_ubuntu.address}"
}
output "TargetPool" {
  value = "${google_compute_target_pool.google_target_pool.name}"
}
output "BackendService" {
  value = "${google_compute_backend_service.google_backend_service.name}"
}
output "RegionBackendService" {
  value = "${google_compute_region_backend_service.google_region_backend_service.name}"
}
output "ILBInstanceGroup" {
  value = "${google_compute_instance_group.google_region_backend_service.name}"
}
output "ServiceAccount" {
  value = "${google_service_account.google_service_account.email}"
}
output "BATsExternalIP" {
  value = "${google_compute_address.google_address_bats_ubuntu.address}"
}
output "BATsStaticIPPair" {
  value = "${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 13)},${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 14)}"
}
output "BATsStaticIP" {
  value = "${cidrhost(google_compute_subnetwork.google_subnetwork.ip_cidr_range, 7)}"
}
