provider "google" {
  version = "~> 2.5"

  project     = var.google_project
  region      = var.google_region
  credentials = var.google_json_key_data
}

provider "random" {
  version = "~> 2.1"
}
