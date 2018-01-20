output "google_project" {
  value = "${var.google_project}"
}

output "google_region" {
  value = "${var.google_region}"
}

output "google_zone" {
  value = "${var.google_zone}"
}

output "google_json_key_data" {
  value = "${var.google_json_key_data}"
}

output "google_auto_network" {
  value = "${google_compute_network.auto.name}"
}

output "google_network" {
  value = "${google_compute_network.manual.name}"
}

output "google_subnetwork" {
  value = "${google_compute_subnetwork.subnetwork.name}"
}

output "google_firewall_internal" {
  value = "${var.google_firewall_internal}"
}

output "google_firewall_external" {
  value = "${var.google_firewall_external}"
}

output "google_backend_service" {
  value = "${google_compute_backend_service.backend_service.name}"
}

output "google_region_backend_service" {
  value = "${google_compute_region_backend_service.region_backend_service.name}"
}

output "google_target_pool" {
  value = "${google_compute_target_pool.regional.name}"
}

output "google_address_director_ubuntu" {
  value = "${google_compute_address.director.name}"
}

output "google_address_bats_ubuntu" {
  value = "${google_compute_address.bats.name}"
}

output "google_address_int_ubuntu" {
  value = "${google_compute_address.int.name}"
}

output "google_service_account" {
  value = "${google_service_account.service_account.email}"
}
