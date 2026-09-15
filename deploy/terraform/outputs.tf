output "droplet_ip" {
  description = "Public IPv4 of the droplet."
  value       = digitalocean_droplet.app.ipv4_address
}

output "ssh_command" {
  description = "Convenience SSH command."
  value       = "ssh -i ~/.ssh/shortlink_deploy root@${digitalocean_droplet.app.ipv4_address}"
}
