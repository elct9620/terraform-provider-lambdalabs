package provider_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_SSHKeyResource(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(r.URL.Path)
		if strings.HasPrefix(path, "/api/v1/ssh-keys") {
			http.NotFoundHandler().ServeHTTP(w, r)
		}

		switch r.Method {
		case http.MethodGet:
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
		case http.MethodPost:
			var input struct {
				Name string `json:"name"`
			}

			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &input) //nolint:errcheck

			resBody := fmt.Sprintf(`
			{
				"data": {
					"id": "0920582c7ff041399e34823a0be62548",
					"name": %[1]q,
					"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDfKpav4ILY54InZe27G user",
					"private_key": "-----BEGIN RSA PRIVATE KEY-----\nKEY CONTENT-----END RSA PRIVATE KEY-----\n"
				}
			}
			`, input.Name)
			w.Write([]byte(resBody)) //nolint:errcheck
		case http.MethodDelete:
			w.Write(json.RawMessage(`{ "data": {} }`)) //nolint:errcheck
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(server.URL) + `
				resource "lambdalabs_ssh_key" "default" {
					name = "terraform"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lambdalabs_ssh_key.default", "name", "terraform"),
					resource.TestCheckResourceAttr("lambdalabs_ssh_key.default", "id", "0920582c7ff041399e34823a0be62548"),
					resource.TestCheckResourceAttr("lambdalabs_ssh_key.default", "public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDfKpav4ILY54InZe27G user"),
					resource.TestCheckResourceAttr("lambdalabs_ssh_key.default", "private_key", ""), // Private key is sensitive and should not be stored
				),
			},
			{
				ResourceName:      "lambdalabs_ssh_key.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
