output "gce_project_id" {
  value = "${var.gce_project_id}"
}

output "gce_region" {
  value = "${var.gce_region}"
}

output "gce_zone" {
  value = "${var.gce_zone}"
}

output "network_name" {
  value = "${google_compute_network.skeletal.name}"
}

output "subnetwork_name" {
  value = "${google_compute_subnetwork.skeletal.name}"
}

output "skeletal_external_ip" {
  value = "${google_compute_address.skeletal.address}"
}

output "skeletal_firewall_tag" {
  value = "skeletal-external-${var.env_name}"
}
