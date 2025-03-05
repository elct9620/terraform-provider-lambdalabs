package provider_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_SSHKeyData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSpace(r.URL.Path) == "/ssh-keys" {
			resBody := `
			{
				"data": [
				{
					"id": "0920582c7ff041399e34823a0be62548",
					"name": "terraform",
					"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDfKpav4ILY54InZe27G user",
					"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEA5IGybv8/wdQM6Y4yYTGiEem4TscBZiAW+9xyW2pDt8S7VDtm\n...\neCW4938W9u8N3R/kpGwi1tZYiGMLBU4Ks0qKFi/VeEaE9OLeP5WQ8Pk=\n-----END RSA PRIVATE KEY-----\n"
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
				data "lambdalabs_ssh_key" "default" {
					name = "terraform"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdalabs_ssh_key.default", "name", "terraform"),
					resource.TestCheckResourceAttr("data.lambdalabs_ssh_key.default", "id", "0920582c7ff041399e34823a0be62548"),
					resource.TestCheckResourceAttr("data.lambdalabs_ssh_key.default", "public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDfKpav4ILY54InZe27G user"),
				),
			},
		},
	})
}
