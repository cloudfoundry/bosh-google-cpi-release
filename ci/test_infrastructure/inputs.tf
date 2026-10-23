variable "google_project" {
}

variable "google_region" {
}

variable "google_zone" {
}

variable "google_json_key_data" {
}

variable "google_subnetwork_range" {
}

variable "google_firewall_internal" {
}

variable "google_firewall_external" {
}

variable "prefix" {
}

variable "google_external_access_cidrs" {
  description = "CIDR ranges of trusted CI worker egress addresses allowed to access test VMs."
  type        = list(string)

  validation {
    condition = length(var.google_external_access_cidrs) > 0 && alltrue([
      for cidr in var.google_external_access_cidrs :
      can(cidrhost(cidr, 0)) && try(tonumber(split("/", cidr)[1]), 0) > 0
    ])
    error_message = "Provide non-empty, explicit CI worker CIDRs; unrestricted IPv4 or IPv6 ingress is not allowed."
  }
}
