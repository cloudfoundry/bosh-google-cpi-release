provider "google" {
  project     = var.google_project
  region      = var.google_region
  credentials = var.google_json_key_data
}

provider "random" {
}
