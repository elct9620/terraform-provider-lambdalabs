package provider_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_InstanceTypesData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSpace(r.URL.Path) == "/instance-types" {
			resBody := `
			{
				"data": {
					"gpu_1x_a100": {
						"instance_type": {
							"name": "gpu_1x_a100",
							"description": "1x NVIDIA A100 (80 GB)",
							"gpu_description": "NVIDIA A100 (80 GB)",
							"price_cents_per_hour": 199,
							"specs": {
								"vcpus": 30,
								"memory_gib": 200,
								"storage_gib": 1024,
								"gpus": 1
							}
						},
						"regions_with_capacity_available": [
							{
								"name": "us-west-1",
								"description": "California, USA"
							}
						]
					}
				}
			}
			`
			w.Write([]byte(resBody)) //nolint:errcheck
		}
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(server.URL) + `
				data "lambdalabs_instance_types" "default" {
					filter = {
						region = "us-west-1"
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdalabs_instance_types.default", "instance_types.gpu_1x_a100.name", "gpu_1x_a100"),
					resource.TestCheckResourceAttr("data.lambdalabs_instance_types.default", "instance_types.gpu_1x_a100.description", "1x NVIDIA A100 (80 GB)"),
					resource.TestCheckResourceAttr("data.lambdalabs_instance_types.default", "instance_types.gpu_1x_a100.gpu_description", "NVIDIA A100 (80 GB)"),
					resource.TestCheckResourceAttr("data.lambdalabs_instance_types.default", "instance_types.gpu_1x_a100.price_cents_per_hour", "199"),
					resource.TestCheckResourceAttr("data.lambdalabs_instance_types.default", "instance_types.gpu_1x_a100.specs.vcpus", "30"),
					resource.TestCheckResourceAttr("data.lambdalabs_instance_types.default", "instance_types.gpu_1x_a100.specs.memory_gib", "200"),
					resource.TestCheckResourceAttr("data.lambdalabs_instance_types.default", "instance_types.gpu_1x_a100.specs.storage_gib", "1024"),
					resource.TestCheckResourceAttr("data.lambdalabs_instance_types.default", "instance_types.gpu_1x_a100.specs.gpus", "1"),
				),
			},
		},
	})
}
