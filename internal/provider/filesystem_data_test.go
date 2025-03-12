package provider_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_FilesystemsData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/file-systems" {
			resBody := `
			{
				"data": [
					{
						"id": "fs-12345678",
						"name": "data-storage",
						"mount_point": "/mnt/data",
						"created": "2023-01-01T00:00:00.000Z",
						"created_by": {
							"id": "user-12345",
							"email": "user@example.com",
							"status": "active"
						},
						"is_in_use": true,
						"bytes_used": 1073741824,
						"region": {
							"name": "us-west-1",
							"description": "California, USA"
						}
					},
					{
						"id": "fs-87654321",
						"name": "model-storage",
						"mount_point": "/mnt/models",
						"created": "2023-02-01T00:00:00.000Z",
						"created_by": {
							"id": "user-12345",
							"email": "user@example.com",
							"status": "active"
						},
						"is_in_use": false,
						"bytes_used": 2147483648,
						"region": {
							"name": "us-east-1",
							"description": "Virginia, USA"
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
				data "lambdalabs_filesystems" "all" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "id", "filesystems"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.#", "2"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.0.id", "fs-12345678"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.0.name", "data-storage"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.0.mount_point", "/mnt/data"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.0.is_in_use", "true"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.0.bytes_used", "1073741824"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.0.region.name", "us-west-1"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.0.created_by.email", "user@example.com"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.1.id", "fs-87654321"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.1.name", "model-storage"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.all", "filesystems.1.region.name", "us-east-1"),
				),
			},
			{
				Config: providerConfig(server.URL) + `
				data "lambdalabs_filesystems" "filtered_by_region" {
					filter = {
						region = "us-west-1"
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.filtered_by_region", "id", "filesystems"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.filtered_by_region", "filesystems.#", "1"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.filtered_by_region", "filesystems.0.id", "fs-12345678"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.filtered_by_region", "filesystems.0.name", "data-storage"),
					resource.TestCheckResourceAttr("data.lambdalabs_filesystems.filtered_by_region", "filesystems.0.region.name", "us-west-1"),
				),
			},
		},
	})
}
