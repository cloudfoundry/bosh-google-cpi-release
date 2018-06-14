resource "random_string" "account_suffix" {
  length  = 4
  upper   = false
  special = false
  lower   = true
  number  = true
}

resource "google_service_account" "service_account" {
  account_id = "${var.prefix}-sa-${random_string.account_suffix.result}"
}

resource "google_compute_address" "director" {
  name = "${var.prefix}-dir"
}

resource "google_compute_address" "bats" {
  name = "${var.prefix}-bats"
}

resource "google_compute_address" "int" {
  name = "${var.prefix}-int"
}

resource "google_compute_network" "auto" {
  name                    = "${var.google_auto_network}"
  auto_create_subnetworks = true
}

resource "google_compute_network" "manual" {
  name                    = "${var.google_network}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "${var.prefix}"
  ip_cidr_range = "${var.google_subnetwork_range}"
  network       = "${google_compute_network.manual.self_link}"
}

resource "google_compute_firewall" "internal" {
  name        = "${var.prefix}-int"
  description = "BOSH CI Internal Traffic"
  network     = "${google_compute_network.manual.self_link}"
  source_tags = ["${var.google_firewall_internal}"]
  target_tags = ["${var.google_firewall_internal}"]

  allow {
    protocol = "tcp"
  }

  allow {
    protocol = "udp"
  }

  allow {
    protocol = "icmp"
  }
}

resource "google_compute_firewall" "external" {
  name        = "${var.prefix}-ext"
  description = "BOSH CI External Traffic"
  network     = "${google_compute_network.manual.self_link}"
  target_tags = ["${var.google_firewall_external}"]

  allow {
    protocol = "tcp"
    ports    = ["22", "443", "4222", "6868", "25250", "25555", "25777"]
  }

  allow {
    protocol = "udp"
    ports    = ["53"]
  }

  allow {
    protocol = "icmp"
  }
}

# Target Pool
resource "google_compute_target_pool" "regional" {
  name   = "${var.prefix}-r"
  region = "${var.google_region}"
}

# Backend Service
resource "google_compute_instance_group" "backend_service" {
  name = "${var.prefix}"
  zone = "${var.google_zone}"
}

resource "google_compute_http_health_check" "backend_service" {
  name = "${var.prefix}"
}

resource "google_compute_backend_service" "backend_service" {
  health_checks = ["${google_compute_http_health_check.backend_service.self_link}"]
  name          = "${var.prefix}"
  port_name     = "http"
  timeout_sec   = "30"

  backend {
    group           = "${google_compute_instance_group.backend_service.self_link}"
    balancing_mode  = "UTILIZATION"
    capacity_scaler = "1"
    max_utilization = "0.8"
  }
}

# Regional Backend Service
resource "google_compute_health_check" "region_backend_service" {
  name = "${var.prefix}-r"

  tcp_health_check {
    port = "8080"
  }
}

resource "google_compute_instance_group" "region_backend_service" {
  name      = "${var.prefix}-r"
  zone      = "${var.google_zone}"
  instances = ["${google_compute_instance.hack.self_link}"]
}

// HACK to work around: googleapi: Error 400: Invalid value for field 'resource.backends[0].group': 'https://www.googleapis.com/compute/v1/projects/pivotal-cloudfoundry/zones/us-west1-a/instanceGroups/ci-bosh-deployment-r'. Instance group must have a network to be attached to a backend service. Add an instance to give the instance group a network., invalid
resource "google_compute_instance" "hack" {
  boot_disk = {
    initialize_params {
      image = "debian-cloud/debian-8"
    }
  }

  machine_type = "f1-micro"
  name         = "${var.prefix}-hack"

  network_interface = {
    subnetwork = "${google_compute_subnetwork.subnetwork.self_link}"
  }

  zone = "${var.google_zone}"
}

resource "google_compute_region_backend_service" "region_backend_service" {
  name          = "${var.prefix}-r"
  health_checks = ["${google_compute_health_check.region_backend_service.self_link}"]
  region        = "${var.google_region}"
  protocol      = "TCP"
  timeout_sec   = "30"

  backend {
    group = "${google_compute_instance_group.region_backend_service.self_link}"
  }
}


resource "google_compute_backend_service" "collision_backend_service" {
  health_checks = ["${google_compute_http_health_check.backend_service.self_link}"]
  name          = "${var.prefix}"
  port_name     = "http"
  timeout_sec   = "30"

  backend {
    group           = "${google_compute_instance_group.backend_service.self_link}"
    balancing_mode  = "UTILIZATION"
    capacity_scaler = "1"
    max_utilization = "0.8"
  }
}

resource "google_compute_region_backend_service" "collision_region_backend_service" {
  name          = "${var.prefix}"
  health_checks = ["${google_compute_health_check.region_backend_service.self_link}"]
  region        = "${var.google_region}"
  protocol      = "TCP"
  timeout_sec   = "30"

  backend {
    group = "${google_compute_instance_group.region_backend_service.self_link}"
  }
}