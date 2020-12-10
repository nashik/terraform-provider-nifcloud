terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_load_balancer" "l4lb" {
  accounting_type = "1"
  load_balancer_name = "nl4lb"
  load_balancer_port = "80"
  instance_port = "80"
  filter_type = "1"
}
