variable "project_id" {
    type = "string"
}

variable "network_project_id" {
    type = "string"
}

variable "region" {
    type = "string"
    default = "us-east1"
}

variable "region_compilation" {
    type = "string"
    default = "us-central1"
}

variable "zone" {
    type = "string"
    default = "us-east1-d"
}

variable "zone_compilation" {
    type = "string"
    default = "us-central1-b"
}

variable "network" {
    type = "string"
}

variable "prefix" {
    type = "string"
    default = ""
}

provider "google" {
    credentials = ""
    project = "${var.project_id}"
    region = "${var.region}"
}

// Subnet for the public Cloud Foundry components
resource "google_compute_subnetwork" "cf-compilation-subnet-1" {
  name          = "${var.prefix}cf-compilation-${var.region_compilation}"
  region        = "${var.region_compilation}"
  ip_cidr_range = "10.200.0.0/16"
  network       = "https://www.googleapis.com/compute/v1/projects/${var.network_project_id}/global/networks/${var.network}"
  project       = "${var.network_project_id}"
}

// Subnet for the private Cloud Foundry components
resource "google_compute_subnetwork" "cf-private-subnet-1" {
  name          = "${var.prefix}cf-private-${var.region}"
  ip_cidr_range = "192.168.0.0/16"
  network       = "https://www.googleapis.com/compute/v1/projects/${var.network_project_id}/global/networks/${var.network}"
  project       = "${var.network_project_id}"
}

// Allow access to CloudFoundry HTTP router
resource "google_compute_firewall" "cf-public" {
  name    = "${var.prefix}cf-public"
  network = "${var.network}"
  project = "${var.network_project_id}"

  allow {
    protocol = "tcp"
    ports    = ["80", "443", "2222", "4443"]
  }

  target_tags = ["cf-public"]
}

// Allow access to CloudFoundry TCP router
resource "google_compute_firewall" "cf-tcp-public" {
  name    = "${var.prefix}cf-tcp-public"
  network = "${var.network}"
  project = "${var.network_project_id}"

  allow {
    protocol = "tcp"
    ports    = ["80", "1024-1123"]
  }

  target_tags = ["cf-tcp-public"]
}

// Static IP address for HTTP forwarding rule
resource "google_compute_address" "cf" {
  name = "${var.prefix}cf"
}

// Static IP address for TCP forwarding rule
resource "google_compute_address" "cf-tcp" {
  name = "${var.prefix}cf-tcp"
}

// HTTP Router Health check
resource "google_compute_http_health_check" "cf-public" {
  name                = "${var.prefix}cf-public"
  host                = "api.${google_compute_address.cf.address}.xip.io"
  request_path        = "/info"
  check_interval_sec  = 30 
  timeout_sec         = 5
  healthy_threshold   = 10
  unhealthy_threshold = 2
  port = 80
}

// TCP Router Health check
resource "google_compute_http_health_check" "cf-tcp-public" {
  name                = "${var.prefix}cf-tcp-public"
  request_path        = "/health"
  check_interval_sec  = 30 
  timeout_sec         = 5
  healthy_threshold   = 10
  unhealthy_threshold = 2
  port = 80
}

// HTTP Load balancing target pool
resource "google_compute_target_pool" "cf-public" {
  name = "${var.prefix}cf-public"

  health_checks = [
    "${google_compute_http_health_check.cf-public.name}"
  ]
}

// TCP Router Load balancing target pool
resource "google_compute_target_pool" "cf-tcp-public" {
  name = "${var.prefix}cf-tcp-public"

  health_checks = [
    "${google_compute_http_health_check.cf-tcp-public.name}"
  ]
}

// HTTP forwarding rule
resource "google_compute_forwarding_rule" "cf-http" {
  name        = "${var.prefix}cf-http"
  target      = "${google_compute_target_pool.cf-public.self_link}"
  port_range  = "80"
  ip_protocol = "TCP"
  ip_address  = "${google_compute_address.cf.address}"
}

// HTTPS forwarding rule
resource "google_compute_forwarding_rule" "cf-https" {
  name        = "${var.prefix}cf-https"
  target      = "${google_compute_target_pool.cf-public.self_link}"
  port_range  = "443"
  ip_protocol = "TCP"
  ip_address  = "${google_compute_address.cf.address}"
}

// SSH forwarding rule
resource "google_compute_forwarding_rule" "cf-ssh" {
  name        = "${var.prefix}cf-ssh"
  target      = "${google_compute_target_pool.cf-public.self_link}"
  port_range  = "2222"
  ip_protocol = "TCP"
  ip_address  = "${google_compute_address.cf.address}"
}

// WSS forwarding rule
resource "google_compute_forwarding_rule" "cf-wss" {
  name        = "${var.prefix}cf-wss"
  target      = "${google_compute_target_pool.cf-public.self_link}"
  port_range  = "4443"
  ip_protocol = "TCP"
  ip_address  = "${google_compute_address.cf.address}"
}

// TCP forwarding rule
resource "google_compute_forwarding_rule" "cf-tcp" {
  name        = "${var.prefix}cf-tcp"
  target      = "${google_compute_target_pool.cf-tcp-public.self_link}"
  port_range  = "1024-1123"
  ip_protocol = "TCP"
  ip_address  = "${google_compute_address.cf-tcp.address}"
}

output "ip" {
    value = "${google_compute_address.cf.address}"
}

output "tcp_ip" {
    value = "${google_compute_address.cf-tcp.address}"
}

output "network" {
   value = "${var.network}"
}

output "private_subnet" {
   value = "${google_compute_subnetwork.cf-private-subnet-1.name}"
}

output "compilation_subnet" {
   value = "${google_compute_subnetwork.cf-compilation-subnet-1.name}"
}

output "zone" {
  value = "${var.zone}"
}

output "zone_compilation" {
  value = "${var.zone_compilation}"
}

output "region" {
  value = "${var.region}"
}

output "region_compilation" {
  value = "${var.region_compilation}"
}


