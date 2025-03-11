package provider_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_ImagesData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/images" {
			resBody := `
			{
				"data": [
					{
						"id": "43336648-096d-4cba-9aa2-f9bb7727639d",
						"created_time": "1970-01-01T00:00:00.000Z",
						"updated_time": "1970-01-01T00:00:00.000Z",
						"name": "ubuntu-24.04.01",
						"description": "Ubuntu LTS",
						"family": "ubuntu-lts",
						"version": "24.04.01",
						"architecture": "x86_64",
						"region": {
							"name": "us-west-1",
							"description": "California, USA"
						}
					}
				]
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
				data "lambdalabs_images" "all" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "id", "images"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.#", "1"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.id", "43336648-096d-4cba-9aa2-f9bb7727639d"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.name", "ubuntu-24.04.01"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.description", "Ubuntu LTS"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.family", "ubuntu-lts"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.version", "24.04.01"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.created_time", "1970-01-01T00:00:00.000Z"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.updated_time", "1970-01-01T00:00:00.000Z"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.region.name", "us-west-1"),
					resource.TestCheckResourceAttr("data.lambdalabs_images.all", "images.0.region.description", "California, USA"),
				),
			},
		},
	})
}
