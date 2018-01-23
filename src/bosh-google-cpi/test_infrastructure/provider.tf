provider "google" {
  version = "v1.5.0"

  project     = "${var.google_project}"
  region      = "${var.google_region}"
  credentials = "${var.google_json_key_data}"
}

provider "random" {
  version = "~> 1.1.0"
}
