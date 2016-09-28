provider "google" {
    project = "${var.gce_project_id}"
    region = "${var.gce_region}"
    credentials = "${var.gce_credentials_json}"
}

resource "google_compute_network" "skeletal" {
  name = "stemcell-ci-${var.env_name}"
}

// Subnet for the skeletal deployment
resource "google_compute_subnetwork" "skeletal" {
  name          = "stemcell-ci-${var.env_name}"
  ip_cidr_range = "10.0.0.0/24"
  network       = "${google_compute_network.skeletal.self_link}"
}

// Allow SSH & MBus to skeletal deployment
resource "google_compute_firewall" "allow-ssh-mbus" {
  name    = "allow-ssh-mbus-${var.env_name}"
  network = "${google_compute_network.skeletal.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22", "6868"]
  }

  target_tags = ["skeletal-external-${var.env_name}"]
}

resource "google_compute_address" "skeletal" {
  name = "skeletal-ip-${var.env_name}"
}
