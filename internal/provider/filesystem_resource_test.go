package provider_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_FilesystemResource(t *testing.T) {
	t.Parallel()

	filesystemId := "fs-12345678"
	filesystemName := "test-filesystem"
	region := "us-west-1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/file-systems":
			if r.Method == http.MethodGet {
				// List file systems
				resBody := fmt.Sprintf(`
				{
					"data": [
						{
							"id": "%s",
							"name": "%s",
							"mount_point": "/mnt/data",
							"created": "2023-01-01T00:00:00.000Z",
							"created_by": {
								"id": "user-12345",
								"email": "user@example.com",
								"status": "active"
							},
							"is_in_use": false,
							"bytes_used": 0,
							"region": {
								"name": "%s",
								"description": "California, USA"
							}
						}
					]
				}
				`, filesystemId, filesystemName, region)
				w.Write([]byte(resBody)) //nolint:errcheck
			}
		case "/filesystems":
			if r.Method == http.MethodPost {
				// Create file system
				resBody := fmt.Sprintf(`
				{
					"data": {
						"id": "%s",
						"name": "%s",
						"mount_point": "/mnt/data",
						"created": "2023-01-01T00:00:00.000Z",
						"created_by": {
							"id": "user-12345",
							"email": "user@example.com",
							"status": "active"
						},
						"is_in_use": false,
						"bytes_used": 0,
						"region": {
							"name": "%s",
							"description": "California, USA"
						}
					}
				}
				`, filesystemId, filesystemName, region)
				w.Write([]byte(resBody)) //nolint:errcheck
			}
		case fmt.Sprintf("/filesystems/%s", filesystemId):
			if r.Method == http.MethodDelete {
				// Delete file system
				resBody := fmt.Sprintf(`
				{
					"data": {
						"deleted_ids": ["%s"]
					}
				}
				`, filesystemId)
				w.Write([]byte(resBody)) //nolint:errcheck
			}
		}
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(server.URL) + fmt.Sprintf(`
				resource "lambdalabs_filesystem" "test" {
					name   = "%s"
					region = "%s"
				}
				`, filesystemName, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lambdalabs_filesystem.test", "id", filesystemId),
					resource.TestCheckResourceAttr("lambdalabs_filesystem.test", "name", filesystemName),
					resource.TestCheckResourceAttr("lambdalabs_filesystem.test", "region", region),
				),
			},
		},
	})
}
