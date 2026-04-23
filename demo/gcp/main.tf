provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

locals {
  instance_names = [
    "tinytorrent-bootstrap",
    "tinytorrent-peer-b",
    "tinytorrent-peer-c",
  ]

  common_tags = ["tinytorrent-demo"]
}

resource "google_compute_network" "tinytorrent_demo" {
  name                    = "tinytorrent-demo-network"
  auto_create_subnetworks = true
}

resource "google_compute_firewall" "tinytorrent_demo_ingress" {
  name    = "tinytorrent-demo-ingress"
  network = google_compute_network.tinytorrent_demo.name

  allow {
    protocol = "tcp"
    ports    = ["22", "4001-4010"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = local.common_tags
}

resource "google_compute_instance" "tinytorrent_demo" {
  for_each     = toset(local.instance_names)
  name         = each.value
  machine_type = var.machine_type
  zone         = var.zone
  tags         = local.common_tags

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/family/debian-12"
      size  = 20
      type  = "pd-balanced"
    }
  }

  network_interface {
    network = google_compute_network.tinytorrent_demo.id

    access_config {}
  }
}
