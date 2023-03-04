package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_SSHKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "lambdalabs_ssh_key" "default" {
					name = "terraform"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lambdalabs_ssh_key.default", "name", "terraform"),
				),
			},
			{
				Config: `
				resource "lambdalabs_ssh_key" "default" {
					name = "terraform-auto"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lambdalabs_ssh_key.default", "name", "terraform-auto"),
				),
			},
		},
	})
}
