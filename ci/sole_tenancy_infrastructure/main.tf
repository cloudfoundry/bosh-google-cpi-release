resource "google_compute_node_template" "soletenant-tmpl" {
  name      = "${var.prefix}-node-group-template"
  region    = var.google_region
  node_type = "c2-node-60-240"
}

resource "google_compute_node_group" "nodes" {
  name          = "${var.prefix}-node-group"
  zone          = var.google_zone
  initial_size  = 1
  node_template = google_compute_node_template.soletenant-tmpl.id
}
