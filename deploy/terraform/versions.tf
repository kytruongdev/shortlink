terraform {
  required_version = ">= 1.6"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# The token is read from the DIGITALOCEAN_TOKEN environment variable,
# so it never lives in a file. Export it before running terraform:
#   export DIGITALOCEAN_TOKEN=dop_v1_...
provider "digitalocean" {}
