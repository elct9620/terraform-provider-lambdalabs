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

output "key" {
  value = lambdalabs_ssh_key.primary
  sensitive = true
}
