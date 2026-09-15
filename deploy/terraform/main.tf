# Register the deploy SSH public key, create the droplet, and firewall it.

resource "digitalocean_ssh_key" "deploy" {
  name       = var.ssh_key_name
  public_key = file(pathexpand(var.ssh_public_key_path))
}

resource "digitalocean_droplet" "app" {
  name     = var.project_name
  region   = var.region
  size     = var.droplet_size
  image    = var.droplet_image
  ssh_keys = [digitalocean_ssh_key.deploy.fingerprint]
  tags     = [var.project_name]
}

resource "digitalocean_firewall" "app" {
  name        = "${var.project_name}-fw"
  droplet_ids = [digitalocean_droplet.app.id]

  # SSH (lock to your IP via allowed_ssh_cidr in real use).
  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = [var.allowed_ssh_cidr]
  }
  # HTTP / HTTPS (public app).
  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }
  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }
  # Grafana (3000) and Prometheus (9090) for the monitoring demo.
  inbound_rule {
    protocol         = "tcp"
    port_range       = "3000"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }
  inbound_rule {
    protocol         = "tcp"
    port_range       = "9090"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  # Allow all outbound.
  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}
