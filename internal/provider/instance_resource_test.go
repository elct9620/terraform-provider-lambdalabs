package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_InstanceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "lambdalabs_instance" "default" {
					region_name        = "us-west-1"
					instance_type_name = "gpu_1x_a10"
					ssh_key_names = [
						"terraform"
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "region_name", "us-west-1"),
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "instance_type_name", "gpu_1x_a10"),
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "ssh_key_names[0]", "terraform"),
				),
			},
			{
				Config: `
				resource "lambdalabs_instance" "default" {
					region_name        = "us-west-1"
					instance_type_name = "gpu_1x_a100"
					ssh_key_names = [
						"terraform"
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "instance_type_name", "gpu_1x_a100"),
				),
			},
		},
	})
}
