variable "project_name" {
  description = "Name/tag for the droplet and firewall."
  type        = string
  default     = "shortlink"
}

variable "region" {
  description = "DigitalOcean region slug (sgp1 = Singapore, closest to Vietnam)."
  type        = string
  default     = "sgp1"
}

variable "droplet_size" {
  description = "Droplet size slug. s-2vcpu-4gb fits the app + monitoring + Jenkins."
  type        = string
  default     = "s-2vcpu-4gb"
}

variable "droplet_image" {
  description = "Base image slug."
  type        = string
  default     = "ubuntu-24-04-x64"
}

variable "ssh_public_key_path" {
  description = "Path to the deploy SSH public key registered on the droplet."
  type        = string
  default     = "~/.ssh/shortlink_deploy.pub"
}

variable "ssh_key_name" {
  description = "Name of the SSH key resource in DigitalOcean."
  type        = string
  default     = "shortlink-deploy"
}

variable "allowed_ssh_cidr" {
  description = "CIDR allowed to reach SSH (22). Lock to your IP in real use."
  type        = string
  default     = "0.0.0.0/0"
}
