package provider_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_InstanceResource(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/instances/0920582c7ff041399e34823a0be62549":
			resBody := `
			{
				"data": {
					"id": "0920582c7ff041399e34823a0be62549",
					"name": "",
					"ip": "10.10.10.1",
					"status": "active",
					"ssh_key_names": [
						"terraform"
					],
					"region": {
						"name": "us-tx-1",
						"description": "Austin, Texas"
					},
					"instance_type": {
						"name": "gpu_1x_a100",
						"description": "1x RTX A100 (24 GB)",
						"price_cents_per_hour": 110,
						"specs": {
							"vcpus": 24,
							"memory_gib": 800,
							"storage_gib": 512
						}
					}
				}
			}
			`
			w.Write([]byte(resBody)) //nolint:errcheck
		case "/instance-operations/launch":
			resBody := `
			{
				"data": {
					"instance_ids": [
					"0920582c7ff041399e34823a0be62549"
					]
				}
			}
			`
			w.Write([]byte(resBody)) //nolint:errcheck
		case "/instance-operations/terminate":
			resBody := `
			{
				"data": {
					"terminated_instances": [
					{
						"id": "0920582c7ff041399e34823a0be62549",
						"name": "training-node-1",
						"ip": "10.10.10.1",
						"status": "active",
						"ssh_key_names": [],
						"file_system_names": [],
						"region": {},
						"instance_type": {}
					}
					]
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
				resource "lambdalabs_instance" "default" {
					region_name        = "us-tx-1"
					instance_type_name = "gpu_1x_a100"
					ssh_key_names = [
						"terraform"
					]
					timeouts {
						create = "10s"
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "id", "0920582c7ff041399e34823a0be62549"),
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "ip", "10.10.10.1"),
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "region_name", "us-tx-1"),
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "instance_type_name", "gpu_1x_a100"),
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "ssh_key_names.0", "terraform"),
					resource.TestCheckResourceAttr("lambdalabs_instance.default", "timeouts.create", "10s"),
				),
			},
		},
	})
}
