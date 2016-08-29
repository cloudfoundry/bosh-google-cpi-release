// Subnet for the public Concourse components
resource "google_compute_subnetwork" "concourse-public-subnet-1" {
  name          = "concourse-public-${var.region}-1"
  ip_cidr_range = "10.150.0.0/16"
  network       = "${google_compute_network.network.self_link}"
}

// Subnet for the public Concourse components
resource "google_compute_subnetwork" "concourse-public-subnet-2" {
  name          = "concourse-public-${var.region}-2"
  ip_cidr_range = "10.160.0.0/16"
  network       = "${google_compute_network.network.self_link}"
}

// Allow access to CloudFoundry router
resource "google_compute_firewall" "concourse-public" {
  name    = "concourse-public"
  network = "${google_compute_network.network.name}"

  allow {
    protocol = "tcp"
    ports    = ["8080"]
  }

  target_tags = ["concourse-public"]
}

// Allow open access to between internal VMs
resource "google_compute_firewall" "concourse-internal" {
  name    = "concourse-internal"
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

  target_tags = ["concourse-internal", "bosh-internal"]
  source_tags = ["concourse-internal", "bosh-internal"]
}

// Static IP address for forwarding rule
resource "google_compute_global_address" "concourse" {
  name = "concourse"
}

// Instance group 1
resource "google_compute_instance_group" "concourse-1" {
  name        = "concourse-${var.zone-1}-1"
  zone        = "${var.zone-1}"
  named_port {
    name = "http"
    port = "8080"
  }
}

// Instance group 2
resource "google_compute_instance_group" "concourse-2" {
  name        = "concourse-${var.zone-2}-2"
  zone        = "${var.zone-2}"
  named_port {
    name = "http"
    port = "8080"
  }
}

resource "google_compute_backend_service" "concourse" {
  name        = "concourse"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10
  enable_cdn  = false

  backend {
    group = "${google_compute_instance_group.concourse-1.self_link}"
  }
  backend {
    group = "${google_compute_instance_group.concourse-2.self_link}"
  }

  health_checks = ["${google_compute_http_health_check.concourse.self_link}"]
}

resource "google_compute_http_health_check" "concourse" {
  name               = "concourse"
  request_path       = "/login"
  check_interval_sec = 1
  timeout_sec        = 1
  port               = 8080
}

resource "google_compute_target_http_proxy" "concourse" {
  name             = "concourse"
  url_map          = "${google_compute_url_map.concourse.self_link}"
}

resource "google_compute_url_map" "concourse" {
  name        = "concourse"

  default_service = "${google_compute_backend_service.concourse.self_link}"
}

resource "google_compute_global_forwarding_rule" "concourse" {
  name       = "concourse"
  ip_address    = "${google_compute_global_address.concourse.address}"
  target     = "${google_compute_target_http_proxy.concourse.self_link}"
  port_range = "80"
}
