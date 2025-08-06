variable "prefix" {
  description = "A prefix to apply to all created resources."
  type        = string
}

variable "google_project" {
  description = "The Google Cloud project to deploy to."
  type        = string
}

variable "google_region" {
  description = "The Google Cloud region to deploy to."
  type        = string
}

variable "google_zone" {
  description = "The Google Cloud zone to deploy to."
  type        = string
}

variable "google_json_key_data" {
  description = "The GCP service account key."
  type        = string
}
