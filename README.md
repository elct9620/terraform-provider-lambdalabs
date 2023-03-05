Terraform Provider - Lambdalabs
===

The Terraform provider for Lambdalabs.

## Example

To set up a [Stable Diffusion WebUI](https://github.com/AUTOMATIC1111/stable-diffusion-webui) in Lambdalabs, you can use terraform to provision everything.

```hcl
terraform {
  required_providers {
    lambdalabs = {
      source = "elct9620/lambdalabs"
    }
  }
}

provider "lambdalabs" {
  # Or use `LAMBDALABS_API_KEY`
  api_key = "API_KEY"
}

resource "lambdalabs_instance" "stable_diffusion" {
  region_name        = "us-west-1"
  instance_type_name = "gpu_1x_a10"
  # Suggest creating a key instead of using resource to make it reusable
  ssh_key_names = [
    "terraform"
  ]

  connection {
    type     = "ssh"
      user     = "ubuntu"
      private_key = file("~/.ssh/id_ed25519")
      host     = self.ip
  }

  provisioner "remote-exec" {
    inline = [
      "pip3 install xformers",
      "wget -qO- https://raw.githubusercontent.com/AUTOMATIC1111/stable-diffusion-webui/master/webui.sh | COMMANDLINE_ARGS=\"--exit --xformers\" bash",
    ]
  }
}
```
After the server is created, you can ssh into the server and run `./webui.sh --listen --xformers` to use WebUI.
