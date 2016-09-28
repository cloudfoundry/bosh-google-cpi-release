variable "gce_project_id" {
  type = "string"
}

variable "gce_credentials_json" {
  type = "string"
}

variable "gce_region" {
  type = "string"
  default = "us-central1"
}

variable "gce_zone" {
  type = "string"
  default = "us-central1-f"
}

variable "env_name" {
  type = "string"
}

