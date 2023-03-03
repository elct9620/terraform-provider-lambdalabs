terraform {
  required_providers {
    lambdalabs = {
      source = "elct9620/lambdalabs"
    }
  }
}

provider "lambdalabs" {}

resource "lambdalabs_ssh_key" "primary" {
  name = "terraform"
}

resource "lambdalabs_instance" "stable_diffusion" {
  region_name = "us-west-1"
  instance_type_name = "gpu_1x_a10"
  ssh_key_names = [
    lambdalabs_ssh_key.primary.name
  ]
}
